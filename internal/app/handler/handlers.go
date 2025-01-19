package handler

import (
	"github.com/Tairascii/google-docs-documents/internal/app"
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
	r.Mount("/document", handlers(h))
	return r
}

func handlers(h *Handler) http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
			h.CreateDocument(w, r)
		})
	})

	return rg
}

func (h *Handler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	err := h.DI.UseCase.Documents.CreateDocument()
	if err != nil {
		return
	}
}
