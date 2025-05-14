package rest

import (
	"github.com/meaqese/norpn/internal/orch/services"
	"net/http"
)

type Core struct {
	calculator *services.Calculator
	auth       *services.AuthService
}

func New(calculator *services.Calculator, auth *services.AuthService) (*http.ServeMux, *Core) {
	handler := &Core{
		calculator: calculator,
		auth:       auth,
	}

	fs := http.FileServer(http.Dir("./web"))

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", cors(withAuth(auth, handler.HandleExpression)))
	mux.HandleFunc("/api/v1/register", cors(handler.Register))
	mux.HandleFunc("/api/v1/login", cors(handler.Login))
	mux.HandleFunc("/api/v1/expressions/{id}", cors(withAuth(auth, handler.HandleGetExpression)))
	mux.HandleFunc("/api/v1/expressions", cors(withAuth(auth, handler.HandleGetExpressions)))
	mux.Handle("/", fs)

	return mux, handler
}
