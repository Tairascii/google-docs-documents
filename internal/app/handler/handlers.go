package handler

import (
	"encoding/json"
	"errors"
	"github.com/Tairascii/google-docs-documents/internal/app"
	"github.com/Tairascii/google-docs-documents/pkg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
)

// TODO move to apigw and use vault
const (
	accessSecret = "yoS0baK1Ya"
)

var (
	ErrAuth = errors.New("authentication failed")
)

type Handler struct {
	DI *app.DI
}

func NewHandler(di *app.DI) *Handler {
	return &Handler{DI: di}
}

func (h *Handler) InitHandlers() *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))
	r.Use(ParseToken(accessSecret))
	r.Route("/api", func(api chi.Router) {
		api.Route("/v1", func(v1 chi.Router) {
			v1.Mount("/documents", handlers(h))
		})
	})
	return r
}

func handlers(h *Handler) http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", h.GetDocuments)
	})
	rg.Group(func(r chi.Router) {
		r.Post("/create", h.CreateDocument)
	})

	return rg
}

func (h *Handler) GetDocuments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res, err := h.DI.UseCase.Documents.GetDocuments(ctx)
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
