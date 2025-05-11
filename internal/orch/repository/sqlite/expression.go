package repository

import (
	"database/sql"
	"github.com/meaqese/norpn/internal/orch/domain"
)

type ExpressionModel struct {
	domain.Expression
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
