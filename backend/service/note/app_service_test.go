package note

import (
	"testing"

	accsNote "app-template/accessor/note"
	"app-template/data"
	"app-template/lib/tr"
	svcNote "app-template/service/note/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateNote(t *testing.T) {
	Describe("CreateNote", t, func(it It) {
		it("creates note after trimming title", func(t *testing.T, s *State) {
			out := s.service.CreateNote(s.debugCtx, svcNote.CreateNoteIn{
				Trace: &tr.Trace{TraceId: t.Name()},
				Title: "  Template note  ",
				Body:  "Body",
			})

			require.True(t, out.Success)
			assert.Equal(t, svcNote.AppErrNone, out.Error)
			require.NotNil(t, out.Note)
			assert.Equal(t, "Template note", out.Note.Title)
			assert.Equal(t, "Body", out.Note.Body)
			assert.Equal(t, data.NSActive, out.Note.Status)

			got := s.h.GetNote(out.Note.Id)
			assert.Equal(t, out.Note.Id, got.Id)
			assert.Equal(t, "Template note", got.Title)
			assert.Equal(t, "Body", got.Body)
			assert.Equal(t, data.NSActive, got.Status)
		})

		it("rejects blank title without writing", func(t *testing.T, s *State) {
			out := s.service.CreateNote(s.debugCtx, svcNote.CreateNoteIn{
				Trace: &tr.Trace{TraceId: t.Name()},
				Title: " ",
				Body:  "Body",
			})

			assert.False(t, out.Success)
			assert.Equal(t, svcNote.AppErrValidation, out.Error)
			assert.Nil(t, out.Note)

			list := s.service.ListNotes(s.debugCtx, svcNote.ListNotesIn{Trace: &tr.Trace{TraceId: t.Name()}})
			require.True(t, list.Success)
			assert.Empty(t, list.Notes)
		})
	})
}

func TestListNotes(t *testing.T) {
	Describe("ListNotes", t, func(it It) {
		it("lists active notes", func(t *testing.T, s *State) {
			s.h.NoteRows([]*accsNote.TestNoteRow{
				{Idx: 0, Title: "A"},
				{Idx: 1, Title: "Archived", Status: data.NSArchived},
				{Idx: 2, Title: "B"},
			})

			out := s.service.ListNotes(s.debugCtx, svcNote.ListNotesIn{Trace: &tr.Trace{TraceId: t.Name()}})

			require.True(t, out.Success)
			assert.Equal(t, svcNote.AppErrNone, out.Error)
			require.Len(t, out.Notes, 2)
			assert.Equal(t, "A", out.Notes[0].Title)
			assert.Equal(t, "B", out.Notes[1].Title)
		})
	})
}

func TestArchiveNote(t *testing.T) {
	Describe("ArchiveNote", t, func(it It) {
		it("archives active note and hides it from list", func(t *testing.T, s *State) {
			notes := s.h.NoteRows([]*accsNote.TestNoteRow{{Idx: 0, Title: "Active"}})

			out := s.service.ArchiveNote(s.debugCtx, svcNote.ArchiveNoteIn{
				Trace:  &tr.Trace{TraceId: t.Name()},
				NoteId: notes[0].Id,
			})

			require.True(t, out.Success)
			assert.Equal(t, svcNote.AppErrNone, out.Error)
			require.NotNil(t, out.Note)
			assert.Equal(t, data.NSArchived, out.Note.Status)

			list := s.service.ListNotes(s.debugCtx, svcNote.ListNotesIn{Trace: &tr.Trace{TraceId: t.Name()}})
			require.True(t, list.Success)
			assert.Empty(t, list.Notes)
		})
	})
}
