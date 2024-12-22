package rest

import (
	"encoding/json"
	"fmt"
	"github.com/meaqese/norpn/pkg/norpn"
	"log"
	"net/http"
	"os"
	"strings"
)

type RequestData struct {
	Expression string `json:"expression"`
}

type ResponseResult struct {
	Result float64 `json:"result"`
}

type ResponseError struct {
	Error       string `json:"error"`
	Description string `json:"description,omitempty'"`
}

func capitalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	requestData := &RequestData{}
	returnErr := func(err error, description error, statusCode int) {
		w.WriteHeader(statusCode)

		errMain := capitalize(err.Error())

		var errDesc string
		if description != nil {
			errDesc = description.Error()
			encoder.Encode(ResponseError{Error: errMain, Description: errDesc})
		} else {
			encoder.Encode(ResponseError{Error: errMain})
		}

		log.SetOutput(os.Stderr)
		log.Println(fmt.Sprintf("[ERROR] %s (%s). expression: \"%s\"", errMain, errDesc, requestData.Expression))
	}

	defer r.Body.Close()
	defer func() {
		if r := recover(); r != nil {
			returnErr(ErrInternalServerError, nil, http.StatusInternalServerError)
		}
	}()

	err := json.NewDecoder(r.Body).Decode(requestData)
	if err != nil {
		returnErr(ErrExpressionIsInvalid, ErrJsonValidation, http.StatusUnprocessableEntity)
		return
	}

	result, err := norpn.Calc(requestData.Expression)
	if err != nil {
		returnErr(ErrExpressionIsInvalid, err, http.StatusUnprocessableEntity)
		return
	}

	encoder.Encode(ResponseResult{Result: result})
	log.Println(fmt.Sprintf("accepted \"%s\", returned \"%.2f\"", requestData.Expression, result))
}
