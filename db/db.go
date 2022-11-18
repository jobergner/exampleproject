package db

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

var (
	ErrNotFound       = errors.New("not found")
	ErrTooManyResults = errors.New("too many results")
)
