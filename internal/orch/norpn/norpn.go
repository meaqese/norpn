package norpn

import (
	"errors"
	"fmt"
	conf "github.com/meaqese/norpn/internal/orch/config"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

type Task struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     rune    `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

type TaskResult struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}

type Calculator struct {
	Mu     *sync.Mutex
	config *conf.Config

	TaskQueue          []*Task
	TaskResultChannels map[string]*chan float64

	TestMode bool
}

func New(config *conf.Config) *Calculator {
	return &Calculator{
		TaskQueue:          make([]*Task, 0),
		TaskResultChannels: make(map[string]*chan float64),
		Mu:                 &sync.Mutex{},
		config:             config,
		TestMode:           false,
	}
}

func (c *Calculator) DequeueTask() *Task {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	if len(c.TaskQueue) > 0 {
		task := c.TaskQueue[0]
		c.TaskQueue = c.TaskQueue[1:]

		return task
	}
	return nil
}

func (c *Calculator) EnqueueTask(task *Task) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	c.TaskQueue = append(c.TaskQueue, task)

	channel := make(chan float64, 1)
	c.TaskResultChannels[task.ID] = &channel
}

func (c *Calculator) generateID() string {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	for {
		generatedID := strconv.FormatInt(rand.Int63(), 10)
		if _, ok := c.TaskResultChannels[generatedID]; !ok {
			return generatedID
		}
	}
}

func (c *Calculator) getOperationTime(operator rune) int {
	switch operator {
	case '+':
		return c.config.TimeAdditionMs
	case '-':
		return c.config.TimeSubtractionMs
	case '*':
		return c.config.TimeMultiplicationsMs
	case '/':
		return c.config.TimeDivisionsMs
	}

	return 0
}

func isOperator(r rune) bool {
	return r == '/' || r == '+' || r == '-' || r == '*'
}

func GetPrecedence(r rune) int {
	switch r {
	case '/', '*':
		return 2
	case '+', '-':
		return 1
	}

	return 0
}

func getLastValue(stack []float64, secondToLast bool) (float64, error) {
	rangeOutError := errors.New("index out of range")

	if len(stack) < 1 {
		return 0, rangeOutError
	}
	if secondToLast {
		if len(stack) < 2 {
			return 0, rangeOutError
		}
		return stack[len(stack)-2], nil
	}
	return stack[len(stack)-1], nil
}

func getLastOperator(stack []rune) rune {
	return stack[len(stack)-1]
}

func (c *Calculator) solve(valuesStack []float64, operatorsStack []rune) ([]float64, []rune, error) {
	val2, err1 := getLastValue(valuesStack, false)
	val1, err2 := getLastValue(valuesStack, true)

	if err1 != nil || err2 != nil {
		return nil, nil, err1
	}

	valuesStack = valuesStack[:len(valuesStack)-2]

	lastOperator := getLastOperator(operatorsStack)
	operatorsStack = operatorsStack[:len(operatorsStack)-1]

	if lastOperator == '/' && val2 == 0 {
		return valuesStack, operatorsStack, errors.New("division to 0")
	}

	id := c.generateID()

	c.EnqueueTask(&Task{
		ID:            id,
		Arg1:          val1,
		Arg2:          val2,
		Operation:     lastOperator,
		OperationTime: c.getOperationTime(lastOperator),
	})

	select {
	case result := <-*c.TaskResultChannels[id]:
		valuesStack = append(valuesStack, result)
		close(*c.TaskResultChannels[id])
		delete(c.TaskResultChannels, id)
	case <-time.After(300 * time.Second):
		return valuesStack, operatorsStack, errors.New("timeout error")
	}

	return valuesStack, operatorsStack, nil
}

func (c *Calculator) solveSimpleExpression(expression string) (string, error) {
	var valuesStack []float64
	var operatorsStack []rune
	var err error

	for i := 0; i < len(expression); i++ {
		symbol := rune(expression[i])
		if unicode.IsDigit(symbol) || (isOperator(symbol) && len(valuesStack) == len(operatorsStack)) {
			var numStr string
			var floatNum float64

			j := i
			for len(expression)-1 >= j && !isOperator(rune(expression[j])) {
				numStr += string(expression[j])
				j++
			}
			i = j - 1

			floatNum, err = strconv.ParseFloat(numStr, 64)
			if err != nil {
				return "", errors.New("parsing error")
			}
			valuesStack = append(valuesStack, floatNum)
		} else if isOperator(symbol) {
			for len(operatorsStack) > 0 && GetPrecedence(getLastOperator(operatorsStack)) >= GetPrecedence(symbol) {
				valuesStack, operatorsStack, err = c.solve(valuesStack, operatorsStack)
				if err != nil {
					return "", err
				}
			}

			operatorsStack = append(operatorsStack, symbol)
		}
	}

	for len(operatorsStack) > 0 {
		valuesStack, operatorsStack, err = c.solve(valuesStack, operatorsStack)
		if err != nil {
			return "", err
		}
	}

	if len(valuesStack) == 0 {
		return "", errors.New("no one value to return")
	}

	return fmt.Sprintf("%f", valuesStack[0]), nil
}

func (c *Calculator) Calc(expression string) (float64, error) {
	expression = strings.Replace(expression, " ", "", -1)

	rechecks := 1
	for rechecks > 0 {
		for i := 0; i < len(expression); i++ {
			currentSym := expression[i]
			if currentSym == '(' {
				for j := i + 1; j < len(expression); j++ {
					if expression[j] == '(' {
						rechecks += 1

						break
					}
					if expression[j] == ')' {
						bracketExpression := expression[i+1 : j]

						solvedBracketExp, err := c.solveSimpleExpression(bracketExpression)
						if err != nil {
							return 0, err
						}

						expression = strings.Replace(expression, "("+bracketExpression+")", solvedBracketExp, -1)
						break
					}
				}
			}
		}
		rechecks -= 1
	}

	solvedExp, err := c.solveSimpleExpression(expression)
	if err != nil {
		return 0, err
	}

	res, _ := strconv.ParseFloat(solvedExp, 64)

	return res, nil
}
