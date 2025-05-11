package services

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/meaqese/norpn/internal/orch/domain"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	users     domain.UserRepository
	jwtSecret string
}

func NewAuth(userRepo domain.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		users:     userRepo,
		jwtSecret: jwtSecret,
	}
}

func comparePasswordAndHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (a *AuthService) Register(user domain.User) (int64, error) {
	return a.users.Add(user)
}

func (a *AuthService) Login(user domain.User) (string, error) {
	userInDB, err := a.users.GetByLogin(user.Login)
	if err != nil {
		return "", err
	}

	err = comparePasswordAndHash(user.Password, userInDB.Password)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"sub": userInDB.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(a.jwtSecret))
}

func (a *AuthService) Parse(token string) (int64, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) { return []byte(a.jwtSecret), nil })
	if err != nil {
		return 0, domain.ErrInvalidToken
	}

	return int64(t.Claims.(jwt.MapClaims)["sub"].(float64)), nil
}
