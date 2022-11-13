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

type QuizRepository struct{}

func NewQuizRepository() *QuizRepository {
	return new(QuizRepository)
}

func (q *QuizRepository) Get(ctx context.Context, selectors ...expr.Expr) (*entity.Quiz, error) {
	quizzes, err := q.List(ctx, selectors...)
	if err != nil {
		return nil, err
	}
	if len(quizzes) == 0 {
		log.Log(errmsg.NotFound("Quiz"), log.Err(db.ErrNotFound), log.Exprs(selectors...))
		return nil, db.ErrNotFound
	}
	if len(quizzes) > 1 {
		log.Log(errmsg.TooManyResults("Quiz"), log.Err(db.ErrTooManyResults), log.Exprs(selectors...))
		return nil, db.ErrTooManyResults
	}
	return &quizzes[0], nil
}

func (q *QuizRepository) List(ctx context.Context, selectors ...expr.Expr) ([]entity.Quiz, error) {
	builder := squirrel.Select(entity.QuizMeta.Columns...).From(entity.QuizMeta.TableName)

	for _, expr := range selectors {
		builder = builder.Where(expr.SQL)
	}

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		log.Log(errmsg.QueryBuild("list Quizzes"), log.Err(err), log.Exprs(selectors...))
		return nil, err
	}

	var quizzes []entity.Quiz
	if err := db.DB.SelectContext(ctx, quizzes, sqlStr, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return quizzes, nil
		}
		log.Log(errmsg.QuerySelect("Quizzes"), log.Err(err), log.Exprs(selectors...))
		return nil, err
	}

	return quizzes, nil
}

func (q *QuizRepository) Update(ctx context.Context, transaction *db.TSX, setMap map[string]interface{}, selectors ...expr.Expr) error {
	builder := squirrel.Update(entity.QuizMeta.TableName).SetMap(setMap)

	for _, expr := range selectors {
		builder = builder.Where(expr.SQL)
	}

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		log.Log(errmsg.QueryBuild("update Quiz"), log.Err(err), log.Exprs(selectors...))
		return err
	}

	if _, err := transaction.Get().ExecContext(ctx, sqlStr, args...); err != nil {
		log.Log(errmsg.QueryUpdate("Quiz"), log.Err(err), log.Exprs(selectors...))
		return err
	}

	return nil
}

func (q *QuizRepository) Create(ctx context.Context, transaction *db.TSX, quiz *entity.Quiz) (entity.QuizID, error) {
	result, err := transaction.Get().NamedExecContext(ctx, entity.QuizMeta.ToInsertQueryString(), quiz)
	if err != nil {
		log.Log(errmsg.QueryCreate("Quiz"), log.Err(err))
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Log(errmsg.EvalResultID("Quiz"), log.Err(err))
		return 0, err
	}

	return entity.QuizID(id), nil
}
