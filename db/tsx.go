package db

import (
	"context"
	"database/sql"
	"exampleproject/log"

	"github.com/jmoiron/sqlx"
)

// TSX is an optional transaction which resorts to DB when using NoTSX
type TSX struct {
	tsx *sqlx.Tx
}

var NoTSX = new(TSX)

type Executer interface {
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func (t *TSX) Get() Executer {
	if t.tsx == nil {
		return DB
	}

	return t.tsx
}

func BeginTSX(ctx context.Context) (*TSX, error) {
	tsx, err := DB.BeginTxx(ctx, nil)
	if err != nil {
		log.Log(log.BeginTSX, log.Err(err))
		return nil, err
	}

	t := TSX{
		tsx: tsx,
	}

	return &t, nil
}

func (t *TSX) Rollback() error {
	if err := t.tsx.Rollback(); err != nil {
		log.Log(log.RollbackTSX, log.Err(err))
		return err
	}

	return nil
}

func (t *TSX) Commit() error {
	if err := t.tsx.Commit(); err != nil {
		_ = t.Rollback()
		log.Log(log.CommitTSX, log.Err(err))
		return err
	}

	return nil
}
