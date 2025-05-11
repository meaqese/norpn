package workers

import (
	"context"
	"github.com/meaqese/norpn/internal/pb"
	"log"
	"time"
)

func sendResult(grpcClient *pb.OrchServiceClient, id string, result float32) {
	log.Printf("Solved task ID: %s = %f", id, result)

	(*grpcClient).SendTaskResult(context.TODO(), &pb.TaskResult{ID: id, Result: result})
}

func StartWorker(grpcClient *pb.OrchServiceClient) {
	for {
		task, err := (*grpcClient).GetTask(context.TODO(), &pb.Empty{})
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		log.Printf("Received task %v", task)

		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
		switch task.Operation {
		case '+':
			sendResult(grpcClient, task.ID, task.Arg1+task.Arg2)
		case '-':
			sendResult(grpcClient, task.ID, task.Arg1-task.Arg2)
		case '*':
			sendResult(grpcClient, task.ID, task.Arg1*task.Arg2)
		case '/':
			sendResult(grpcClient, task.ID, task.Arg1/task.Arg2)
		}
	}
}
