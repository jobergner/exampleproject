package main

import (
	"exampleproject/db"
	"exampleproject/log"
	"exampleproject/server"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Log(log.ReadingEnv, log.Err(err))
		panic(err)
	}
}

func main() {
	if err := db.Connect(); err != nil {
		panic(err)
	}

	if err := db.MigrateUp(); err != nil {
		panic(err)
	}

	signalShutdown, _ := server.Listen()

	if err := <-signalShutdown; err != nil {
		panic(err)
	}
}
