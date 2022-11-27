package main

import (
	"context"
	"exampleproject/db"
	"exampleproject/entity"
	"exampleproject/entity/selector/category"
	"exampleproject/entity/selector/quiz"
	"exampleproject/log"
	"exampleproject/repository"
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
	var cats []category.CategoryWithQuizCount
	fmt.Println(repository.Default.Category.List(context.TODO(), &cats, category.WithQuizCount(), category.IsNotArchived()))
	fmt.Println(cats)
	t.FailNow()
}
