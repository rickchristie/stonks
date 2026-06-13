package api

type CreateNoteReq struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type CreateNoteResp struct {
	Error ApiErr `json:"error"`
	Note  *Note  `json:"note"`
}

type ListNotesReq struct{}

type ListNotesResp struct {
	Error ApiErr  `json:"error"`
	Notes []*Note `json:"notes"`
}

type ArchiveNoteReq struct {
	NoteId int `json:"noteId"`
}

type ArchiveNoteResp struct {
	Error ApiErr `json:"error"`
	Note  *Note  `json:"note"`
}

type HelloReq struct{}

type HelloResp struct {
	Error           ApiErr `json:"error"`
	Message         string `json:"message"`
	DatabaseMessage string `json:"databaseMessage"`
	NoteCount       int    `json:"noteCount"`
}
