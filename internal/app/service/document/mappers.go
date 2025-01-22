package document

import "github.com/Tairascii/google-docs-documents/internal/app/service/document/repo"

func toDocuments(raw []repo.Document) []Document {
	documents := make([]Document, len(raw))
	for i, doc := range raw {
		documents[i] = Document{
			Id:             doc.Id,
			Title:          doc.Title,
			OwnerId:        doc.OwnerId,
			InitialContent: doc.InitialContent,
			RoomId:         doc.RoomId,
			OrgId:          doc.OrgId,
		}
	}

	return documents
}
