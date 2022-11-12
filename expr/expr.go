package expr

import (
	"exampleproject/entity"

	"github.com/Masterminds/squirrel"
)

type Expr struct {
	SQL  squirrel.Sqlizer
	Name string
}

func QuizIDEquals(quizID entity.QuizID) Expr {
	return Expr{SQL: squirrel.Expr("QuizID = ?", quizID), Name: "QuizIDEquals"}
}

func ChoiceIDEquals(choiceID entity.ChoiceID) Expr {
	return Expr{SQL: squirrel.Expr("ChoiceID = ?", choiceID), Name: "ChoiceIDEquals"}
}

func CategoryIDEquals(categoryID entity.CategoryID) Expr {
	return Expr{SQL: squirrel.Expr("CategoryID = ?", categoryID), Name: "CategoryIDEquals"}
}

func ChoiceIsCorrect() Expr {
	return Expr{SQL: squirrel.Expr("IsCorrect = 1"), Name: "ChoiceIsCorrect"}
}

func IsNotArchived() Expr {
	return Expr{SQL: squirrel.Expr("Archived = 1"), Name: "IsNotArchived"}
}
