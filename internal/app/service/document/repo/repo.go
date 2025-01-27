package repo

import (
	"context"
	"errors"
	"github.com/dancannon/gorethink"
	"log"
)

var (
	ErrCreateDocument = errors.New("error creating document")
)

const (
	ownerIdField = "owner_id"
)

type DocumentsRepo interface {
	CreateDocument(ctx context.Context, title, initialContent, ownerId string) (string, error)
	GetDocuments(ctx context.Context, ownerId string) ([]Document, error)
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

func (r *Repo) CreateDocument(ctx context.Context, title, initialContent, ownerId string) (string, error) {
	doc := Document{
		Title:          title,
		OwnerId:        ownerId,
		InitialContent: initialContent,
	}
	res, err := gorethink.Table(r.documentsTable).Insert(doc).RunWrite(r.session, gorethink.RunOpts{Context: ctx})
	if err != nil {
		return "", errors.Join(ErrCreateDocument, err)
	}

	if len(res.GeneratedKeys) != 1 {
		return "", ErrCreateDocument
	}

	return res.GeneratedKeys[0], nil
}

func (r *Repo) GetDocuments(ctx context.Context, ownerId string) ([]Document, error) {
	filterByOwner := gorethink.Row.Field(ownerIdField).Eq(ownerId)
	cursor, err := gorethink.Table(r.documentsTable).Filter(filterByOwner).Run(r.session, gorethink.RunOpts{Context: ctx})
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
