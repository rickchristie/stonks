package note

import (
	"context"
	"errors"
	"strings"

	"stonks/accessor/db"
	accsNote "stonks/accessor/note"
	"stonks/lib/mend"
	svcNote "stonks/service/note/lib"
)

var logger = mend.NewZerologLogger("svcNote")
var _ svcNote.AppClient = (*appService)(nil)

type appService struct {
	storage accsNote.Storage
}

func NewAppService(storage accsNote.Storage) svcNote.AppClient {
	return &appService{storage: storage}
}

func (s *appService) CreateNote(ctx context.Context, in svcNote.CreateNoteIn) svcNote.CreateNoteOut {
	title := strings.TrimSpace(in.Title)
	if title == "" || len(title) > 200 || len(in.Body) > 20000 {
		return svcNote.CreateNoteOut{Error: svcNote.AppErrValidation}
	}

	tx, err := s.storage.BeginTx(ctx, in.Trace)
	if err != nil {
		logger.ErrorErr(in.Trace, err).Msg("failed to begin tx")
		return svcNote.CreateNoteOut{Error: svcNote.AppErrInternalError}
	}
	defer tx.Rollback()

	note, err := tx.CreateNote(title, in.Body)
	if err != nil {
		logger.ErrorErr(in.Trace, err).Msg("failed to create note")
		return svcNote.CreateNoteOut{Error: svcNote.AppErrInternalError}
	}

	if err := tx.Commit(); err != nil {
		logger.ErrorErr(in.Trace, err).Msg("failed to commit")
		return svcNote.CreateNoteOut{Error: svcNote.AppErrInternalError}
	}

	return svcNote.CreateNoteOut{
		Success: true,
		Error:   svcNote.AppErrNone,
		Note:    note,
	}
}

func (s *appService) ListNotes(ctx context.Context, in svcNote.ListNotesIn) svcNote.ListNotesOut {
	tx, err := s.storage.BeginTxReader(ctx, in.Trace)
	if err != nil {
		logger.ErrorErr(in.Trace, err).Msg("failed to begin read tx")
		return svcNote.ListNotesOut{Error: svcNote.AppErrInternalError}
	}
	defer tx.Rollback()

	notes, err := tx.ListActiveNotes()
	if err != nil {
		logger.ErrorErr(in.Trace, err).Msg("failed to list active notes")
		return svcNote.ListNotesOut{Error: svcNote.AppErrInternalError}
	}

	return svcNote.ListNotesOut{
		Success: true,
		Error:   svcNote.AppErrNone,
		Notes:   notes,
	}
}

func (s *appService) ArchiveNote(ctx context.Context, in svcNote.ArchiveNoteIn) svcNote.ArchiveNoteOut {
	if in.NoteId <= 0 {
		return svcNote.ArchiveNoteOut{Error: svcNote.AppErrValidation}
	}

	tx, err := s.storage.BeginTx(ctx, in.Trace)
	if err != nil {
		logger.ErrorErr(in.Trace, err).Msg("failed to begin tx")
		return svcNote.ArchiveNoteOut{Error: svcNote.AppErrInternalError}
	}
	defer tx.Rollback()

	note, err := tx.ArchiveNote(in.NoteId)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return svcNote.ArchiveNoteOut{Error: svcNote.AppErrNotFound}
		}
		logger.ErrorErr(in.Trace, err).Int("noteId", in.NoteId).Msg("failed to archive note")
		return svcNote.ArchiveNoteOut{Error: svcNote.AppErrInternalError}
	}

	if err := tx.Commit(); err != nil {
		logger.ErrorErr(in.Trace, err).Int("noteId", in.NoteId).Msg("failed to commit")
		return svcNote.ArchiveNoteOut{Error: svcNote.AppErrInternalError}
	}

	return svcNote.ArchiveNoteOut{
		Success: true,
		Error:   svcNote.AppErrNone,
		Note:    note,
	}
}
