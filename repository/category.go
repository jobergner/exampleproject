package repository

import (
	"context"
	"database/sql"
	"errors"
	"exampleproject/db"
	"exampleproject/entity"
	"exampleproject/log"
	"exampleproject/repository/expression"

	"github.com/Masterminds/squirrel"
)

type CategoryRepository struct{}

func NewCategoryRepository() *CategoryRepository {
	return new(CategoryRepository)
}

func (q *CategoryRepository) Get(ctx context.Context, selectors ...expression.Expr) (*entity.Category, error) {
	categories, err := q.List(ctx, selectors...)
	if err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		log.Log(log.NotFound, log.Err(db.ErrNotFound), log.Exprs(selectors...))
		return nil, db.ErrNotFound
	}

	if len(categories) > 1 {
		log.Log(log.TooManyResults, log.Err(db.ErrTooManyResults), log.Exprs(selectors...))
		return nil, db.ErrTooManyResults
	}

	return &categories[0], nil
}

func (q *CategoryRepository) List(ctx context.Context, selectors ...expression.Expr) ([]entity.Category, error) {
	builder := squirrel.Select(entity.CategoryMeta.Columns...).From(entity.CategoryMeta.TableName)

	for _, expr := range selectors {
		builder = builder.Where(expr.SQL)
	}

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		log.Log(log.QueryBuild, log.Err(err), log.Exprs(selectors...))
		return nil, err
	}

	var categories []entity.Category
	if err := db.DB.SelectContext(ctx, categories, sqlStr, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return categories, nil
		}
		log.Log(log.QuerySelect, log.Err(err), log.SQLQuery(builder))
		return nil, err
	}

	return categories, nil
}

func (q *CategoryRepository) Update(ctx context.Context, transaction *db.TSX, setMap map[string]interface{}, selectors ...expression.Expr) error {
	builder := squirrel.Update(entity.CategoryMeta.TableName).SetMap(setMap)

	for _, expr := range selectors {
		builder = builder.Where(expr.SQL)
	}

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		log.Log(log.QueryBuild, log.Err(err), log.Exprs(selectors...))
		return err
	}

	if _, err := transaction.Get().ExecContext(ctx, sqlStr, args...); err != nil {
		log.Log(log.QueryUpdate, log.Err(err), log.SQLQuery(builder))
		return err
	}

	return nil
}

func (q *CategoryRepository) Create(ctx context.Context, transaction *db.TSX, category *entity.Category) (entity.CategoryID, error) {
	result, err := transaction.Get().NamedExecContext(ctx, entity.CategoryMeta.ToInsertQueryString(), category)
	if err != nil {
		log.Log(log.QueryCreate, log.Err(err))
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Log(log.EvalResultID, log.Err(err))
		return 0, err
	}

	return entity.CategoryID(id), nil
}
