package api

import (
	"context"
	"exampleproject/db"
	"exampleproject/entity"
	"exampleproject/repository"
	"exampleproject/repository/selector"
)

func IsChoiceCorrect(ctx context.Context, choiceID entity.ChoiceID) (bool, error) {
	var choice entity.Choice
	if err := repository.Default.Choice.Get(ctx, &choice, selector.IDEquals(int(choiceID))); err != nil {
		return false, err
	}

	return choice.IsCorrect, nil
}

func CreateQuiz(ctx context.Context, newQuiz NewQuiz) error {
	tsx, err := db.BeginTSX(ctx)
	if err != nil {
		return err
	}

	q := entity.Quiz{
		CategoryID:  newQuiz.CategoryID,
		Title:       newQuiz.Title,
		Description: newQuiz.Description,
	}

	quizID, err := repository.Default.Quiz.Create(ctx, tsx, &q)
	if err != nil {
		tsx.Rollback()
		return err
	}

	for _, newChoice := range newQuiz.Choices {
		choice := entity.Choice{
			QuizID:    entity.QuizID(quizID),
			IsCorrect: newChoice.IsCorrect,
			Content:   newChoice.Content,
		}

		if _, err := repository.Default.Choice.Create(ctx, tsx, &choice); err != nil {
			tsx.Rollback()
			return err
		}
	}

	if err := tsx.Commit(); err != nil {
		tsx.Rollback()
		return err
	}

	return nil
}
