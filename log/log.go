package log

import (
	"go.uber.org/zap"
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
