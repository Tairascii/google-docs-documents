package document

import (
	"context"
	"errors"
	"github.com/Tairascii/google-docs-documents/internal/app/service/document/repo"
)

var (
	ErrInvalidOwnerId = errors.New("invalid owner id")
)

type DocumentsService interface {
	CreateDocument(ctx context.Context, title, initialContent string) (string, error)
	GetDocuments(ctx context.Context) ([]Document, error)
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

func (s *Service) GetDocuments(ctx context.Context) ([]Document, error) {
	ownerId, ok := ctx.Value("id").(string)
	if !ok {
		return nil, ErrInvalidOwnerId
	}
	raw, err := s.repo.GetDocuments(ctx, ownerId)
	if err != nil {
		return nil, err
	}
	return toDocuments(raw), nil
}
