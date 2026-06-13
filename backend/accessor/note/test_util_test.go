package note

import (
	"testing"

	"app-template/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoteRows(t *testing.T) {
	Describe("NoteRows", t, func(it It) {
		it("returns inserted rows at fixture index", func(t *testing.T, s *State) {
			notes := s.h.NoteRows([]*TestNoteRow{
				{Idx: 2, Title: "Third", Body: "Third body"},
				{Idx: 0, Title: "First", Body: "First body"},
				{Idx: 1, Title: "Archived", Body: "Archived body", Status: data.NSArchived},
			})

			require.Len(t, notes, 3)
			assert.Equal(t, "First", notes[0].Title)
			assert.Equal(t, "Archived", notes[1].Title)
			assert.Equal(t, "Third", notes[2].Title)
			assert.Equal(t, data.NSActive, notes[0].Status)
			assert.Equal(t, data.NSArchived, notes[1].Status)
			assert.Equal(t, data.NSActive, notes[2].Status)
			assert.NotZero(t, notes[0].Id)
			assert.NotZero(t, notes[1].Id)
			assert.NotZero(t, notes[2].Id)

			got := s.h.GetNote(notes[1].Id)
			assert.Equal(t, "Archived body", got.Body)
			assert.Equal(t, data.NSArchived, got.Status)
		})
	})
}
