package repository

import (
	"database/sql"
	"github.com/meaqese/norpn/internal/orch/domain"
	"testing"
)

func TestExpressionRepo(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo, err := NewExpressionRepo(db)
	if err != nil {
		t.Fatal(repo)
	}

	userRepo, err := NewUserRepo(db)
	if err != nil {
		t.Fatal(userRepo)
	}

	testUser := domain.User{
		Login:    "test",
		Password: "test",
	}
	testExpression := domain.Expression{
		Expression: "2+2",
	}

	userId, err := userRepo.Add(testUser)

	_, err = repo.Add(testExpression)
	if err == nil {
		t.Fatal("Expression without user should throw error")
	}

	testExpression.UserID = userId
	expId, err := repo.Add(testExpression)
	if err != nil {
		t.Fatal(err)
	}

	expression, err := repo.GetById(expId)
	if err != nil {
		t.Fatal(err)
	}

	if expression.Expression != testExpression.Expression {
		t.Fatalf("Expression should be %s", testExpression.Expression)
	}

	expression.Status = "completed"
	expression.Result = 12

	err = repo.Update(*expression)
	if err != nil {
		t.Fatal(err)
	}

	expList, err := repo.GetAll(userId)
	if err != nil {
		t.Fatal(err)
	}

	if expList[0].Result != expression.Result {
		t.Fatalf("Expression result should be %f, not %f", expression.Result, expList[0].Result)
	}
}
