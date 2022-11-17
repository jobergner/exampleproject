package migration

import (
	"errors"
	"exampleproject/db"
	"exampleproject/log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// TODO: logging
func Migrate() error {
	driver, err := postgres.WithInstance(db.DB.DB, &postgres.Config{})
	if err != nil {
		log.Log(log.Migration, log.Err(err))
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migration/",
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
