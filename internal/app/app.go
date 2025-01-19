package app

import "github.com/Tairascii/google-docs-documents/internal/app/usecase"

type UseCase struct {
	Documents usecase.DocumentsUseCase
}

type DI struct {
	UseCase UseCase
}
