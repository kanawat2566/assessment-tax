//go:build integration

package handlers_test

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	models "github.com/kanawat2566/assessment-tax/model"
	"github.com/stretchr/testify/assert"
)

type Response struct {
	*http.Response
	err error
}

func clientRequest(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	//req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(v)
}

func uri(paths ...string) string {
	baseURL := os.Getenv("TEST_URL")

	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	if paths == nil {
		return baseURL
	}
	return baseURL + "/" + strings.Join(paths, "/")
}

func TestPostTaxCalulation_Invalid(t *testing.T) {

	var result models.TaxResponse

	res := clientRequest("POST", uri("tax/calculations"), nil)
	err := res.Decode(&result)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, res.StatusCode)

}

func TestPostTaxCalulation_valid(t *testing.T) {

	var result models.TaxResponse

	res := clientRequest("POST", uri("tax/calculations"), strings.NewReader(
		`{
		"totalIncome": 500000.0,
		"wht": 0.0,
		"allowances": [
		  {
			"allowanceType": "donation",
			"amount": 0.0
		  }
		]
	  }
	`))
	err := res.Decode(&result)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)

}
