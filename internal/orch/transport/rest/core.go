package rest

import (
	"crypto/md5"
	"encoding/hex"
	conf "github.com/meaqese/norpn/internal/orch/config"
	"github.com/meaqese/norpn/internal/orch/norpn"
	"sync"
)

type Core struct {
	mu *sync.Mutex

	config          *conf.Config
	expressionStore map[string]*Expression

	calculator *norpn.Calculator
}

func New(config *conf.Config) *Core {
	handler := &Core{
		config:          config,
		expressionStore: make(map[string]*Expression),
		calculator:      norpn.New(config),
		mu:              &sync.Mutex{},
	}
	return handler
}

func (c *Core) GenerateHash(expression string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	md5Hash := md5.Sum([]byte(expression))
	return hex.EncodeToString(md5Hash[:])
}

func (c *Core) AddExpression(expHash string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	newExpression := &Expression{
		ID:     expHash,
		Status: "processing",
	}

	c.expressionStore[newExpression.ID] = newExpression
}

func (c *Core) UpdateExpression(expHash string, newStatus string, result float64, reason string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expression := c.expressionStore[expHash]
	expression.Status = newStatus
	expression.Result = result

	if reason != "" {
		expression.Reason = reason
	}
}

func (c *Core) StartCalc(expression, expHash string) {
	result, err := c.calculator.Calc(expression)
	if err != nil {
		c.UpdateExpression(expHash, "error", 0, err.Error())
		return
	}

	c.UpdateExpression(expHash, "completed", result, "")
}
