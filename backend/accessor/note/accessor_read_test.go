package note

import (
	"errors"
	"testing"

	"app-template/accessor/db"
	"app-template/data"
	"app-template/lib/tr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetNoteById(t *testing.T) {
	Describe("GetNoteById", t, func(it It) {
		it("returns full note fields", func(t *testing.T, s *State) {
			notes := s.h.NoteRows([]*TestNoteRow{
				{Idx: 0, Title: "First", Body: "Readable body"},
			})

			tx, err := s.storage.BeginTxReader(s.debugCtx, &tr.Trace{TraceId: t.Name()})
			require.NoError(t, err)
			defer tx.Rollback()

			got, errType, err := tx.GetNoteById(notes[0].Id)
			require.NoError(t, err)
			assert.Equal(t, db.GetErrUnknown, errType)
			assert.Equal(t, notes[0].Id, got.Id)
			assert.Equal(t, "First", got.Title)
			assert.Equal(t, "Readable body", got.Body)
			assert.Equal(t, data.NSActive, got.Status)
			assert.False(t, got.CreatedTs.IsZero())
			assert.False(t, got.LastUpdatedTs.IsZero())
		})

		it("returns not found error for missing note", func(t *testing.T, s *State) {
			tx, err := s.storage.BeginTxReader(s.debugCtx, &tr.Trace{TraceId: t.Name()})
			require.NoError(t, err)
			defer tx.Rollback()

			got, errType, err := tx.GetNoteById(9999)
			assert.True(t, errors.Is(err, db.ErrNotFound))
			assert.Equal(t, db.GetErrNotFound, errType)
			assert.Nil(t, got)
		})
	})
}

func TestListActiveNotes(t *testing.T) {
	Describe("ListActiveNotes", t, func(it It) {
		it("returns active notes ordered by id", func(t *testing.T, s *State) {
			s.h.NoteRows([]*TestNoteRow{
				{Idx: 0, Title: "A", Status: data.NSActive},
				{Idx: 1, Title: "Archived", Status: data.NSArchived},
				{Idx: 2, Title: "B", Status: data.NSActive},
			})

			tx, err := s.storage.BeginTxReader(s.debugCtx, &tr.Trace{TraceId: t.Name()})
			require.NoError(t, err)
			defer tx.Rollback()

			got, err := tx.ListActiveNotes()
			require.NoError(t, err)
			require.Len(t, got, 2)
			assert.Equal(t, "A", got[0].Title)
			assert.Equal(t, "B", got[1].Title)
			assert.Equal(t, data.NSActive, got[0].Status)
			assert.Equal(t, data.NSActive, got[1].Status)
		})
	})
}
