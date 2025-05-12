package services

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/meaqese/norpn/internal/orch/domain"
	repository "github.com/meaqese/norpn/internal/orch/repository/sqlite"
	"testing"
)

func TestAuthService(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")

	userRepo, _ := repository.NewUserRepo(db)

	authSvc := NewAuth(userRepo, "secret")

	_, err := authSvc.Login(domain.User{Login: "user", Password: "pass"})
	if err == nil {
		t.Fatal("User is not exist, should return error")
	}

	testUser := domain.User{
		Login:    "user",
		Password: "pass",
	}

	userId, err := authSvc.Register(testUser)
	if err != nil {
		t.Fatal(err)
	}

	token, err := authSvc.Login(testUser)
	if err != nil {
		t.Fatal(err)
	}

	parsedUserId, err := authSvc.Parse(token)
	if err != nil {
		t.Fatal(err)
	}

	if parsedUserId != userId {
		t.Fatalf("User ID in JWT sub should be %d, not %d", userId, parsedUserId)
	}
}
