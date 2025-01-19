package usecase

type DocumentsUseCase interface {
	CreateDocument() error
}

type UseCase struct {
}

func NewDocumentsUseCase() DocumentsUseCase {
	return &UseCase{}
}

func (u *UseCase) CreateDocument() error {
	return nil
}
