package document

type Document struct {
	Id             string
	Title          string
	OwnerId        string
	InitialContent string
	RoomId         string
	OrgId          string
	Content        []byte
}
