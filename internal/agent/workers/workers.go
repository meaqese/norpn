package workers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/meaqese/norpn/internal/orch/domain"
	"github.com/meaqese/norpn/internal/orch/transport/rest"
	"log"
	"net/http"
	"time"
)

func sendResult(id string, result float64) {
	data, _ := json.Marshal(rest.TaskResult{ID: id, Result: result})
	dataReader := bytes.NewReader(data)

	log.Printf("Solved task ID: %s = %f", id, result)

	_, err := http.Post("http://localhost:8080/internal/task", "application/json", dataReader)
	if err != nil {
		fmt.Println(err)
	}
}

func StartWorker() {
	for {
		resp, err := http.Get("http://localhost:8080/internal/task")
		if err != nil || resp.StatusCode != 200 {
			time.Sleep(1 * time.Second)
			continue
		}

		task := &domain.Task{}
		err = json.NewDecoder(resp.Body).Decode(task)
		if err != nil {
			log.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}

		log.Printf("Received task %v", task)

		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
		switch task.Operation {
		case '+':
			sendResult(task.ID, task.Arg1+task.Arg2)
		case '-':
			sendResult(task.ID, task.Arg1-task.Arg2)
		case '*':
			sendResult(task.ID, task.Arg1*task.Arg2)
		case '/':
			sendResult(task.ID, task.Arg1/task.Arg2)
		}
	}
}
