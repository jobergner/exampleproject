package category

import "exampleproject/entity"

type (
	CategoryWithQuizCount struct {
		ID        entity.CategoryID `db:"id"`
		Title     string            `db:"title"`
		Archived  bool              `db:"archived"`
		QuizCount int64             `db:"quiz_count"`
	}
)
