package document

import "github.com/Tairascii/google-docs-documents/internal/app/service/document/repo"

type DocumentsService interface {
	CreateDocument() error
}
type Service struct {
	repo repo.DocumentsRepo
}

func New(repo repo.DocumentsRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateDocument() error {
	return s.repo.CreateDocument()
}
