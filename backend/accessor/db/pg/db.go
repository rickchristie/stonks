package pg

import (
	"time"

	"github.com/jackc/pgx/v5"
)

type dbContextKey string

var timeoutKey dbContextKey = "_to"

const DefaultTxTimeout = 10 * time.Second

func StandardTxOptions(readOnly bool) pgx.TxOptions {
	opts := pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		DeferrableMode: pgx.NotDeferrable,
	}

	if readOnly {
		opts.AccessMode = pgx.ReadOnly
	} else {
		opts.AccessMode = pgx.ReadWrite
	}

	return opts
}
