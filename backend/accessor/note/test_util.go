package note

import (
	"context"
	"testing"

	"app-template/accessor/db"
	"app-template/accessor/db/pg/pgtest"
	"app-template/data"
	"app-template/lib/tr"
	"github.com/stretchr/testify/require"
)

type TestHelper struct {
	Storage Storage

	t        *testing.T
	db       *pgtest.TestDb
	debugCtx context.Context
}

type TestNoteRow struct {
	Idx    int
	Title  string
	Body   string
	Status data.NoteStatus
}

func NewTestHelper(t *testing.T, db *pgtest.TestDb, storage Storage) *TestHelper {
	return &TestHelper{
		Storage:  storage,
		t:        t,
		db:       db,
		debugCtx: db.DebugCtx(),
	}
}

func (h *TestHelper) Accessor() TransactionWriter {
	tx, err := h.Storage.BeginTx(h.debugCtx, &tr.Trace{TraceId: h.t.Name()})
	require.NoError(h.t, err)
	return tx
}

func (h *TestHelper) NoteRows(rows []*TestNoteRow) []*data.Note {
	h.t.Helper()
	ret := make([]*data.Note, noteRowsLen(rows))

	tx := h.Accessor()
	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	for _, row := range rows {
		require.GreaterOrEqual(h.t, row.Idx, 0)

		title := row.Title
		if title == "" {
			title = "Note " + string(rune('A'+row.Idx))
		}
		status := row.Status
		if status == "" {
			status = data.NSActive
		}

		note, err := tx.CreateNote(title, row.Body)
		require.NoError(h.t, err)
		switch status {
		case data.NSActive:
		case data.NSArchived:
			note, err = tx.ArchiveNote(note.Id)
			require.NoError(h.t, err)
		default:
			require.FailNowf(h.t, "unsupported note fixture status", "status=%s", status)
		}

		// Callers address fixture rows by Idx, including archived rows that
		// ListActiveNotes intentionally hides.
		require.Nil(h.t, ret[row.Idx])
		ret[row.Idx] = note
	}

	require.NoError(h.t, tx.Commit())
	committed = true
	return ret
}

func noteRowsLen(rows []*TestNoteRow) int {
	maxIdx := -1
	for _, row := range rows {
		if row.Idx > maxIdx {
			maxIdx = row.Idx
		}
	}
	return maxIdx + 1
}

func (h *TestHelper) GetNote(noteId int) *data.Note {
	h.t.Helper()
	tx := h.Accessor()
	defer tx.Rollback()

	note, errType, err := tx.GetNoteById(noteId)
	require.NoError(h.t, err)
	require.Equal(h.t, db.GetErrUnknown, errType)
	return note
}
