package document

import (
	"context"
	"errors"
	"github.com/Tairascii/google-docs-documents/internal/app/service/document/repo"
)

var (
	ErrInvalidOwnerId = errors.New("invalid owner id")
	ErrQueryDocument  = errors.New("error querying document")
	ErrNotAllowed     = errors.New("not allowed")
)

type DocumentsService interface {
	CreateDocument(ctx context.Context, title, initialContent string) (string, error)
	GetDocuments(ctx context.Context, search string) ([]Document, error)
	GetDocumentById(ctx context.Context, id string) (Document, error)
	DeleteDocument(ctx context.Context, id string) error
	EditDocument(ctx context.Context, id string, title string) error
	SaveDocumentContent(ctx context.Context, id string, content []byte) error
	CheckPermission(ctx context.Context, id string) error
}
type Service struct {
	repo repo.DocumentsRepo
}

func New(repo repo.DocumentsRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateDocument(ctx context.Context, title, initialContent string) (string, error) {
	ownerId, ok := ctx.Value("id").(string)
	if !ok {
		return "", ErrInvalidOwnerId
	}
	return s.repo.CreateDocument(ctx, title, initialContent, ownerId)
}

func (s *Service) GetDocuments(ctx context.Context, search string) ([]Document, error) {
	ownerId, ok := ctx.Value("id").(string)
	if !ok {
		return nil, ErrInvalidOwnerId
	}
	raw, err := s.repo.GetDocuments(ctx, ownerId, search)
	if err != nil {
		return nil, err
	}
	return toDocuments(raw), nil
}

func (s *Service) GetDocumentById(ctx context.Context, id string) (Document, error) {
	ownerId, ok := ctx.Value("id").(string)
	if !ok {
		return Document{}, ErrInvalidOwnerId
	}
	raw, err := s.repo.GetDocumentById(ctx, id)
	if err != nil {
		return Document{}, errors.Join(ErrQueryDocument, err)
	}

	if raw.OwnerId != ownerId {
		return Document{}, ErrNotAllowed
	}
	return Document(raw), nil
}

func (s *Service) DeleteDocument(ctx context.Context, id string) error {
	return s.repo.DeleteDocument(ctx, id)
}

func (s *Service) EditDocument(ctx context.Context, id string, title string) error {
	return s.repo.EditDocument(ctx, id, title)
}

func (s *Service) SaveDocumentContent(ctx context.Context, id string, content []byte) error {
	return s.repo.SaveDocumentContent(ctx, id, content)
}

func (s *Service) CheckPermission(ctx context.Context, id string) error {
	userID, ok := ctx.Value("id").(string)
	if !ok {
		return ErrInvalidOwnerId
	}
	doc, err := s.repo.GetDocumentById(ctx, id)
	if err != nil {
		return err
	}
	if doc.OwnerId != userID {
		return ErrNotAllowed
	}
	return nil
}
