package db

import "errors"

// Transaction is embedded by domain accessor transaction interfaces.
//
// Keep transaction lifetimes short. Use row locks in a consistent order when
// multiple records must be locked; if lock acquisition fails, rollback and
// retry the whole service workflow.
type Transaction interface {
	Commit() error
	Rollback()
}

type Scannable interface {
	Scan(dest ...interface{}) error
}

type GetErrorType int

const (
	GetErrUnknown  GetErrorType = 0
	GetErrNotFound GetErrorType = 1
)

var ErrNotFound = errors.New("not found")

type InsertErrorType int

const (
	InsertErrUnknown      InsertErrorType = 0
	InsertErrDuplicateKey InsertErrorType = 1
)
