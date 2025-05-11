package domain

type Expression struct {
	ID         int64
	Expression string
	Status     string
	Result     float64
	Reason     string
	UserID     int64
}
