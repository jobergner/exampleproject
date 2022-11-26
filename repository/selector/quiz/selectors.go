package quiz

import (
	"exampleproject/entity"
	"exampleproject/repository/selector"

	"github.com/Masterminds/squirrel"
)

func IDEquals(id entity.QuizID) selector.Selector {
	return selector.Selector{Where: squirrel.Expr("ID = ?", id), Name: "IDEquals"}
}

func TitleEquals(title string) selector.Selector {
	return selector.Selector{Where: squirrel.Expr("Title = ?", title), Name: "TitleEquals"}
}

func IsNotArchived() selector.Selector {
	return selector.Selector{Where: squirrel.Expr("Archived = FALSE"), Name: "IsNotArchived"}
}

func CategoryIDEquals(categoryID entity.CategoryID) selector.Selector {
	return selector.Selector{Where: squirrel.Expr("CategoryID = ?", categoryID), Name: "CategoryIDEquals"}
}

func OfArchivedCategorys() selector.Selector {
	return selector.Selector{
		Join:  squirrel.Expr("INNER JOIN Categories ON Quizzes.CategoryID = Categories.ID"),
		Where: squirrel.Expr("Categories.Archived = TRUE"),
		Name:  "OfArchivedCategory",
	}
}
