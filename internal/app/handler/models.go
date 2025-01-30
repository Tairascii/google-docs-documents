package handler

type CreateDocumentPayload struct {
	Title          string `json:"title,omitempty"`
	InitialContent string `json:"initialContent,omitempty"`
}

type CreateDocumentResponse struct {
	DocumentID string `json:"documentId"`
}

type Document struct {
	Id             string `json:"id"`
	Title          string `json:"title"`
	OwnerId        string `json:"ownerId"`
	InitialContent string `json:"initialContent"`
	RoomId         string `json:"roomId"`
	OrgId          string `json:"orgId"`
}
