package action

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
)
