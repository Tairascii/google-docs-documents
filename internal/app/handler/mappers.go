package handler

import "github.com/Tairascii/google-docs-documents/internal/app/service/document"

func toDocuments(raw []document.Document) []Document {
	documents := make([]Document, len(raw))
	for i, doc := range raw {
		documents[i] = Document(doc)
	}
	return documents
}
