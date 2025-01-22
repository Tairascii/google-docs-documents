package repo

import (
	"github.com/dancannon/gorethink"
	"log"
)

type DocumentsRepo interface {
	CreateDocument() error
	GetDocuments() ([]Document, error)
}

type Repo struct {
	session        *gorethink.Session
	documentsTable string
}

type Params struct {
	Session        *gorethink.Session
	DocumentsTable string
}

func NewRepo(params Params) *Repo {
	return &Repo{
		session:        params.Session,
		documentsTable: params.DocumentsTable,
	}
}

func (r *Repo) CreateDocument() error {
	return nil
}

func (r *Repo) GetDocuments() ([]Document, error) {
	cursor, err := gorethink.Table(r.documentsTable).Run(r.session)
	if err != nil {
		log.Fatalf("Error querying table: %v", err)
	}
	defer cursor.Close()

	result := make([]Document, 0)
	var doc Document
	for cursor.Next(&doc) {
		result = append(result, doc)
	}

	if err := cursor.Err(); err != nil {
		log.Fatalf("Error reading cursor: %v", err)
		return nil, err
	}
	return result, nil
}
