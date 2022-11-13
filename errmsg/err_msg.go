package errmsg

import "fmt"

func NotFound(name string) string {
	return fmt.Sprintf("couln't not find %s", name)
}

func TooManyResults(name string) string {
	return fmt.Sprintf("found more than one %s", name)
}

func QueryBuild(name string) string {
	return fmt.Sprintf("failed building %s query", name)
}

func QuerySelect(name string) string {
	return fmt.Sprintf("failed selecting %s", name)
}

func QueryUpdate(name string) string {
	return fmt.Sprintf("failed updating %s", name)
}

func QueryCreate(name string) string {
	return fmt.Sprintf("failed creating %s", name)
}

func EvalResultID(name string) string {
	return fmt.Sprintf("failed evaluating last inserted %s ID", name)
}
