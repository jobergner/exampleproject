package category

import (
	"exampleproject/entity"
	"exampleproject/repository/selector"

	"github.com/Masterminds/squirrel"
)

func IDEquals(id entity.CategoryID) selector.Selector {
	return selector.Selector{Where: squirrel.Expr("ID = ?", id), Name: "IDEquals"}
}

func TitleEquals(title string) selector.Selector {
	return selector.Selector{Where: squirrel.Expr("Title = ?", title), Name: "TitleEquals"}
}

func IsNotArchived() selector.Selector {
	return selector.Selector{Where: squirrel.Expr("Archived = FALSE"), Name: "IsNotArchived"}
}
