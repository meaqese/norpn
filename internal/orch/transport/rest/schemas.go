package rest

type RequestUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ResponseUser struct {
	Status   bool   `json:"status"`
	Error    string `json:"error,omitempty"`
	JWTToken string `json:"jwt_token,omitempty"`
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
