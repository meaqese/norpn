package agent

import (
	"github.com/meaqese/norpn/internal/agent/workers"
	"log"
	"os"
	"os/signal"
	"strconv"
)

type Agent struct {
	ComputingPower int
}

func New() *Agent {
	computingPower, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil {
		computingPower = 1
	}

	return &Agent{
		ComputingPower: computingPower,
	}
}

func (o *Agent) Run() {
	for i := 0; i < o.ComputingPower; i++ {
		go workers.StartWorker()
	}

	log.Printf("Agent started with %d workers", o.ComputingPower)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}
