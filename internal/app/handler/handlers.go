package handler

import (
	"github.com/Tairascii/google-docs-documents/internal/app"
	"github.com/Tairascii/google-docs-documents/internal/app/service/document"
	"github.com/Tairascii/google-docs-documents/pkg"
	"github.com/go-chi/chi"
	"net/http"
)

type Handler struct {
	DI *app.DI
}

func NewHandler(di *app.DI) *Handler {
	return &Handler{DI: di}
}

func (h *Handler) InitHandlers() *chi.Mux {
	r := chi.NewRouter()
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
	res, err := h.DI.UseCase.Documents.GetDocuments()
	if err != nil {
		pkg.JSONErrorResponseWriter(w, err, http.StatusInternalServerError)
		return
	}
	pkg.JSONResponseWriter[[]document.Document](w, res, http.StatusOK)
}

func (h *Handler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	err := h.DI.UseCase.Documents.CreateDocument()
	if err != nil {
		return
	}
}
