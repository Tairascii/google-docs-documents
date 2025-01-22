package usecase

import "github.com/Tairascii/google-docs-documents/internal/app/service/document"

type DocumentsUseCase interface {
	CreateDocument() error
	GetDocuments() ([]document.Document, error)
}

type UseCase struct {
	documentsService document.DocumentsService
}

func NewDocumentsUseCase(documentsService document.DocumentsService) DocumentsUseCase {
	return &UseCase{
		documentsService: documentsService,
	}
}

func (u *UseCase) CreateDocument() error {
	return u.documentsService.CreateDocument()
}

func (u *UseCase) GetDocuments() ([]document.Document, error) {
	return u.documentsService.GetDocuments()
}
