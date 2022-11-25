package selector

import (
	"exampleproject/entity"

	"github.com/Masterminds/squirrel"
)

type Selector struct {
	Where   squirrel.Sqlizer
	Columns []string
	Join    squirrel.Sqlizer
	Name    string
}

func (s Selector) ApplyColumns(columns map[string]struct{}) {
	for _, c := range s.Columns {
		columns[c] = struct{}{}
	}
}

func (s Selector) ApplyExpressions(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
	if s.Where != nil {
		builder = builder.Where(s.Where)
	}

	if s.Join != nil {
		builder = builder.JoinClause(s.Join)
	}

	return builder
}

func IDEquals(id int) Selector {
	return Selector{Where: squirrel.Expr("ID = ?", id), Name: "IDEquals"}
}

func QuizIDEquals(quizID entity.QuizID) Selector {
	return Selector{Where: squirrel.Expr("QuizID = ?", quizID), Name: "QuizIDEquals"}
}

func ChoiceIDEquals(choiceID entity.ChoiceID) Selector {
	return Selector{Where: squirrel.Expr("ChoiceID = ?", choiceID), Name: "ChoiceIDEquals"}
}

func CategoryIDEquals(categoryID entity.CategoryID) Selector {
	return Selector{Where: squirrel.Expr("CategoryID = ?", categoryID), Name: "CategoryIDEquals"}
}

func ChoiceIsCorrect() Selector {
	return Selector{Where: squirrel.Expr("IsCorrect = 1"), Name: "ChoiceIsCorrect"}
}

func TitleEquals(title string) Selector {
	return Selector{Where: squirrel.Expr("Title = ?", title), Name: "ChoiceIsCorrect"}
}

func IsNotArchived() Selector {
	return Selector{Where: squirrel.Expr("Archived = 1"), Name: "IsNotArchived"}
}

func UserNameEquals(name string) Selector {
	return Selector{Where: squirrel.Expr("Name = ?", name), Name: "UserNameEquals"}
}

func QuizOfArchivedCategory() Selector {
	return Selector{
		Join:  squirrel.Expr("INNER JOIN Categories ON Quiz.CategoryID = Categories.ID"),
		Where: squirrel.Expr("Categories.Archived = 1"),
		Name:  "QuizOfArchivedCategory",
	}
}
