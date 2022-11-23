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

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return new(UserRepository)
}

func (q *UserRepository) Get(ctx context.Context, selectors ...expression.Expr) (*entity.User, error) {
	users, err := q.List(ctx, selectors...)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		log.Log(log.NotFound, log.Err(db.ErrNotFound), log.Exprs(selectors...))
		return nil, db.ErrNotFound
	}

	if len(users) > 1 {
		log.Log(log.TooManyResults, log.Err(db.ErrTooManyResults), log.Exprs(selectors...))
		return nil, db.ErrTooManyResults
	}

	return &users[0], nil
}

func (q *UserRepository) List(ctx context.Context, selectors ...expression.Expr) ([]entity.User, error) {
	builder := squirrel.Select(entity.UserMeta.Columns...).From(entity.UserMeta.TableName)

	for _, expr := range selectors {
		builder = builder.Where(expr.SQL)
	}

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		log.Log(log.QueryBuild, log.Err(err), log.Exprs(selectors...))
		return nil, err
	}

	var users []entity.User
	if err := db.DB.SelectContext(ctx, users, sqlStr, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return users, nil
		}
		log.Log(log.QuerySelect, log.Err(err), log.SQLQuery(builder))
		return nil, err
	}

	return users, nil
}

func (q *UserRepository) Update(ctx context.Context, transaction *db.TSX, setMap map[string]interface{}, selectors ...expression.Expr) error {
	builder := squirrel.Update(entity.UserMeta.TableName).SetMap(setMap)

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

func (q *UserRepository) Create(ctx context.Context, transaction *db.TSX, user *entity.User) (entity.UserID, error) {
	result, err := transaction.Get().NamedExecContext(ctx, entity.UserMeta.ToInsertQueryString(), user)
	if err != nil {
		log.Log(log.QueryCreate, log.Err(err))
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Log(log.EvalResultID, log.Err(err))
		return 0, err
	}

	return entity.UserID(id), nil
}
