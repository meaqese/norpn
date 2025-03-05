package orch

import (
	conf "github.com/meaqese/norpn/internal/orch/config"
	"github.com/meaqese/norpn/internal/orch/transport/rest"
	"log"
	"net/http"
)

type Orch struct {
	config *conf.Config
}

func New() *Orch {
	return &Orch{
		config: conf.FromEnv(),
	}
}

func (o *Orch) Run() {
	addr := ":" + o.config.Port

	handler := rest.New(o.config)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", handler.HandleExpression)
	mux.HandleFunc("/api/v1/expressions/{id}", handler.HandleGetExpression)
	mux.HandleFunc("/api/v1/expressions", handler.HandleGetExpressions)

	mux.HandleFunc("/internal/task", handler.HandleTask)

	log.Println("Orchestrator server started on " + addr)
	http.ListenAndServe(addr, mux)
}
