package repo

type Document struct {
	Id             string `gorethink:"id,omitempty"`
	Title          string `gorethink:"title"`
	OwnerId        string `gorethink:"owner_id"`
	InitialContent string `gorethink:"initial_content,omitempty"`
	RoomId         string `gorethink:"room_id,omitempty"`
	OrgId          string `gorethink:"org_id,omitempty"`
}
