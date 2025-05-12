package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	conf "github.com/meaqese/norpn/internal/orch/config"
	memory "github.com/meaqese/norpn/internal/orch/repository/memory"
	sqlite "github.com/meaqese/norpn/internal/orch/repository/sqlite"
	"github.com/meaqese/norpn/internal/orch/services"
	"github.com/meaqese/norpn/internal/orch/transport/rest"
	"log"
	"net/http/httptest"
	"testing"
	"time"
)

const ApiEndpoint = "http://norpn.io"

func sendResult(core *rest.Core, id string, result float64) {
	var buf bytes.Buffer

	json.NewEncoder(&buf).Encode(services.TaskResult{
		ID:     id,
		Result: result,
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/internal/task", &buf)

	core.HandleTask(w, r)
}

func worker(core *rest.Core) {
	for {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/internal/task", nil)

		core.HandleTask(w, r)

		resp := w.Result()

		if resp.StatusCode != 200 {
			time.Sleep(1 * time.Second)
			resp.Body.Close()
			continue
		}

		task := &services.Task{}
		err := json.NewDecoder(resp.Body).Decode(task)
		resp.Body.Close()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		switch task.Operation {
		case '+':
			sendResult(core, task.ID, task.Arg1+task.Arg2)
		case '-':
			sendResult(core, task.ID, task.Arg1-task.Arg2)
		case '*':
			sendResult(core, task.ID, task.Arg1*task.Arg2)
		case '/':
			sendResult(core, task.ID, task.Arg1/task.Arg2)
		}
	}
}

func getExpression(core *rest.Core, id string) rest.Expression {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", fmt.Sprintf("%s/api/v1/expressions/%s", ApiEndpoint, id), nil)
	r.SetPathValue("id", id)
	core.HandleGetExpression(w, r)

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

	core.HandleExpression(w, req)
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
			t.Fatalf("not valid return - %v", expression.Result)
		}
	} else {
		if !tc.ShouldFail {
			t.Fatalf("success test not should fail")
		}
	}
}

func TestCalcHandler(t *testing.T) {
	cases := services.GetTestCases()

	db, _ := sql.Open("sqlite3", ":memory:")

	taskRepo := memory.NewTaskRepo()
	userRepo, _ := sqlite.NewUserRepo(db)
	expressionRepo, _ := sqlite.NewExpressionRepo(db)

	calcSvc := services.New(taskRepo, expressionRepo, services.CalcOptions{
		TimeAdditionMs:        0,
		TimeSubtractionMs:     0,
		TimeMultiplicationsMs: 0,
		TimeDivisionsMs:       0,
	})

	core := rest.New(calcSvc)

	go worker(core)

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			addingExpression(tc, t, core)
		})
	}
}
