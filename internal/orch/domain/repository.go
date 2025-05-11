package domain

type UserRepository interface {
	Add() (int64, error)
	GetByLogin()
}

type ExpressionRepository interface {
	Add(expression Expression) (int64, error)
	Update(expression Expression) error
}

type TaskRepository interface {
	Enqueue(task *Task)
	Dequeue() *Task
	GenerateID() string
	GetChannelByID(id string) (*chan float64, bool)
	RemoveChannelByID(id string)
}
