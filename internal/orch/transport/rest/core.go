package rest

import (
	"github.com/meaqese/norpn/internal/orch/services"
)

type Core struct {
	calculator *services.Calculator
}

func New(calculator *services.Calculator) *Core {
	handler := &Core{
		calculator: calculator,
	}
	return handler
}
