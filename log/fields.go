package log

import (
	"exampleproject/entity"
	"exampleproject/repository/expression"
	"fmt"

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

func Exprs(exprs ...expression.Selector) multi {
	fields := make([]zap.Field, 0, len(exprs))
	for _, expr := range exprs {
		fields = append(fields, zap.String(expr.Name, squirrel.DebugSqlizer(expr.Where)))
	}
	return newMulti(fields...)
}

func SQLQuery(builder squirrel.Sqlizer) single {
	return newSingle(zap.String("Query", squirrel.DebugSqlizer(builder)))
}

func Err(err error) single {
	return newSingle(zap.String("Error", err.Error()))
}

func UnexpectedType(actual any, expected any) multi {
	return newMulti(
		zap.String("UnexpectedType", fmt.Sprintf("%T", actual)),
		zap.String("ExpectedType", fmt.Sprintf("%T", expected)),
	)
}
