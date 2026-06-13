package note

import (
	"app-template/accessor/db"
	"app-template/accessor/db/pg"
	"app-template/data"
	"app-template/lib/mend"
)

var _ TransactionWriter = (*pgWriterTx)(nil)

type pgWriterTx struct {
	*pg.ConnTxHelper
	*pgReaderTx
}

func (p *pgWriterTx) CreateNote(title string, body string) (*data.Note, error) {
	query := `
		INSERT INTO note (title, body, status)
		VALUES ($1, $2, 'Active')
		RETURNING id, title, body, status, created_ts, last_updated_ts
	`

	note := &data.Note{}
	if err := scanNote(p.Tx.QueryRow(p.Ctx, query, title, body), note); err != nil {
		return nil, mend.Wrap(err, true)
	}

	return note, nil
}

func (p *pgWriterTx) ArchiveNote(noteId int) (*data.Note, error) {
	query := `
		UPDATE note
		SET status = 'Archived',
		    last_updated_ts = NOW()
		WHERE id = $1
		  AND status = 'Active'
		RETURNING id, title, body, status, created_ts, last_updated_ts
	`

	note := &data.Note{}
	if err := scanNote(p.Tx.QueryRow(p.Ctx, query, noteId), note); err != nil {
		if pg.IsNotFound(err) {
			return nil, db.ErrNotFound
		}
		return nil, mend.Wrap(err, true)
	}

	return note, nil
}
