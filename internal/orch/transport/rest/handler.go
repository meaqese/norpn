package rest

import (
	"encoding/json"
	"github.com/meaqese/norpn/internal/orch/domain"
	"log"
	"net/http"
	"strconv"
)

func (c *Core) Register(w http.ResponseWriter, r *http.Request) {
	request := &RequestUser{}
	err := json.NewDecoder(r.Body).Decode(request)
	defer r.Body.Close()

	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	if err != nil || request.Login == "" || request.Password == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		encoder.Encode(ResponseUser{Status: false, Error: "check login and password"})
		return
	}

	_, err = c.auth.Register(domain.User{
		Login:    request.Login,
		Password: request.Password,
	})

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		encoder.Encode(ResponseUser{Status: false, Error: err.Error()})
		return
	}

	encoder.Encode(ResponseUser{Status: true})
}

func (c *Core) Login(w http.ResponseWriter, r *http.Request) {
	request := &RequestUser{}
	err := json.NewDecoder(r.Body).Decode(request)
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	if err != nil || request.Login == "" || request.Password == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		encoder.Encode(ResponseUser{Status: false, Error: "check login and password"})
		return
	}

	token, err := c.auth.Login(domain.User{
		Login:    request.Login,
		Password: request.Password,
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		encoder.Encode(ResponseUser{Status: false, Error: err.Error()})
		return
	}

	encoder.Encode(ResponseUser{Status: true, JWTToken: token})
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
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
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
