package rest

import (
	"bytes"
	"encoding/json"
	"github.com/meaqese/norpn/internal/orch/norpn"
	"net/http/httptest"
	"strings"
	"testing"
)

const ApiEndpoint = "http://norpn.io/api/v1/calculate"

func TestInvalidJson(t *testing.T) {
	buf := make([]byte, 100)
	buf = append(buf, []byte("{\"expression\":}")...)
	buffer := bytes.NewReader(buf)

	req := httptest.NewRequest("POST", ApiEndpoint, buffer)
	w := httptest.NewRecorder()

	CalcHandler(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	decoder := json.NewDecoder(w.Body)

	resResult := &ResponseError{}
	err := decoder.Decode(resResult)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 422 || strings.ToLower(resResult.Error) != ErrExpressionIsInvalid.Error() {
		t.Fatalf("when sending invalid json, should be error and 422 status code")
	}
}

func testingCase(tc norpn.TestCase, t *testing.T) {
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(RequestData{Expression: tc.Expression})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", ApiEndpoint, &buffer)
	w := httptest.NewRecorder()

	CalcHandler(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	decoder := json.NewDecoder(w.Body)

	if resp.StatusCode == 200 {
		if tc.ShouldFail {
			t.Fatalf("test should fail")
		}

		resResult := &ResponseResult{}
		err = decoder.Decode(resResult)
		if err != nil {
			t.Fatal(err)
		}

		if resResult.Result != tc.Expected {
			t.Fatalf("not valid return")
		}
	} else {
		if !tc.ShouldFail {
			t.Fatalf("success test not should fail")
		}
	}
}

func TestCalcHandler(t *testing.T) {
	cases := norpn.GetTestCases()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			testingCase(tc, t)
		})
	}
}
