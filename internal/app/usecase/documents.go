package usecase

import (
	"context"
	"github.com/Tairascii/google-docs-documents/internal/app/service/document"
)

const (
	defaultTitle = "Untitled document"
)

type DocumentsUseCase interface {
	CreateDocument(ctx context.Context, title, initialContent string) (string, error)
	GetDocuments(ctx context.Context, search string) ([]document.Document, error)
	DeleteDocument(ctx context.Context, id string) error
	EditDocument(ctx context.Context, id, title string) error
	SaveDocumentContent(ctx context.Context, id string, content []byte) error
	CheckPermission(ctx context.Context, docID string) error
	WatchDocument(ctx context.Context, documentId string, ch chan<- []byte) error
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
		title = defaultTitle
	}
	return u.documentsService.CreateDocument(ctx, title, initialContent)
}

func (u *UseCase) GetDocuments(ctx context.Context, search string) ([]document.Document, error) {
	return u.documentsService.GetDocuments(ctx, search)
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

func (u *UseCase) SaveDocumentContent(ctx context.Context, id string, content []byte) error {
	return u.documentsService.SaveDocumentContent(ctx, id, content)
}

func (u *UseCase) CheckPermission(ctx context.Context, docID string) error {
	return u.documentsService.CheckPermission(ctx, docID)
}

func (u *UseCase) WatchDocument(ctx context.Context, documentId string, ch chan<- []byte) error {
	return u.documentsService.WatchDocument(ctx, documentId, ch)
}
