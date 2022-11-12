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

type ChoiceRepository struct{}

func NewChoiceRepository() *ChoiceRepository {
	return new(ChoiceRepository)
}

func (q *ChoiceRepository) Get(ctx context.Context, selectors ...expr.Expr) (*entity.Choice, error) {
	choices, err := q.List(ctx, selectors...)
	if err != nil {
		return nil, err
	}
	if len(choices) == 0 {
		log.Log(errmsg.NotFound("choice"), log.Err(db.ErrNotFound), log.Exprs(selectors...))
		return nil, db.ErrNotFound
	}
	if len(choices) > 1 {
		log.Log(errmsg.TooManyResults("choices"), log.Err(db.ErrTooManyResults), log.Exprs(selectors...))
		return nil, db.ErrTooManyResults
	}
	return &choices[0], nil
}

func (q *ChoiceRepository) List(ctx context.Context, selectors ...expr.Expr) ([]entity.Choice, error) {
	builder := squirrel.Select(entity.ChoiceMeta.Columns...).From(entity.ChoiceMeta.TableName)

	for _, expr := range selectors {
		builder = builder.Where(expr.SQL)
	}

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		log.Log(errmsg.QueryBuild("list choices"), log.Err(err), log.Exprs(selectors...))
		return nil, err
	}

	var choices []entity.Choice
	if err := db.DB.SelectContext(ctx, choices, sqlStr, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return choices, nil
		}
		log.Log(errmsg.QuerySelect("choices"), log.Err(err), log.Exprs(selectors...))
		return nil, err
	}

	return choices, nil
}

func (q *ChoiceRepository) Update(ctx context.Context, transaction *db.TSX, setMap map[string]interface{}, selectors ...expr.Expr) error {
	builder := squirrel.Update(entity.ChoiceMeta.TableName).SetMap(setMap)

	for _, expr := range selectors {
		builder = builder.Where(expr.SQL)
	}

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		log.Log(errmsg.QueryBuild("update choices"), log.Err(err), log.Exprs(selectors...))
		return err
	}

	if _, err := transaction.Get().ExecContext(ctx, sqlStr, args...); err != nil {
		_ = transaction.Rollback()
		log.Log(errmsg.QueryUpdate("choices"), log.Err(err), log.Exprs(selectors...))
		return err
	}

	return nil
}

func (q *ChoiceRepository) Create(ctx context.Context, transaction *db.TSX, choice *entity.Choice) error {
	if _, err := transaction.Get().NamedExecContext(ctx, entity.ChoiceMeta.ToInsertQueryString(), choice); err != nil {
		log.Log(errmsg.QueryCreate("choice"), log.Err(err))
		return err
	}

	return nil
}
