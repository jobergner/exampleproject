package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

var (
	ErrNotFound       = errors.New("not found")
	ErrTooManyResults = errors.New("too many results")
)

// TSX is an optional transaction which resorts to DB when using NoTSX
type TSX struct {
	tsx *sqlx.Tx
}

var NoTSX = new(TSX)

type Executer interface {
	sqlx.ExtContext
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

func (t *TSX) Get() Executer {
	if t.tsx == nil {
		return DB
	}

	return t.tsx
}

func (t *TSX) Begin(ctx context.Context) (sqlx.ExecerContext, error) {
	tsx, err := DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	t.tsx = tsx

	return tsx, nil
}

func (t *TSX) Rollback() error {
	return t.tsx.Rollback()
}

func (t *TSX) Commit() error {
	if err := t.tsx.Commit(); err != nil {
		_ = t.Rollback()
		return err
	}

	return nil
}
