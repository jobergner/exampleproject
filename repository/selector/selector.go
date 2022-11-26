package selector

import (
	"github.com/Masterminds/squirrel"
)

type Selector struct {
	Where   squirrel.Sqlizer
	Columns []string
	Join    squirrel.Sqlizer
	Name    string
}
