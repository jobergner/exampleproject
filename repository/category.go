package repository

import (
	"context"
	"database/sql"
	"errors"
	"exampleproject/db"
	"exampleproject/entity"
	"exampleproject/errmsg"
	"exampleproject/expr"
	"exampleproject/log"

	"github.com/Masterminds/squirrel"
)

type CategoryRepository struct{}

func NewCategoryRepository() *CategoryRepository {
	return new(CategoryRepository)
}

func (q *CategoryRepository) Get(ctx context.Context, selectors ...expr.Expr) (*entity.Category, error) {
	categories, err := q.List(ctx, selectors...)
	if err != nil {
		return nil, err
	}
	if len(categories) == 0 {
		log.Log(errmsg.NotFound("Category"), log.Err(db.ErrNotFound), log.Exprs(selectors...))
		return nil, db.ErrNotFound
	}
	if len(categories) > 1 {
		log.Log(errmsg.TooManyResults("Category"), log.Err(db.ErrTooManyResults), log.Exprs(selectors...))
		return nil, db.ErrTooManyResults
	}
	return &categories[0], nil
}

func (q *CategoryRepository) List(ctx context.Context, selectors ...expr.Expr) ([]entity.Category, error) {
	builder := squirrel.Select(entity.CategoryMeta.Columns...).From(entity.CategoryMeta.TableName)

	for _, expr := range selectors {
		builder = builder.Where(expr.SQL)
	}

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		log.Log(errmsg.QueryBuild("list Categories"), log.Err(err), log.Exprs(selectors...))
		return nil, err
	}

	var categories []entity.Category
	if err := db.DB.SelectContext(ctx, categories, sqlStr, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return categories, nil
		}
		log.Log(errmsg.QuerySelect("Categories"), log.Err(err), log.Exprs(selectors...))
		return nil, err
	}

	return categories, nil
}

func (q *CategoryRepository) Update(ctx context.Context, transaction *db.TSX, setMap map[string]interface{}, selectors ...expr.Expr) error {
	builder := squirrel.Update(entity.CategoryMeta.TableName).SetMap(setMap)

	for _, expr := range selectors {
		builder = builder.Where(expr.SQL)
	}

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		log.Log(errmsg.QueryBuild("update Category"), log.Err(err), log.Exprs(selectors...))
		return err
	}

	if _, err := transaction.Get().ExecContext(ctx, sqlStr, args...); err != nil {
		log.Log(errmsg.QueryUpdate("Category"), log.Err(err), log.Exprs(selectors...))
		return err
	}

	return nil
}

func (q *CategoryRepository) Create(ctx context.Context, transaction *db.TSX, category *entity.Category) (entity.CategoryID, error) {
	result, err := transaction.Get().NamedExecContext(ctx, entity.CategoryMeta.ToInsertQueryString(), category)
	if err != nil {
		log.Log(errmsg.QueryCreate("Category"), log.Err(err))
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Log(errmsg.EvalResultID("Category"), log.Err(err))
		return 0, err
	}

	return entity.CategoryID(id), nil
}
