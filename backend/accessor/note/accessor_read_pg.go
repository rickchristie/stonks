package note

import (
	"app-template/accessor/db"
	"app-template/accessor/db/pg"
	"app-template/data"
	"app-template/lib/mend"
)

var _ TransactionReader = (*pgReaderTx)(nil)

type pgReaderTx struct {
	*pg.ConnTxHelper
}

func (p *pgReaderTx) GetNoteById(noteId int) (*data.Note, db.GetErrorType, error) {
	query := `
		SELECT id, title, body, status, created_ts, last_updated_ts
		FROM note
		WHERE id = $1
	`

	row := p.Tx.QueryRow(p.Ctx, query, noteId)
	note := &data.Note{}
	if err := scanNote(row, note); err != nil {
		if pg.IsNotFound(err) {
			return nil, db.GetErrNotFound, db.ErrNotFound
		}
		return nil, db.GetErrUnknown, mend.Wrap(err, true)
	}

	return note, db.GetErrUnknown, nil
}

func (p *pgReaderTx) ListActiveNotes() ([]*data.Note, error) {
	query := `
		SELECT id, title, body, status, created_ts, last_updated_ts
		FROM note
		WHERE status = 'Active'
		ORDER BY id ASC
	`

	rows, err := p.Tx.Query(p.Ctx, query)
	if err != nil {
		return nil, mend.Wrap(err, true)
	}

	return pg.ConvertRowList(rows, scanNote)
}

func scanNote(row db.Scannable, note *data.Note) error {
	return row.Scan(
		&note.Id,
		&note.Title,
		&note.Body,
		&note.Status,
		&note.CreatedTs,
		&note.LastUpdatedTs,
	)
}
