package entity

import (
	"fmt"
	"strings"
)

type Meta struct {
	TableName  string
	Columns    []string
	primaryKey string
}

func (m Meta) ToInsertQueryString() string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("INSERT INTO %s SET\n", m.TableName))

	for i, c := range m.Columns {
		if c == m.primaryKey {
			continue
		}

		if i == len(m.Columns)-1 {
			sb.WriteString(fmt.Sprintf("%s = : %s\n", c, c))
		} else {
			sb.WriteString(fmt.Sprintf("%s = : %s\n,", c, c))
		}
	}

	return sb.String()
}

func (m Meta) ColumnMap() map[string]struct{} {
	cols := make(map[string]struct{})

	for _, c := range m.Columns {
		name := fmt.Sprintf("%s.%s", m.TableName, c)
		cols[name] = struct{}{}
	}

	return cols
}
