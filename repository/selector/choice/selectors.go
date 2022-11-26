package choice

import (
	"exampleproject/entity"
	"exampleproject/repository/selector"

	"github.com/Masterminds/squirrel"
)

func IDEquals(id entity.ChoiceID) selector.Selector {
	return selector.Selector{Where: squirrel.Expr("ID = ?", id), Name: "IDEquals"}
}

func QuizIDEquals(quizID entity.QuizID) selector.Selector {
	return selector.Selector{Where: squirrel.Expr("QuizID = ?", quizID), Name: "QuizIDEquals"}
}

func IsNotArchived() selector.Selector {
	return selector.Selector{Where: squirrel.Expr("Archived = FALSE"), Name: "IsNotArchived"}
}

func IsCorrect() selector.Selector {
	return selector.Selector{Where: squirrel.Expr("IsCorrect = TRUE"), Name: "IsCorrect"}
}
