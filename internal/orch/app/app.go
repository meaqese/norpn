package orch

import (
	"database/sql"
	conf "github.com/meaqese/norpn/internal/orch/config"
	memory "github.com/meaqese/norpn/internal/orch/repository/memory"
	sqlite "github.com/meaqese/norpn/internal/orch/repository/sqlite"
	"github.com/meaqese/norpn/internal/orch/services"
	"github.com/meaqese/norpn/internal/orch/transport/rest"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Orch struct {
	config *conf.Config
}

func New() *Orch {
	return &Orch{
		config: conf.FromEnv(),
	}
}

func (o *Orch) Run() error {
	addr := ":" + o.config.Port

	db, err := sql.Open("sqlite3", "storage.db")
	if err != nil {
		return err
	}

	_, err = sqlite.NewUserRepo(db)
	if err != nil {
		return err
	}

	expressionRepo, err := sqlite.NewExpressionRepo(db)
	if err != nil {
		return err
	}

	taskRepo := memory.NewTaskRepo()

	calc := services.New(taskRepo, expressionRepo, services.CalcOptions{
		TimeAdditionMs:        o.config.TimeAdditionMs,
		TimeSubtractionMs:     o.config.TimeSubtractionMs,
		TimeMultiplicationsMs: o.config.TimeMultiplicationsMs,
		TimeDivisionsMs:       o.config.TimeDivisionsMs,
	})

	handler := rest.New(calc)

	fs := http.FileServer(http.Dir("./web"))

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", handler.Cors(handler.HandleExpression))
	//mux.HandleFunc("/api/v1/expressions/{id}", handler.Cors(handler.HandleGetExpression))
	//mux.HandleFunc("/api/v1/expressions", handler.Cors(handler.HandleGetExpressions))

	mux.HandleFunc("/internal/task", handler.HandleTask)

	mux.Handle("/", fs)

	log.Println("Orchestrator server started on " + addr)
	http.ListenAndServe(addr, mux)

	return nil
}
