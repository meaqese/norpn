package domain

type UserRepository interface {
	Add(user User) (int64, error)
	GetByLogin(login string) (User, error)
}

type ExpressionRepository interface {
	Add(expression Expression) (int64, error)
	Update(expression Expression) error
	GetById(id int64) (*Expression, error)
	GetAll(userId int64) ([]*Expression, error)
}

type TaskRepository interface {
	Enqueue(task *Task)
	Dequeue() *Task
	GenerateID() string
	GetChannelByID(id string) (*chan float64, bool)
	RemoveChannelByID(id string)
}
