package repository

import (
	"context"
	"database/sql"
	"errors"
	"exampleproject/db"
	"exampleproject/entity"
	"exampleproject/log"
	"exampleproject/repository/query"
	"exampleproject/repository/selector"
)

type (
	Repository struct {
		meta entity.Meta
	}

	Collection struct {
		Category *Repository
		Quiz     *Repository
		Choice   *Repository
		User     *Repository
	}
)

var Default = Collection{
	Category: NewRepository(entity.CategoryMeta),
	Quiz:     NewRepository(entity.QuizMeta),
	Choice:   NewRepository(entity.ChoiceMeta),
	User:     NewRepository(entity.UserMeta),
}

func NewRepository(meta entity.Meta) *Repository {
	return &Repository{meta}
}

func (r *Repository) Get(ctx context.Context, dest any, selectors ...selector.Selector) error {
	builder := query.NewSelectBuilder(r.meta, selectors...)

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		log.Log(log.QueryBuild, log.Err(err), log.Selectors(selectors...))
		return err
	}

	if err := db.DB.SelectContext(ctx, dest, sqlStr, args...); err != nil {
		log.Log(log.QuerySelect, log.Err(err), log.SQLQuery(builder))
		return err
	}

	return nil
}

func (r *Repository) List(ctx context.Context, dest any, selectors ...selector.Selector) error {
	builder := query.NewSelectBuilder(r.meta, selectors...)

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		log.Log(log.QueryBuild, log.Err(err), log.Selectors(selectors...))
		return err
	}

	if err := db.DB.SelectContext(ctx, dest, sqlStr, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		log.Log(log.QuerySelect, log.Err(err), log.SQLQuery(builder))
		return err
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, transaction *db.TSX, setMap map[string]interface{}, selectors ...selector.Selector) error {
	builder := query.NewUpdateBuilder(r.meta, setMap, selectors...)

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		log.Log(log.QueryBuild, log.Err(err), log.Selectors(selectors...))
		return err
	}

	if _, err := transaction.Get().ExecContext(ctx, sqlStr, args...); err != nil {
		log.Log(log.QueryUpdate, log.Err(err), log.SQLQuery(builder))
		return err
	}

	return nil
}

func (r *Repository) Create(ctx context.Context, transaction *db.TSX, item any) (int64, error) {
	result, err := transaction.Get().NamedExecContext(ctx, r.meta.ToInsertQueryString(), item)
	if err != nil {
		log.Log(log.QueryCreate, log.Err(err))
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Log(log.EvalResultID, log.Err(err))
		return 0, err
	}

	return id, nil
}
