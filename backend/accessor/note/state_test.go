package note

import (
	"context"
	"testing"

	"stonks/accessor/db/pg/pgtest"
	"stonks/lib/lb"
)

type State struct {
	debugCtx context.Context
	storage  Storage
	h        *TestHelper
}

type It func(spec string, fn func(t *testing.T, s *State))

func Describe(describe string, t *testing.T, fn func(It It)) {
	t.Parallel()
	it := func(spec string, itFn func(t *testing.T, s *State)) {
		t.Run(spec, func(t *testing.T) {
			t.Parallel()
			pgtest.WithDb(t, func(db *pgtest.TestDb) {
				storage := NewPgStorage(
					lb.Randomized([]string{db.ConnStr()}),
					lb.Randomized([]string{db.ConnStr()}),
				)
				s := &State{
					debugCtx: db.DebugCtx(),
					storage:  storage,
					h:        NewTestHelper(t, db, storage),
				}
				itFn(t, s)
			})
		})
	}
	fn(it)
}
