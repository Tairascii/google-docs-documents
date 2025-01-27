package usecase

import (
	"context"
	"github.com/Tairascii/google-docs-documents/internal/app/service/document"
)

type DocumentsUseCase interface {
	CreateDocument(ctx context.Context, title, initialContent string) (string, error)
	GetDocuments(ctx context.Context) ([]document.Document, error)
}

type UseCase struct {
	documentsService document.DocumentsService
}

func NewDocumentsUseCase(documentsService document.DocumentsService) DocumentsUseCase {
	return &UseCase{
		documentsService: documentsService,
	}
}

func (u *UseCase) CreateDocument(ctx context.Context, title, initialContent string) (string, error) {
	if title == "" {
		title = "Untitled document"
	}
	return u.documentsService.CreateDocument(ctx, title, initialContent)
}

func (u *UseCase) GetDocuments(ctx context.Context) ([]document.Document, error) {
	return u.documentsService.GetDocuments(ctx)
}
