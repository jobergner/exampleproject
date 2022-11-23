package api

import "exampleproject/entity"

type (
	NewQuiz struct {
		CategoryID  entity.CategoryID
		Title       string
		Description string
		Choices     []NewChoice
	}
	NewChoice struct {
		IsCorrect bool
		Content   string
	}
	LoginData struct {
		Name     string
		Password string
	}
)
