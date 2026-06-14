package note

import (
	"context"
	"testing"

	"stonks/accessor/db/pg/pgtest"
	accsNote "stonks/accessor/note"
	"stonks/lib/lb"
	svcNote "stonks/service/note/lib"
)

type State struct {
	debugCtx context.Context
	storage  accsNote.Storage
	service  svcNote.AppClient
	h        *accsNote.TestHelper
}

type It func(spec string, fn func(t *testing.T, s *State))

func Describe(describe string, t *testing.T, fn func(It It)) {
	t.Parallel()
	it := func(spec string, itFn func(t *testing.T, s *State)) {
		t.Run(spec, func(t *testing.T) {
			t.Parallel()
			pgtest.WithDb(t, func(db *pgtest.TestDb) {
				storage := accsNote.NewPgStorage(
					lb.Randomized([]string{db.ConnStr()}),
					lb.Randomized([]string{db.ConnStr()}),
				)
				s := &State{
					debugCtx: db.DebugCtx(),
					storage:  storage,
					service:  NewAppService(storage),
					h:        accsNote.NewTestHelper(t, db, storage),
				}
				itFn(t, s)
			})
		})
	}
	fn(it)
}
