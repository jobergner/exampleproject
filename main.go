package main

import (
	"exampleproject/db"
	"exampleproject/log"
	"exampleproject/serve"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println(os.Getenv("GOMOD"))
	err := godotenv.Load(".env")
	if err != nil {
		log.Log(log.ReadingEnv, log.Err(err))
		panic(err)
	}

	if err := db.Connect(); err != nil {
		panic(err)
	}

	if err := db.MigrateUp(); err != nil {
		panic(err)
	}

	signalShutdown, _ := serve.Listen()

	if err := <-signalShutdown; err != nil {
		panic(err)
	}

}
