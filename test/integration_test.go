package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	memory "github.com/meaqese/norpn/internal/orch/repository/memory"
	sqlite "github.com/meaqese/norpn/internal/orch/repository/sqlite"
	"github.com/meaqese/norpn/internal/orch/services"
	grpcsrv "github.com/meaqese/norpn/internal/orch/transport/grpc"
	"github.com/meaqese/norpn/internal/orch/transport/rest"
	"github.com/meaqese/norpn/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const ApiEndpoint = "http://norpn.io"

func sendResult(grpcClient *pb.OrchServiceClient, id string, result float32) {
	log.Printf("Solved task ID: %s = %f", id, result)

	_, err := (*grpcClient).SendTaskResult(context.TODO(), &pb.TaskResult{ID: id, Result: result})
	if err != nil {
		fmt.Println(err)
	}
}

func worker(grpcClient *pb.OrchServiceClient, t *testing.T) {
	t.Log("started worker")
	for {
		task, err := (*grpcClient).GetTask(context.TODO(), &pb.Empty{})

		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		t.Log("task received")

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

func getExpression(core *rest.Core, id int64) rest.Expression {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", fmt.Sprintf("%s/api/v1/expressions/%d", ApiEndpoint, id), nil)
	r.SetPathValue("id", strconv.Itoa(int(id)))

	ctx := context.WithValue(r.Context(), "user_id", int64(1))

	core.HandleGetExpression(w, r.WithContext(ctx))

	resp := w.Result()
	defer resp.Body.Close()

	expResult := rest.Expression{}
	expDecoder := json.NewDecoder(resp.Body)
	expDecoder.Decode(&expResult)

	return expResult
}

func addingExpression(tc services.TestCase, t *testing.T, core *rest.Core) {
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(rest.RequestExpression{Expression: tc.Expression})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", ApiEndpoint+"/api/v1/calculate", &buffer)
	w := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), "user_id", int64(1))

	core.HandleExpression(w, req.WithContext(ctx))
	resp := w.Result()
	defer resp.Body.Close()

	decoder := json.NewDecoder(w.Body)

	if resp.StatusCode == 201 {
		resResult := &rest.ResponseExpression{}
		err = decoder.Decode(resResult)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Duration(tc.TimeoutSeconds) * time.Second)

		expression := getExpression(core, resResult.ID)

		log.Println(expression)

		if expression.Result != tc.Expected {
			t.Fatalf("not valid return - %v, expected %v", expression.Result, tc.Expected)
		}
	} else {
		if !tc.ShouldFail {
			t.Fatalf("success test not should fail")
		}
	}
}

func getClientConn(lis *bufconn.Listener) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		"passthrough://bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func TestCalcHandler(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	cases := services.GetTestCases()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)

	taskRepo := memory.NewTaskRepo()
	userRepo, _ := sqlite.NewUserRepo(db)
	expressionRepo, _ := sqlite.NewExpressionRepo(db)

	calcSvc := services.New(taskRepo, expressionRepo, services.CalcOptions{
		TimeAdditionMs:        0,
		TimeSubtractionMs:     0,
		TimeMultiplicationsMs: 0,
		TimeDivisionsMs:       0,
	})

	authSvc := services.NewAuth(userRepo, "secret")

	_, core := rest.New(calcSvc, authSvc)

	s := grpc.NewServer()
	pb.RegisterOrchServiceServer(s, grpcsrv.New(calcSvc))

	go s.Serve(lis)

	conn, err := getClientConn(lis)
	if err != nil {
		t.Fatal(err)
	}

	client := pb.NewOrchServiceClient(conn)

	go worker(&client, t)

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			addingExpression(tc, t, core)
		})
	}
}
