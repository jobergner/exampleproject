package db

import (
	"database/sql"
	"exampleproject/log"
	"fmt"
	"os"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres() *Postgres {
	return &Postgres{}
}

func (p *Postgres) Connect() error {
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable connect_timeout=5",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Log(log.DBConnection, log.Err(err))
		return err
	}

	if err = db.Ping(); err != nil {
		log.Log(log.DBConnection, log.Err(err))
		return err
	}

	p.db = db

	return nil
}
