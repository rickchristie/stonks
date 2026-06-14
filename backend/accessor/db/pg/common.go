package pg

import (
	"errors"

	"stonks/accessor/db"
	"stonks/lib/mend"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func IsNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func IsDuplicateKeyErr(err error) (bool, string) {
	pgErr := &pgconn.PgError{}
	if !errors.As(err, &pgErr) {
		return false, ""
	}
	if pgErr.Code != "23505" {
		return false, ""
	}
	return true, pgErr.ConstraintName
}

func ConvertRowList[T any](rows pgx.Rows, convert func(db.Scannable, *T) error) ([]*T, error) {
	defer rows.Close()

	ret := make([]*T, 0)
	for rows.Next() {
		dat := new(T)
		if err := convert(rows, dat); err != nil {
			return nil, err
		}
		ret = append(ret, dat)
	}

	if err := rows.Err(); err != nil {
		return nil, mend.Wrap(err, true)
	}

	return ret, nil
}
