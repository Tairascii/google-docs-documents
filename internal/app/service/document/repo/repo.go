package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dancannon/gorethink"
	"log"
)

var (
	ErrCreateDocument = errors.New("error creating document")
	ErrNoData         = errors.New("no data")
)

const (
	ownerIdField    = "owner_id"
	titleField      = "title"
	documentIdField = "document_id"
)

type DocumentsRepo interface {
	CreateDocument(ctx context.Context, title, initialContent, ownerId string) (string, error)
	GetDocuments(ctx context.Context, ownerId, search string) ([]Document, error)
	GetDocumentById(ctx context.Context, documentId string) (Document, error)
	DeleteDocument(ctx context.Context, documentId string) error
	EditDocument(ctx context.Context, documentId string, title string) error
	SaveDocumentContent(ctx context.Context, documentId string, content []byte) error
	WatchTableChange()
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

func (r *Repo) GetDocuments(ctx context.Context, ownerId, search string) ([]Document, error) {
	filterByOwner := gorethink.Row.Field(ownerIdField).Eq(ownerId)
	filterByName := gorethink.Row.Field(titleField).Match(fmt.Sprintf("(?i)%s", search))
	cursor, err := gorethink.Table(r.documentsTable).Filter(filterByOwner).Filter(filterByName).Run(r.session, gorethink.RunOpts{Context: ctx})
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

func (r *Repo) GetDocumentById(ctx context.Context, documentId string) (Document, error) {
	var doc Document
	err := gorethink.Table(r.documentsTable).Get(documentId).ReadOne(&doc, r.session, gorethink.RunOpts{Context: ctx})
	if err != nil {
		return Document{}, err
	}

	return doc, nil
}

func (r *Repo) DeleteDocument(ctx context.Context, documentId string) error {
	doc := gorethink.Table(r.documentsTable).Get(documentId)
	return doc.Delete().Exec(r.session, gorethink.ExecOpts{Context: ctx, NoReply: true})
}

func (r *Repo) EditDocument(ctx context.Context, documentId string, title string) error {
	doc := gorethink.Table(r.documentsTable).Get(documentId)
	return doc.Update(Document{Title: title}).Exec(r.session, gorethink.ExecOpts{Context: ctx, NoReply: true})
}

func (r *Repo) SaveDocumentContent(ctx context.Context, documentId string, content []byte) error {
	doc := gorethink.Table(r.documentsTable)
	_, err := doc.Insert(Document{Id: documentId, Content: content}, gorethink.InsertOpts{Conflict: "replace"}).RunWrite(r.session)
	return err
}

func (r *Repo) WatchTableChange() {
	cursor, err := gorethink.Table(r.documentsTable).Changes().Run(r.session)
	if err != nil {
		log.Println("Error watching changes:", err)
		return
	}
	defer cursor.Close()

	var change map[string]interface{}
	for cursor.Next(&change) {
		fmt.Println("Document updated:", change)

		newValue := change["new_val"].(map[string]interface{})
		docID := newValue["id"].(string)
		content, _ := json.Marshal(newValue["content"])

		fmt.Printf("new data: %v %v \n", docID, string(content))
	}
}
