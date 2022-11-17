package main

import (
	"exampleproject/db"
	"exampleproject/migration"
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(fmt.Errorf("error reading env file: %s", err))
	}

	if err := db.Connect(); err != nil {
		panic(err)
	}

	if err := migration.Migrate(); err != nil {
		panic(err)
	}
}
