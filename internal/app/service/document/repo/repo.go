package repo

import "github.com/dancannon/gorethink"

type DocumentsRepo interface {
	CreateDocument() error
}

type Repo struct {
	session *gorethink.Session
}

func NewRepo(session *gorethink.Session) *Repo {
	return &Repo{
		session: session,
	}
}

func (r *Repo) CreateDocument() error {
	return nil
}
