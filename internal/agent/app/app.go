package agent

import (
	"github.com/meaqese/norpn/internal/agent/workers"
	"github.com/meaqese/norpn/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"os/signal"
	"strconv"
)

type Agent struct {
	ComputingPower int
	grpcPort       string
}

func New() *Agent {
	computingPower, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil {
		computingPower = 1
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "9090"
	}

	return &Agent{
		ComputingPower: computingPower,
		grpcPort:       grpcPort,
	}
}

func (o *Agent) Run() {
	addr := "localhost:" + o.grpcPort
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewOrchServiceClient(conn)

	for i := 0; i < o.ComputingPower; i++ {
		go workers.StartWorker(&client)
	}

	log.Printf("Agent started with %d workers", o.ComputingPower)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}
