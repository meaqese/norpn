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

	fs := http.FileServer(http.Dir("./web"))

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", handler.Cors(handler.HandleExpression))
	mux.HandleFunc("/api/v1/expressions/{id}", handler.Cors(handler.HandleGetExpression))
	mux.HandleFunc("/api/v1/expressions", handler.Cors(handler.HandleGetExpressions))

	mux.HandleFunc("/internal/task", handler.HandleTask)

	mux.Handle("/", fs)

	log.Println("Orchestrator server started on " + addr)
	http.ListenAndServe(addr, mux)
}
