// entity holds our database schema entities and some utility values & functions
package entity

type UserID int
type QuizID int
type CategoryID int
type ChoiceID int

var CategoryMeta = Meta{
	TableName: "categories",
	Columns: []string{
		"id",
		"title",
		"archived",
	},
	primaryKey: "id",
}

type Category struct {
	ID       CategoryID `db:"id"`
	Title    string     `db:"title"`
	Archived bool       `db:"archived"`
}

var QuizMeta = Meta{
	TableName: "quizzes",
	Columns: []string{
		"id",
		"category_id",
		"title",
		"description",
		"archived",
	},
	primaryKey: "id",
}

type Quiz struct {
	ID          QuizID     `db:"id"`
	CategoryID  CategoryID `db:"category_id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Archived    bool       `db:"archived"`
}

var ChoiceMeta = Meta{
	TableName: "choices",
	Columns: []string{
		"id",
		"quiz_id",
		"is_correct",
		"content",
		"archived",
	},
	primaryKey: "id",
}

type Choice struct {
	ID        ChoiceID `db:"id"`
	QuizID    QuizID   `db:"quiz_id"`
	IsCorrect bool     `db:"is_correct"`
	Content   string   `db:"content"`
	Archived  bool     `db:"archived"`
}

var UserMeta = Meta{
	TableName: "users",
	Columns: []string{
		"id",
		"name",
		"password_hash",
		"archived",
	},
	primaryKey: "id",
}

type User struct {
	ID           UserID `db:"id"`
	Name         string `db:"name"`
	PasswordHash []byte `db:"password_hash"`
	Archived     bool   `db:"archived"`
}
