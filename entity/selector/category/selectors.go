package category

import (
	"exampleproject/entity"
	"exampleproject/entity/selector"

	"github.com/Masterminds/squirrel"
)

func IDEquals(id entity.CategoryID) selector.Selector {
	return selector.Selector{Where: squirrel.Expr("id = ?", id), Name: "IDEquals"}
}

func TitleEquals(title string) selector.Selector {
	return selector.Selector{Where: squirrel.Expr("title = ?", title), Name: "TitleEquals"}
}

func IsNotArchived() selector.Selector {
	return selector.Selector{Where: squirrel.Expr("archived = FALSE"), Name: "IsNotArchived"}
}

func WithQuizCount() selector.Selector {
	return selector.Selector{
		Columns: []string{"(SELECT COUNT(*) FROM quizzes WHERE category_id = categories.id) as quiz_count"},
		Name:    "WithQuizCount",
	}
}
