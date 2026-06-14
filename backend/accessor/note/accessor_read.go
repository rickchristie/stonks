package note

import (
	"stonks/accessor/db"
	"stonks/data"
)

type TransactionReader interface {
	db.Transaction

	GetNoteById(noteId int) (*data.Note, db.GetErrorType, error)
	ListActiveNotes() ([]*data.Note, error)
}
