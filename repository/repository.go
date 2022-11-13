package repository

type Repository struct {
	Category *CategoryRepository
	Quiz     *QuizRepository
	Choice   *ChoiceRepository
}

var Default = Repository{
	Category: NewCategoryRepository(),
	Quiz:     NewQuizRepository(),
	Choice:   NewChoiceRepository(),
}
