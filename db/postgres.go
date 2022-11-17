package db

import (
	"exampleproject/log"
	"fmt"
	"os"

	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
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

func MigrateUp() error {
	driver, err := postgres.WithInstance(DB.DB, &postgres.Config{})
	if err != nil {
		log.Log(log.Migration, log.Err(err))
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations/",
		"postgres",
		driver,
	)
	if err != nil {
		log.Log(log.Migration, log.Err(err))
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Log(log.Migration, log.Err(err))
		return err
	}

	return nil
}
