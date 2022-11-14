package db

import (
	"exampleproject/log"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

func Connect() error {
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable connect_timeout=5",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Log(log.DBConnection, log.Err(err))
		return err
	}

	if err = db.Ping(); err != nil {
		log.Log(log.DBConnection, log.Err(err))
		return err
	}

	DB = db

	return nil
}
