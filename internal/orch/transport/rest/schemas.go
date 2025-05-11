package rest

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Expression struct {
	ID     int64   `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result,omitempty"`
	Reason string  `json:"reason,omitempty"`
}

type RequestExpression struct {
	Expression string `json:"expression"`
}

type ResponseExpression struct {
	ID int64 `json:"id"`
}

type ResponseExpressions struct {
	Expressions []*Expression `json:"expressions"`
}

type TaskResult struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}
