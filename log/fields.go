package log

import (
	"exampleproject/entity"
	"exampleproject/repository/expression"

	"github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

func QuizID(id entity.QuizID) single {
	return newSingle(zap.Int("QuizID", int(id)))
}

func CategoryID(id entity.CategoryID) single {
	return newSingle(zap.Int("CategoryID", int(id)))
}

func ChoiceID(id entity.ChoiceID) single {
	return newSingle(zap.Int("ChoiceID", int(id)))
}

func JSONData(data []byte) single {
	return newSingle(zap.ByteString("JSONData", data))
}

func Password(password string) single {
	return newSingle(zap.String("Password", password))
}

func Exprs(exprs ...expression.Expr) multi {
	fields := make([]zap.Field, 0, len(exprs))
	for _, expr := range exprs {
		fields = append(fields, zap.String(expr.Name, squirrel.DebugSqlizer(expr.SQL)))
	}
	return newMulti(fields)
}

func SQLQuery(builder squirrel.Sqlizer) single {
	return newSingle(zap.String("Query", squirrel.DebugSqlizer(builder)))
}

func Err(err error) single {
	return newSingle(zap.String("Error", err.Error()))
}
