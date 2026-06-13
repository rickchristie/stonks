package note

import (
	"context"

	"app-template/accessor/db/pg"
	"app-template/lib/lb"
	"app-template/lib/mend"
	"app-template/lib/tr"
)

type Storage interface {
	BeginTxReader(ctx context.Context, trace *tr.Trace) (TransactionReader, error)
	BeginTx(ctx context.Context, trace *tr.Trace) (TransactionWriter, error)
}

var _ Storage = (*pgStorage)(nil)

type pgStorage struct {
	logger        mend.Logger
	readerConnStr lb.Selector[string]
	writerConnStr lb.Selector[string]
}

func NewPgStorage(readerConnStr lb.Selector[string], writerConnStr lb.Selector[string]) Storage {
	return &pgStorage{
		logger:        mend.NewZerologLogger("accsNote"),
		readerConnStr: readerConnStr,
		writerConnStr: writerConnStr,
	}
}

func (p *pgStorage) BeginTxReader(ctx context.Context, trace *tr.Trace) (TransactionReader, error) {
	tx, err := pg.NewConnTxHelper(ctx, trace, p.readerConnStr, p.logger, true)
	if err != nil {
		return nil, err
	}
	return &pgReaderTx{ConnTxHelper: tx}, nil
}

func (p *pgStorage) BeginTx(ctx context.Context, trace *tr.Trace) (TransactionWriter, error) {
	tx, err := pg.NewConnTxHelper(ctx, trace, p.writerConnStr, p.logger, false)
	if err != nil {
		return nil, err
	}
	return &pgWriterTx{
		ConnTxHelper: tx,
		pgReaderTx:   &pgReaderTx{ConnTxHelper: tx},
	}, nil
}
