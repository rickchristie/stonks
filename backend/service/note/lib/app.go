package lib

import (
	"context"

	"stonks/data"
	"stonks/lib/tr"
)

type AppClient interface {
	CreateNote(ctx context.Context, in CreateNoteIn) CreateNoteOut
	ListNotes(ctx context.Context, in ListNotesIn) ListNotesOut
	ArchiveNote(ctx context.Context, in ArchiveNoteIn) ArchiveNoteOut
}

type AppErr string

const (
	AppErrNone          AppErr = ""
	AppErrValidation    AppErr = "Validation"
	AppErrNotFound      AppErr = "NotFound"
	AppErrInternalError AppErr = "InternalError"
)

type CreateNoteIn struct {
	Trace *tr.Trace
	Title string
	Body  string
}

type CreateNoteOut struct {
	Success bool
	Error   AppErr
	Note    *data.Note
}

type ListNotesIn struct {
	Trace *tr.Trace
}

type ListNotesOut struct {
	Success bool
	Error   AppErr
	Notes   []*data.Note
}

type ArchiveNoteIn struct {
	Trace  *tr.Trace
	NoteId int
}

type ArchiveNoteOut struct {
	Success bool
	Error   AppErr
	Note    *data.Note
}
