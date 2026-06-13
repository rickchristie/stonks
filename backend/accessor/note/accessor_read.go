package note

import (
	"app-template/accessor/db"
	"app-template/data"
)

type TransactionReader interface {
	db.Transaction

	GetNoteById(noteId int) (*data.Note, db.GetErrorType, error)
	ListActiveNotes() ([]*data.Note, error)
}
