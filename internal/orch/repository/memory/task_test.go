package repository

import (
	"github.com/meaqese/norpn/internal/orch/domain"
	"testing"
)

func addTestTask(repo *TaskRepo) *domain.Task {
	task := &domain.Task{
		ID:            "123",
		Arg1:          1,
		Arg2:          2,
		Operation:     43,
		OperationTime: 100,
	}
	repo.Enqueue(task)

	return task
}

func TestTaskRepo_Enqueue(t *testing.T) {
	repo := NewTaskRepo()
	task := addTestTask(repo)

	if repo.TaskQueue[0].ID != task.ID {
		t.Fatalf("Task ID should be 123")
	}
}

func TestTaskRepo_Dequeue(t *testing.T) {
	repo := NewTaskRepo()
	task := addTestTask(repo)

	newTask := repo.Dequeue()

	if newTask.ID != task.ID {
		t.Fatalf("Task ID should be %s", newTask.ID)
	}
}
