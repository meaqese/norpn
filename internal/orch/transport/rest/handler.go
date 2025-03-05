package rest

import (
	"encoding/json"
	"github.com/meaqese/norpn/internal/orch/norpn"
	"net/http"
)

func (c *Core) Cors(next http.HandlerFunc) http.HandlerFunc {
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

func (c *Core) HandleExpression(w http.ResponseWriter, r *http.Request) {
	requestData := &RequestExpression{}
	encoder := json.NewEncoder(w)

	err := json.NewDecoder(r.Body).Decode(requestData)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	hashOfExpression := c.GenerateHash(requestData.Expression)

	c.mu.Lock()
	_, found := c.expressionStore[hashOfExpression]
	c.mu.Unlock()

	if !found {
		c.AddExpression(hashOfExpression)

		go c.StartCalc(requestData.Expression, hashOfExpression)
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(ResponseExpression{ID: hashOfExpression})
}

func (c *Core) HandleGetExpressions(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	var expressions []*Expression
	c.mu.Lock()
	for _, val := range c.expressionStore {
		expressions = append(expressions, val)
	}
	c.mu.Unlock()

	encoder.Encode(ResponseExpressions{Expressions: expressions})
}

func (c *Core) HandleGetExpression(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	c.mu.Lock()
	defer c.mu.Unlock()

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if val, ok := c.expressionStore[id]; ok {
		encoder.Encode(val)
		return
	}

	w.WriteHeader(http.StatusNotFound)
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
	} else if r.Method == "POST" {
		c.calculator.Mu.Lock()
		defer c.calculator.Mu.Unlock()

		taskResult := &norpn.TaskResult{}
		err := json.NewDecoder(r.Body).Decode(taskResult)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
		defer r.Body.Close()

		if ch, ok := c.calculator.TaskResultChannels[taskResult.ID]; ok {
			*ch <- taskResult.Result
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
