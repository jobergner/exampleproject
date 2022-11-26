package query

import (
	"exampleproject/entity"
	"exampleproject/repository/selector"

	"github.com/Masterminds/squirrel"
	"golang.org/x/exp/maps"
)

func NewSelectBuilder(meta entity.Meta, selectors ...selector.Selector) squirrel.SelectBuilder {
	cols := meta.ColumnMap()
	for _, s := range selectors {
		for _, c := range s.Columns {
			cols[c] = struct{}{}
		}
	}

	builder := squirrel.Select(maps.Keys(cols)...).From(meta.TableName)

	for _, s := range selectors {
		if s.Where != nil {
			builder = builder.Where(s.Where)
		}

		if s.Join != nil {
			builder = builder.JoinClause(s.Join)
		}
	}

	return builder
}

func NewUpdateBuilder(meta entity.Meta, setMap map[string]interface{}, selectors ...selector.Selector) squirrel.UpdateBuilder {
	builder := squirrel.Update(meta.TableName).SetMap(setMap)

	for _, s := range selectors {
		if s.Where != nil {
			builder = builder.Where(s.Where)
		}
	}

	return builder
}
