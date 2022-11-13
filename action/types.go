package action

type (
	NewQuiz struct {
		Category    string
		Title       string
		Description string
		Choices     []NewChoice
	}
	NewChoice struct {
		IsCorrect bool
		Content   string
	}
)
