package note

import "app-template/data"

type TransactionWriter interface {
	TransactionReader

	CreateNote(title string, body string) (*data.Note, error)
	ArchiveNote(noteId int) (*data.Note, error)
}
