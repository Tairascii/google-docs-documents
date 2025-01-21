package usecase

import "github.com/Tairascii/google-docs-documents/internal/app/service/document"

type DocumentsUseCase interface {
	CreateDocument() error
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
