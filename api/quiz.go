package api

import (
	"context"
	"exampleproject/db"
	"exampleproject/entity"
	"exampleproject/entity/selector/choice"
	"exampleproject/log"
	"exampleproject/repository"
)

func IsChoiceCorrect(ctx context.Context, choiceID entity.ChoiceID) (bool, error) {
	var c entity.Choice
	if err := repository.Default.Choice.Get(ctx, &c, choice.IDEquals(choiceID)); err != nil {
		return false, err
	}

	return c.IsCorrect, nil
}

func CreateQuiz(ctx context.Context, newQuiz NewQuiz) error {
	tsx, err := db.DB.BeginTxx(ctx, nil)
	if err != nil {
		log.Log(log.BeginTSX, log.Err(err))
		return err
	}

	q := entity.Quiz{
		CategoryID:  newQuiz.CategoryID,
		Title:       newQuiz.Title,
		Description: newQuiz.Description,
	}

	quizID, err := repository.Default.Quiz.CreateTsx(ctx, tsx, &q)
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

		if _, err := repository.Default.Choice.CreateTsx(ctx, tsx, &choice); err != nil {
			tsx.Rollback()
			return err
		}
	}

	if err := tsx.Commit(); err != nil {
		tsx.Rollback()
		log.Log(log.CommitTSX, log.Err(err))
		return err
	}

	return nil
}
