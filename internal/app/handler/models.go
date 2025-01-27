package handler

type CreateDocumentPayload struct {
	Title          string `json:"title,omitempty"`
	InitialContent string `json:"initial_content,omitempty"`
}

type CreateDocumentResponse struct {
	DocumentID string `json:"document_id"`
}

type Document struct {
	Id             string `json:"id"`
	Title          string `json:"title"`
	OwnerId        string `json:"owner_id"`
	InitialContent string `json:"initial_content"`
	RoomId         string `json:"room_id"`
	OrgId          string `json:"org_id"`
}
