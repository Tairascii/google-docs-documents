package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Tairascii/google-docs-documents/internal/app"
	"github.com/Tairascii/google-docs-documents/internal/app/service/document"
	"github.com/Tairascii/google-docs-documents/pkg"
	"github.com/dancannon/gorethink"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

// TODO move to apigw and use vault
const (
	accessSecret = "yoS0baK1Ya"
)

var (
	ErrAuth           = errors.New("authentication failed")
	ErrInvalidRequest = errors.New("invalid request")
	ErrNotAllowed     = errors.New("not allowed")
)

type Handler struct {
	DI       *app.DI
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
	mu       *sync.Mutex
	session  *gorethink.Session
}

func NewHandler(di *app.DI, session *gorethink.Session) *Handler {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	var clients = make(map[*websocket.Conn]bool)
	return &Handler{
		DI:       di,
		upgrader: upgrader,
		clients:  clients,
		mu:       &sync.Mutex{},
		session:  session,
	}
}

func (h *Handler) InitHandlers() *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))
	r.Route("/api", func(api chi.Router) {
		api.Route("/v1", func(v1 chi.Router) {
			v1.Mount("/documents", handlers(h))
			v1.HandleFunc("/document/ws/{id}", h.ConnectWS)
		})
	})
	return r
}

func handlers(h *Handler) http.Handler {
	rg := chi.NewRouter()
	rg.Use(ParseToken(accessSecret))
	rg.Group(func(r chi.Router) {
		r.Get("/", h.GetDocuments)
		r.Post("/", h.CreateDocument)
		r.Delete("/{id}", h.DeleteDocument)
		r.Put("/{id}", h.EditDocument)
	})
	return rg
}

func (h *Handler) GetDocuments(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	ctx := r.Context()
	res, err := h.DI.UseCase.Documents.GetDocuments(ctx, search)
	if err != nil {
		pkg.JSONErrorResponseWriter(w, err, http.StatusInternalServerError)
		return
	}
	pkg.JSONResponseWriter[[]Document](w, toDocuments(res), http.StatusOK)
}

func (h *Handler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	var payload CreateDocumentPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		pkg.JSONErrorResponseWriter(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ctx := r.Context()
	id, err := h.DI.UseCase.Documents.CreateDocument(ctx, payload.Title, payload.InitialContent)
	if err != nil {
		pkg.JSONErrorResponseWriter(w, err, http.StatusBadRequest)
		return
	}

	pkg.JSONResponseWriter[CreateDocumentResponse](w, CreateDocumentResponse{DocumentID: id}, http.StatusOK)
}

func (h *Handler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		pkg.JSONErrorResponseWriter(w, ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err := h.DI.UseCase.Documents.DeleteDocument(ctx, id)
	if err != nil {
		if errors.Is(err, document.ErrNotAllowed) {
			pkg.JSONErrorResponseWriter(w, ErrNotAllowed, http.StatusForbidden)
			return
		}
		pkg.JSONErrorResponseWriter(w, err, http.StatusInternalServerError)
		return
	}

	pkg.JSONResponseWriter[any](w, nil, http.StatusNoContent)
}

func (h *Handler) EditDocument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		pkg.JSONErrorResponseWriter(w, ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	var payload EditDocumentPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		pkg.JSONErrorResponseWriter(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if payload.Title == "" {
		pkg.JSONErrorResponseWriter(w, ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err := h.DI.UseCase.Documents.EditDocument(ctx, id, payload.Title)
	if err != nil {
		if errors.Is(err, document.ErrNotAllowed) {
			pkg.JSONErrorResponseWriter(w, ErrNotAllowed, http.StatusForbidden)
			return
		}
		pkg.JSONErrorResponseWriter(w, err, http.StatusInternalServerError)
		return
	}

	pkg.JSONResponseWriter[any](w, nil, http.StatusNoContent)
}

func (h *Handler) ConnectWS(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		pkg.JSONErrorResponseWriter(w, ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	err := h.DI.UseCase.Documents.CheckPermission(ctx, id)
	if err != nil {
		pkg.JSONErrorResponseWriter(w, err, http.StatusForbidden)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		pkg.JSONErrorResponseWriter(w, err, http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()
	defer func() {
		h.mu.Lock()
		delete(h.clients, conn)
		h.mu.Unlock()
	}()

	go h.WatchTableChange(conn, id)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read Error:", err)
			break
		}
		log.Printf("Received: %s\n", msg)

		if err := h.DI.UseCase.Documents.SaveDocumentContent(context.TODO(), id, msg); err != nil {
			log.Printf("Save Error: %s\n", err)
			break
		}
	}
}

// TODO think of more clear way
func (h *Handler) WatchTableChange(conn *websocket.Conn, id string) {
	cursor, err := gorethink.Table("documents").Get(id).Changes().Run(h.session)
	if err != nil {
		log.Println("Error watching changes:", err)
		return
	}
	defer cursor.Close()

	var change map[string]interface{}
	for cursor.Next(&change) {
		fmt.Println("Document updated:", change)

		newValue := change["new_val"].(map[string]interface{})
		docID := newValue["id"].(string)
		content, ok := newValue["content"].([]byte)
		if !ok {
			log.Println("Error converting content to bytes:", newValue["content"])
			break
		}
		if err := conn.WriteMessage(websocket.BinaryMessage, content); err != nil {
			log.Println("Write Error:", err)
			break
		}
		fmt.Println("content send")
		fmt.Printf("new data: %v %v \n", docID, string(content))
	}
	return
}
