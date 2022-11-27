package repository

import (
	"context"
	"database/sql"
	"errors"
	"exampleproject/db"
	"exampleproject/entity"
	"exampleproject/entity/selector"
	"exampleproject/log"

	"github.com/jmoiron/sqlx"
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

	QueryExecuter interface {
		NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
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
	return r.get(ctx, db.DB, dest, selectors...)
}

func (r *Repository) GetTsx(ctx context.Context, tsx *sqlx.Tx, dest any, selectors ...selector.Selector) error {
	return r.get(ctx, tsx, dest, selectors...)
}

func (r *Repository) get(ctx context.Context, db QueryExecuter, dest any, selectors ...selector.Selector) error {
	builder := selector.NewSelectBuilder(r.meta, selectors...)

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		log.Log(log.QueryBuild, log.Err(err), log.Selectors(selectors...))
		return err
	}

	if err := db.SelectContext(ctx, dest, sqlStr, args...); err != nil {
		log.Log(log.QuerySelect, log.Err(err), log.SQLQuery(builder))
		return err
	}

	return nil
}

func (r *Repository) List(ctx context.Context, dest any, selectors ...selector.Selector) error {
	return r.list(ctx, db.DB, dest, selectors...)
}

func (r *Repository) ListTsx(ctx context.Context, tsx *sqlx.Tx, dest any, selectors ...selector.Selector) error {
	return r.list(ctx, tsx, dest, selectors...)
}

func (r *Repository) list(ctx context.Context, db QueryExecuter, dest any, selectors ...selector.Selector) error {
	builder := selector.NewSelectBuilder(r.meta, selectors...)

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		log.Log(log.QueryBuild, log.Err(err), log.Selectors(selectors...))
		return err
	}

	if err := db.SelectContext(ctx, dest, sqlStr, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		log.Log(log.QuerySelect, log.Err(err), log.SQLQuery(builder))
		return err
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, setMap map[string]interface{}, selectors ...selector.Selector) error {
	return r.update(ctx, db.DB, setMap, selectors...)
}

func (r *Repository) UpdateTsx(ctx context.Context, tsx *sqlx.Tx, setMap map[string]interface{}, selectors ...selector.Selector) error {
	return r.update(ctx, tsx, setMap, selectors...)
}

func (r *Repository) update(ctx context.Context, db QueryExecuter, setMap map[string]interface{}, selectors ...selector.Selector) error {
	builder := selector.NewUpdateBuilder(r.meta, setMap, selectors...)

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		log.Log(log.QueryBuild, log.Err(err), log.Selectors(selectors...))
		return err
	}

	if _, err := db.ExecContext(ctx, sqlStr, args...); err != nil {
		log.Log(log.QueryUpdate, log.Err(err), log.SQLQuery(builder))
		return err
	}

	return nil
}

func (r *Repository) Create(ctx context.Context, item any) (int64, error) {
	return r.create(ctx, db.DB, item)
}

func (r *Repository) CreateTsx(ctx context.Context, tsx *sqlx.Tx, item any) (int64, error) {
	return r.create(ctx, tsx, item)
}

func (r *Repository) create(ctx context.Context, db QueryExecuter, item any) (int64, error) {
	result, err := db.NamedExecContext(ctx, r.meta.ToInsertQueryString(), item)
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
