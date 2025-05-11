package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/meaqese/norpn/internal/orch/domain"
	"github.com/meaqese/norpn/internal/orch/services"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	}
}

func withAuth(authService *services.AuthService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prefix := "Bearer "
		header := r.Header.Get("Authorization")

		encoder := json.NewEncoder(w)

		if !strings.HasPrefix(header, prefix) {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		id, err := authService.Parse(strings.TrimPrefix(header, prefix))
		if err != nil {
			fmt.Println(err)
			encoder.Encode(&ResponseUser{
				Status: "bad",
				Error:  err.Error(),
			})
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (c *Core) Register(w http.ResponseWriter, r *http.Request) {
	request := &RequestUser{}
	err := json.NewDecoder(r.Body).Decode(request)
	defer r.Body.Close()

	if err != nil || request.Login == "" || request.Password == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	encoder := json.NewEncoder(w)

	_, err = c.auth.Register(domain.User{
		Login:    request.Login,
		Password: request.Password,
	})

	if err != nil {
		encoder.Encode(ResponseUser{Status: "bad", Error: err.Error()})
		return
	}

	encoder.Encode(ResponseUser{Status: "ok"})
}

func (c *Core) Login(w http.ResponseWriter, r *http.Request) {
	request := &RequestUser{}
	err := json.NewDecoder(r.Body).Decode(request)
	defer r.Body.Close()

	if err != nil || request.Login == "" || request.Password == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	token, err := c.auth.Login(domain.User{
		Login:    request.Login,
		Password: request.Password,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(ResponseUser{Status: "ok", JWTToken: token})
}

func (c *Core) HandleExpression(w http.ResponseWriter, r *http.Request) {
	requestData := &RequestExpression{}
	encoder := json.NewEncoder(w)

	err := json.NewDecoder(r.Body).Decode(requestData)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	userId := r.Context().Value("user_id").(int64)

	expressionID, err := c.calculator.RegisterExpression(requestData.Expression, userId)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	go c.calculator.StartCalc(expressionID, requestData.Expression)

	log.Printf("Received expression: '%s'", requestData.Expression)
	w.WriteHeader(http.StatusCreated)
	encoder.Encode(ResponseExpression{ID: expressionID})
}

func (c *Core) HandleGetExpressions(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	var expressions []*Expression

	userID := r.Context().Value("user_id").(int64)
	dbExpressions, err := c.calculator.GetAllExpressions(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	for _, expression := range dbExpressions {
		expressions = append(expressions, &Expression{
			ID:     expression.ID,
			Status: expression.Status,
			Result: expression.Result,
			Reason: expression.Reason,
		})
	}

	encoder.Encode(ResponseExpressions{Expressions: expressions})
}

func (c *Core) HandleGetExpression(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	userId := r.Context().Value("user_id").(int64)
	idInt, _ := strconv.Atoi(id)
	expDb, err := c.calculator.GetExpressionById(int64(idInt))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if expDb.UserID != userId {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	encoder.Encode(Expression{
		ID:     expDb.ID,
		Status: expDb.Status,
		Result: expDb.Result,
		Reason: expDb.Reason,
	})
}

func (c *Core) HandleTask(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	if r.Method == "GET" {
		task := c.calculator.DequeueTask()

		if task == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		encoder.Encode(task)
		log.Printf("Dequeue task %s from store", task.ID)
	} else if r.Method == "POST" {
		taskResult := &TaskResult{}
		err := json.NewDecoder(r.Body).Decode(taskResult)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
		defer r.Body.Close()

		if ch, ok := c.calculator.GetChannelByID(taskResult.ID); ok {
			*ch <- taskResult.Result
			log.Printf("Received solve for task %s", taskResult.ID)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
