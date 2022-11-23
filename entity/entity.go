// entity holds our database schema entities and some utility values & functions
package entity

type UserID int
type QuizID int
type CategoryID int
type ChoiceID int

var CategoryMeta = Meta{
	TableName: "Categories",
	Columns: []string{
		"ID",
		"Title",
		"Archived",
	},
	primaryKey: "ID",
}

type Category struct {
	ID       CategoryID
	Title    string
	Archived bool
}

var QuizMeta = Meta{
	TableName: "Quizzes",
	Columns: []string{
		"ID",
		"CategoryID",
		"Title",
		"Description",
		"Archived",
	},
	primaryKey: "ID",
}

type Quiz struct {
	ID          QuizID
	CategoryID  CategoryID
	Title       string
	Description string
	Archived    bool
}

var ChoiceMeta = Meta{
	TableName: "Choices",
	Columns: []string{
		"ID",
		"QuizID",
		"IsCorrect",
		"Content",
		"Archived",
	},
	primaryKey: "ID",
}

type Choice struct {
	ID        ChoiceID
	QuizID    QuizID
	IsCorrect bool
	Content   string
	Archived  bool
}

var UserMeta = Meta{
	TableName: "Users",
	Columns: []string{
		"ID",
		"Name",
		"PasswordHash",
		"Archived",
	},
	primaryKey: "ID",
}

type User struct {
	ID           UserID
	Name         string
	PasswordHash []byte
	Archived     bool
}
