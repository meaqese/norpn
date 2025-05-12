package repository

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/meaqese/norpn/internal/orch/domain"
	"testing"
)

func TestUserRepo(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo, err := NewUserRepo(db)
	if err != nil {
		t.Fatal(repo)
	}

	testUser := domain.User{
		Login:    "testuser",
		Password: "test",
	}

	userId, err := repo.Add(testUser)
	if err != nil {
		t.Fatal(err)
	}

	user, err := repo.GetByLogin(testUser.Login)
	if err != nil {
		t.Fatal(err)
	}

	if user.ID != userId {
		t.Fatalf("User ID should be %d, not %d", userId, user.ID)
	}

	if user.Password == testUser.Password {
		t.Fatal("User password should be crypted")
	}

	_, err = repo.GetByLogin("nonexist")
	if err == nil {
		t.Fatal("Found user that not exists")
	}
}
