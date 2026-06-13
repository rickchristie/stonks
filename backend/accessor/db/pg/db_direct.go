package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

	"app-template/accessor/db"
	"app-template/lib/lb"
	"app-template/lib/mend"
	"app-template/lib/tr"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var _ db.Transaction = (*ConnTxHelper)(nil)

type ConnTxHelper struct {
	Ctx    context.Context
	Tx     pgx.Tx
	Logger mend.Logger
	Trace  *tr.Trace
	cancel context.CancelFunc
	conn   *pgx.Conn
}

func NewConnTxHelper(
	ctx context.Context,
	trace *tr.Trace,
	conns lb.Selector[string],
	logger mend.Logger,
	readOnly bool,
) (*ConnTxHelper, error) {
	timeout := DefaultTxTimeout
	if override, ok := ctx.Value(timeoutKey).(time.Duration); ok {
		timeout = override
	}
	txCtx, cancel := context.WithTimeout(ctx, timeout)

	conn, err := pgx.Connect(txCtx, conns.Get())
	if err != nil {
		defer cancel()
		return nil, mend.Wrap(err, true)
	}

	if _, err = conn.Exec(txCtx, fmt.Sprintf(`SET idle_in_transaction_session_timeout = %v`, timeout.Milliseconds())); err != nil {
		defer cancel()
		defer closeConn(trace, logger, conn)
		return nil, mend.Wrap(err, true)
	}

	tx, err := conn.BeginTx(txCtx, StandardTxOptions(readOnly))
	if err != nil {
		defer cancel()
		defer closeConn(trace, logger, conn)
		return nil, mend.Wrap(err, true)
	}

	helper := &ConnTxHelper{
		Trace:  trace,
		Ctx:    txCtx,
		Tx:     tx,
		Logger: logger,
		conn:   conn,
		cancel: cancel,
	}

	if _, err = tx.Exec(txCtx, `SET DateStyle = 'ISO, YMD';`); err != nil {
		defer helper.Rollback()
		return nil, mend.Wrap(err, true)
	}

	return helper, nil
}

func closeConn(trace *tr.Trace, logger mend.Logger, conn *pgx.Conn) {
	if err := conn.Close(context.Background()); err != nil && err.Error() != "conn closed" {
		logger.FatalErr(trace, mend.Wrap(err, true)).Msg("failed to close connection")
	}
}

func (t *ConnTxHelper) close() {
	if t.conn != nil {
		closeConn(t.Trace, t.Logger, t.conn)
	}
}

func (t *ConnTxHelper) Commit() error {
	defer t.cancel()
	defer t.close()

	err := t.Tx.Commit(t.Ctx)
	if err != nil {
		if errors.Is(pgx.ErrTxClosed, err) {
			return nil
		}
		t.Logger.FatalErr(t.Trace, err).Msg("failed to commit transaction")
		return mend.Wrap(err, true)
	}

	return nil
}

func (t *ConnTxHelper) Rollback() {
	defer t.cancel()
	defer t.close()

	rollbackCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := t.Tx.Rollback(rollbackCtx)
	if err != nil {
		if errors.Is(err, pgx.ErrTxClosed) || err.Error() == "conn closed" {
			return
		}
		if err.Error() == "failed to deallocate cached statement(s): conn closed" {
			return
		}
		t.Logger.FatalErr(t.Trace, mend.Wrap(err, true)).Msg("failed to rollback transaction")
	}
}

func (t *ConnTxHelper) AffectedOneRow(res pgconn.CommandTag) error {
	return t.AffectedRows(res, 1)
}

func (t *ConnTxHelper) AffectedRows(res pgconn.CommandTag, expected int64) error {
	n := res.RowsAffected()
	if n != expected {
		return mend.Err("unexpected num of rows affected", true).
			Str("expected", expected).
			Str("rowsAffected", n)
	}
	return nil
}
