package log

import (
	"go.uber.org/zap"
)

const (
	NotFound       = "could not find item"
	TooManyResults = "found more than one item"
	QueryBuild     = "failed building query"
	QuerySelect    = "failed selecting items"
	QueryUpdate    = "failed updating item"
	QueryCreate    = "failed creating item"
	EvalResultID   = "failed evaluating last inserted item ID"
	Unmarshal      = "failed unmarshalling item"
	ReadBody       = "failed reading request body"
	Serve          = "failed starting server"
	DBConnection   = "failed connecting to DB"
)

var Logger, _ = zap.NewDevelopment(zap.AddStacktrace(zap.ErrorLevel), zap.AddCallerSkip(1))

func Log(msg string, loggables ...Loggable) {
	var fields []zap.Field
	for _, l := range loggables {
		fields = append(fields, l.Fields()...)
	}
	Logger.Error(msg, fields...)
}

type Loggable interface {
	Fields() []zap.Field
}

type single struct {
	field zap.Field
}

func newSingle(field zap.Field) single {
	return single{field: field}
}

func (s single) Fields() []zap.Field {
	return []zap.Field{s.field}
}

type multi struct {
	fields []zap.Field
}

func newMulti(fields []zap.Field) multi {
	return multi{fields: fields}
}

func (m multi) Fields() []zap.Field {
	return m.fields
}
