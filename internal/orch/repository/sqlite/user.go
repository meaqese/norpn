package repository

import (
	"database/sql"
	"github.com/meaqese/norpn/internal/orch/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	ID int64
	domain.User
}

func (u UserModel) generateHashOfPassword() (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) (*UserRepo, error) {
	createTable := `
		CREATE TABLE IF NOT EXISTS users (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    login TEXT NOT NULL,
		    password TEXT NOT NULL
		)
	`

	repo := &UserRepo{db: db}

	_, err := db.Exec(createTable)
	if err != nil {
		return repo, err
	}

	return repo, nil
}

func (ur *UserRepo) Add(user UserModel) (int64, error) {
	q := "INSERT INTO user (login, password) VALUES ($1, $2)"

	passwordHash, err := user.generateHashOfPassword()
	if err != nil {
		return 0, err
	}

	res, err := ur.db.Exec(q, user.Login, passwordHash)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (ur *UserRepo) GetByLogin(login string) (UserModel, error) {
	var user UserModel

	q := "SELECT id, login, password FROM users WHERE login = $1"
	err := ur.db.QueryRow(q, login).Scan(user.ID, user.Login, user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}
