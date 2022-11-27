package quiz

import (
	"exampleproject/entity"
	"exampleproject/entity/selector"

	"github.com/Masterminds/squirrel"
)

func IDEquals(id entity.QuizID) selector.Selector {
	return selector.Selector{Where: squirrel.Expr("id = ?", id), Name: "IDEquals"}
}

func TitleEquals(title string) selector.Selector {
	return selector.Selector{Where: squirrel.Expr("title = ?", title), Name: "TitleEquals"}
}

func IsNotArchived() selector.Selector {
	return selector.Selector{Where: squirrel.Expr("archived = FALSE"), Name: "IsNotArchived"}
}

func CategoryIDEquals(categoryID entity.CategoryID) selector.Selector {
	return selector.Selector{Where: squirrel.Expr("category_id = ?", categoryID), Name: "CategoryIDEquals"}
}

func OfArchivedCategorys() selector.Selector {
	return selector.Selector{
		Join:  squirrel.Expr("INNER JOIN categories ON quizzes.category_id = categories.id"),
		Where: squirrel.Expr("categories.archived = TRUE"),
		Name:  "OfArchivedCategory",
	}
}
