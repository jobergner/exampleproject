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

type QuizRepository struct{}

func NewQuizRepository() *QuizRepository {
	return new(QuizRepository)
}

func (q *QuizRepository) Get(ctx context.Context, selectors ...expression.Expr) (*entity.Quiz, error) {
	quizzes, err := q.List(ctx, selectors...)
	if err != nil {
		return nil, err
	}

	if len(quizzes) == 0 {
		log.Log(log.NotFound, log.Err(db.ErrNotFound), log.Exprs(selectors...))
		return nil, db.ErrNotFound
	}

	if len(quizzes) > 1 {
		log.Log(log.TooManyResults, log.Err(db.ErrTooManyResults), log.Exprs(selectors...))
		return nil, db.ErrTooManyResults
	}

	return &quizzes[0], nil
}

func (q *QuizRepository) List(ctx context.Context, selectors ...expression.Expr) ([]entity.Quiz, error) {
	builder := squirrel.Select(entity.QuizMeta.Columns...).From(entity.QuizMeta.TableName)

	for _, expr := range selectors {
		builder = builder.Where(expr.SQL)
	}

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		log.Log(log.QueryBuild, log.Err(err), log.Exprs(selectors...))
		return nil, err
	}

	var quizzes []entity.Quiz
	if err := db.DB.SelectContext(ctx, quizzes, sqlStr, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return quizzes, nil
		}
		log.Log(log.QuerySelect, log.Err(err), log.SQLQuery(builder))
		return nil, err
	}

	return quizzes, nil
}

func (q *QuizRepository) Update(ctx context.Context, transaction *db.TSX, setMap map[string]interface{}, selectors ...expression.Expr) error {
	builder := squirrel.Update(entity.QuizMeta.TableName).SetMap(setMap)

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

func (q *QuizRepository) Create(ctx context.Context, transaction *db.TSX, quiz *entity.Quiz) (entity.QuizID, error) {
	result, err := transaction.Get().NamedExecContext(ctx, entity.QuizMeta.ToInsertQueryString(), quiz)
	if err != nil {
		log.Log(log.QueryCreate, log.Err(err))
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Log(log.EvalResultID, log.Err(err))
		return 0, err
	}

	return entity.QuizID(id), nil
}
