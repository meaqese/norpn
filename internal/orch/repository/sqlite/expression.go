package repository

import (
	"database/sql"
	"errors"
	"github.com/meaqese/norpn/internal/orch/domain"
)

type ExpressionModel struct {
	ID         int64
	Expression string
	Status     string
	Result     sql.NullFloat64
	Reason     sql.NullString
	UserID     int64
}

func (em ExpressionModel) toDomain() *domain.Expression {
	exp := &domain.Expression{
		ID: em.ID, Expression: em.Expression,
		Status: em.Status, UserID: em.UserID,
	}

	if em.Result.Valid {
		exp.Result = em.Result.Float64
	}

	if em.Reason.Valid {
		exp.Reason = em.Reason.String
	}

	return exp
}

type ExpressionRepo struct {
	db *sql.DB
}

func NewExpressionRepo(db *sql.DB) (*ExpressionRepo, error) {
	createTable := `
		CREATE TABLE IF NOT EXISTS expressions (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    expression TEXT,
		    status TEXT,
		    result FLOAT,
		    reason TEXT,
		    user_id INTEGER,
		    
		    FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`

	repo := &ExpressionRepo{db: db}

	_, err := db.Exec(createTable)
	if err != nil {
		return repo, err
	}

	return repo, nil
}

func (e *ExpressionRepo) Add(expression domain.Expression) (int64, error) {
	if expression.UserID < 1 {
		return 0, errors.New("user id must be greater than zero")
	}

	query := "INSERT INTO expressions (expression, status, user_id) VALUES ($1, $2, $3)"
	result, err := e.db.Exec(query, expression.Expression, "processing", expression.UserID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (e *ExpressionRepo) Update(expression domain.Expression) error {
	q := "UPDATE expressions SET status = $1, result = $2, reason = $3 WHERE id = $4"
	_, err := e.db.Exec(q, expression.Status, expression.Result, expression.Reason, expression.ID)
	if err != nil {
		return err
	}

	return nil
}

func (e *ExpressionRepo) GetById(id int64) (*domain.Expression, error) {
	q := "SELECT id, expression, status, result, reason, user_id FROM expressions WHERE id = $1"

	exp := ExpressionModel{}
	err := e.db.QueryRow(q, id).Scan(&exp.ID, &exp.Expression, &exp.Status, &exp.Result, &exp.Reason, &exp.UserID)
	if err != nil {
		return &domain.Expression{}, err
	}

	return exp.toDomain(), nil
}

func (e *ExpressionRepo) GetAll(userId int64) ([]*domain.Expression, error) {
	q := "SELECT id, expression, status, result, reason, user_id FROM expressions WHERE user_id = $1"

	var expressions []*domain.Expression
	rows, err := e.db.Query(q, userId)
	if err != nil {
		return expressions, err
	}
	defer rows.Close()

	for rows.Next() {
		exp := ExpressionModel{}
		err = rows.Scan(&exp.ID, &exp.Expression, &exp.Status, &exp.Result, &exp.Reason, &exp.UserID)
		if err != nil {
			return expressions, err
		}
		expressions = append(expressions, exp.toDomain())
	}

	return expressions, nil
}
