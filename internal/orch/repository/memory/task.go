package repository

import (
	"github.com/meaqese/norpn/internal/orch/domain"
	"math/rand"
	"strconv"
	"sync"
)

type TaskRepo struct {
	Mu sync.Mutex

	TaskQueue          []*domain.Task
	TaskResultChannels map[string]*chan float64
}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{
		Mu:                 sync.Mutex{},
		TaskQueue:          make([]*domain.Task, 0),
		TaskResultChannels: make(map[string]*chan float64),
	}
}

func (tr *TaskRepo) Dequeue() *domain.Task {
	tr.Mu.Lock()
	defer tr.Mu.Unlock()
	if len(tr.TaskQueue) > 0 {
		task := tr.TaskQueue[0]
		tr.TaskQueue = tr.TaskQueue[1:]

		return task
	}
	return nil
}

func (tr *TaskRepo) Enqueue(task *domain.Task) {
	tr.Mu.Lock()
	defer tr.Mu.Unlock()

	tr.TaskQueue = append(tr.TaskQueue, task)

	channel := make(chan float64, 1)
	tr.TaskResultChannels[task.ID] = &channel
}

func (tr *TaskRepo) GenerateID() string {
	tr.Mu.Lock()
	defer tr.Mu.Unlock()
	for {
		generatedID := strconv.FormatInt(rand.Int63(), 10)
		if _, ok := tr.TaskResultChannels[generatedID]; !ok {
			return generatedID
		}
	}
}

func (tr *TaskRepo) GetChannelByID(id string) (*chan float64, bool) {
	tr.Mu.Lock()
	defer tr.Mu.Unlock()
	value, ok := tr.TaskResultChannels[id]
	return value, ok
}

func (tr *TaskRepo) RemoveChannelByID(id string) {
	tr.Mu.Lock()
	defer tr.Mu.Unlock()
	delete(tr.TaskResultChannels, id)
}
