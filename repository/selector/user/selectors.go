package category

import (
	"exampleproject/entity"
	"exampleproject/repository/selector"

	"github.com/Masterminds/squirrel"
)

func IDEquals(id entity.CategoryID) selector.Selector {
	return selector.Selector{Where: squirrel.Expr("ID = ?", id), Name: "IDEquals"}
}

func NameEquals(name string) selector.Selector {
	return selector.Selector{Where: squirrel.Expr("Name = ?", name), Name: "NameEquals"}
}

func IsNotArchived() selector.Selector {
	return selector.Selector{Where: squirrel.Expr("Archived = FALSE"), Name: "IsNotArchived"}
}
