package usecase

import (
	"context"
	"github.com/Tairascii/google-docs-documents/internal/app/service/document"
)

type DocumentsUseCase interface {
	CreateDocument(ctx context.Context, title, initialContent string) (string, error)
	GetDocuments(ctx context.Context) ([]document.Document, error)
	DeleteDocument(ctx context.Context, id string) error
	EditDocument(ctx context.Context, id, title string) error
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

func (u *UseCase) DeleteDocument(ctx context.Context, id string) error {
	_, err := u.documentsService.GetDocumentById(ctx, id)

	if err != nil {
		return err
	}

	return u.documentsService.DeleteDocument(ctx, id)
}

func (u *UseCase) EditDocument(ctx context.Context, id, title string) error {
	_, err := u.documentsService.GetDocumentById(ctx, id)

	if err != nil {
		return err
	}

	return u.documentsService.EditDocument(ctx, id, title)
}
