package note

import (
	"testing"

	"stonks/data"
	"stonks/lib/tr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateNote(t *testing.T) {
	Describe("CreateNote", t, func(it It) {
		it("creates note with title and body", func(t *testing.T, s *State) {
			tx, err := s.storage.BeginTx(s.debugCtx, &tr.Trace{TraceId: t.Name()})
			require.NoError(t, err)

			note, err := tx.CreateNote("Created", "Created body")
			require.NoError(t, err)
			require.NoError(t, tx.Commit())

			got := s.h.GetNote(note.Id)
			assert.Equal(t, note.Id, got.Id)
			assert.Equal(t, "Created", got.Title)
			assert.Equal(t, "Created body", got.Body)
			assert.Equal(t, data.NSActive, got.Status)
			assert.False(t, got.CreatedTs.IsZero())
			assert.False(t, got.LastUpdatedTs.IsZero())
		})
	})
}

func TestArchiveNote(t *testing.T) {
	Describe("ArchiveNote", t, func(it It) {
		it("archives active note", func(t *testing.T, s *State) {
			notes := s.h.NoteRows([]*TestNoteRow{{Idx: 0, Title: "Active"}})

			tx, err := s.storage.BeginTx(s.debugCtx, &tr.Trace{TraceId: t.Name()})
			require.NoError(t, err)
			archived, err := tx.ArchiveNote(notes[0].Id)
			require.NoError(t, err)
			require.NoError(t, tx.Commit())

			got := s.h.GetNote(notes[0].Id)
			assert.Equal(t, archived.Id, got.Id)
			assert.Equal(t, data.NSArchived, got.Status)
			assert.Equal(t, "Active", got.Title)
		})
	})
}
