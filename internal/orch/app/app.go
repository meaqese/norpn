package orch

import (
	"database/sql"
	conf "github.com/meaqese/norpn/internal/orch/config"
	memory "github.com/meaqese/norpn/internal/orch/repository/memory"
	sqlite "github.com/meaqese/norpn/internal/orch/repository/sqlite"
	"github.com/meaqese/norpn/internal/orch/services"
	grpcServerHandler "github.com/meaqese/norpn/internal/orch/transport/grpc"
	"github.com/meaqese/norpn/internal/orch/transport/rest"
	"github.com/meaqese/norpn/internal/pb"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

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

	userRepo, err := sqlite.NewUserRepo(db)
	if err != nil {
		return err
	}

	expressionRepo, err := sqlite.NewExpressionRepo(db)
	if err != nil {
		return err
	}

	taskRepo := memory.NewTaskRepo()

	calcSvc := services.New(taskRepo, expressionRepo, services.CalcOptions{
		TimeAdditionMs:        o.config.TimeAdditionMs,
		TimeSubtractionMs:     o.config.TimeSubtractionMs,
		TimeMultiplicationsMs: o.config.TimeMultiplicationsMs,
		TimeDivisionsMs:       o.config.TimeDivisionsMs,
	})

	authSrv := services.NewAuth(userRepo, o.config.JWTSecret)

	handler := rest.New(calcSvc, authSrv)

	grpcAddr := "localhost:" + o.config.GRPCPort
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return err
	}

	grpcHandler := grpcServerHandler.New(calcSvc)
	grpcServer := grpc.NewServer()

	pb.RegisterOrchServiceServer(grpcServer, grpcHandler)

	log.Println("Orchestrator server started on " + addr)
	log.Println("Orchestrator gRPC server started on " + grpcAddr)
	go http.ListenAndServe(addr, handler)
	go grpcServer.Serve(lis)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	grpcServer.GracefulStop()

	return nil
}
