package pgtest

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"stonks/accessor/db/pg"
	"stonks/accessor/db/pg/setup"
	"github.com/jackc/pgx/v5"
	"github.com/rickchristie/govner/pgflock/client"
)

const (
	pgflockPort     = 8001
	pgflockPassword = "pgflock"
)

type TestDb struct {
	t        *testing.T
	connStr  string
	debugCtx context.Context
}

func WithDb(t *testing.T, fn func(db *TestDb)) {
	t.Helper()
	db := &TestDb{t: t}
	db.initialize()
	defer db.cleanup()
	fn(db)
}

func (db *TestDb) initialize() {
	db.t.Helper()

	connStr, err := client.Lock(pgflockPort, db.t.Name(), pgflockPassword)
	if err != nil {
		if os.Getenv("STONKS_REQUIRE_PGFLOCK") == "1" {
			db.t.Fatalf("pgflock lock failed: %v", err)
		}
		if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "connect") {
			db.t.Skipf("pgflock is not running on port %d: %v", pgflockPort, err)
		}
		db.t.Fatalf("pgflock lock failed: %v", err)
	}
	db.connStr = connStr

	if err := db.runMigrations(); err != nil {
		db.cleanup()
		db.t.Fatalf("schema migration failed: %v", err)
	}

	db.debugCtx = pg.TimeoutOverrideForTesting(context.Background(), 30*time.Minute)
}

func (db *TestDb) cleanup() {
	if db.connStr != "" {
		_ = client.Unlock(pgflockPort, pgflockPassword, db.connStr)
		db.connStr = ""
	}
}

func (db *TestDb) ConnStr() string {
	return db.connStr
}

func (db *TestDb) DebugCtx() context.Context {
	return db.debugCtx
}

func (db *TestDb) Exec(sql string, args ...any) {
	db.t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, db.connStr)
	if err != nil {
		db.t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close(ctx)

	if _, err := conn.Exec(ctx, sql, args...); err != nil {
		db.t.Fatalf("exec failed: %v", err)
	}
}

func (db *TestDb) runMigrations() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, db.connStr)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	for _, sql := range setup.Init {
		if _, err := conn.Exec(ctx, sql.Value); err != nil {
			return err
		}
	}

	return nil
}
