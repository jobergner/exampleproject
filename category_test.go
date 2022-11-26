package main

import (
	"context"
	"exampleproject/db"
	"exampleproject/entity"
	"exampleproject/log"
	"exampleproject/repository"
	"exampleproject/repository/selector/quiz"
	"fmt"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Log(log.ReadingEnv, log.Err(err))
		panic(err)
	}
}

func TestCate(t *testing.T) {
	if err := db.Connect(); err != nil {
		panic(err)
	}

	if err := db.MigrateUp(); err != nil {
		panic(err)
	}

	var quizzes []entity.Quiz
	fmt.Println(repository.Default.Quiz.List(context.TODO(), &quizzes, quiz.OfArchivedCategorys()))
	fmt.Println(quizzes)
	t.FailNow()
}
