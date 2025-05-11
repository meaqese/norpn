package grpc

import (
	"context"
	"github.com/meaqese/norpn/internal/orch/services"
	"github.com/meaqese/norpn/internal/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type Server struct {
	pb.UnimplementedOrchServiceServer
	calculator *services.Calculator
}

func New(calculator *services.Calculator) *Server {
	return &Server{
		calculator: calculator,
	}
}

func (s *Server) GetTask(_ context.Context, _ *pb.Empty) (*pb.TaskResponse, error) {
	task := s.calculator.DequeueTask()
	var response *pb.TaskResponse
	if task == nil {
		return response, status.Errorf(codes.NotFound, "tasks not available now")
	}

	return &pb.TaskResponse{
		ID:            task.ID,
		Arg1:          float32(task.Arg1),
		Arg2:          float32(task.Arg2),
		Operation:     task.Operation,
		OperationTime: int32(task.OperationTime),
	}, nil
}

func (s *Server) SendTaskResult(_ context.Context, taskResult *pb.TaskResult) (*pb.Empty, error) {
	if ch, ok := s.calculator.GetChannelByID(taskResult.ID); ok {
		*ch <- float64(taskResult.Result)
		log.Printf("Received solve for task %s", taskResult.ID)
	}

	return &pb.Empty{}, nil
}
