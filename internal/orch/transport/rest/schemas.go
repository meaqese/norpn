package rest

type Expression struct {
	ID     string  `json:"ID"`
	Status string  `json:"status"`
	Result float64 `json:"result,omitempty"`
	Reason string  `json:"reason,omitempty"`
}

type RequestExpression struct {
	Expression string `json:"expression"`
}

type ResponseExpression struct {
	ID string `json:"id"`
}

type ResponseExpressions struct {
	Expressions []*Expression `json:"expressions"`
}
