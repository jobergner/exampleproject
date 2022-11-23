package repository

type Repository struct {
	Category *CategoryRepository
	Quiz     *QuizRepository
	Choice   *ChoiceRepository
	User     *UserRepository
}

var Default = Repository{
	Category: NewCategoryRepository(),
	Quiz:     NewQuizRepository(),
	Choice:   NewChoiceRepository(),
	User:     NewUserRepository(),
}
