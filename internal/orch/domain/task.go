package domain

type Task struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     rune    `json:"operation"`
	OperationTime int     `json:"operation_time"`
}
