package choice

import (
	"exampleproject/entity"
	"exampleproject/entity/selector"

	"github.com/Masterminds/squirrel"
)

func IDEquals(id entity.ChoiceID) selector.Selector {
	return selector.Selector{Where: squirrel.Expr("id = ?", id), Name: "IDEquals"}
}

func QuizIDEquals(quizID entity.QuizID) selector.Selector {
	return selector.Selector{Where: squirrel.Expr("quiz_id = ?", quizID), Name: "QuizIDEquals"}
}

func IsNotArchived() selector.Selector {
	return selector.Selector{Where: squirrel.Expr("archived = FALSE"), Name: "IsNotArchived"}
}

func IsCorrect() selector.Selector {
	return selector.Selector{Where: squirrel.Expr("is_correct = TRUE"), Name: "IsCorrect"}
}
