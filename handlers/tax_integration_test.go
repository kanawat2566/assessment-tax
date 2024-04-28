//go:build integration

package handlers_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	ct "github.com/kanawat2566/assessment-tax/constants"
	md "github.com/kanawat2566/assessment-tax/model"
	models "github.com/kanawat2566/assessment-tax/model"
	"github.com/stretchr/testify/assert"
)

type Response struct {
	*http.Response
	err error
}

func clientRequest(method, url string, body io.Reader) *Response {

	req, _ := http.NewRequest(method, url, body)
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

type TaxCalculationTestCase struct {
	name       string
	request    string
	statusCode int
	response   md.TaxResponse
}

func TestPostTaxCalulation_Integration(t *testing.T) {

	testCases := []TaxCalculationTestCase{
		{
			name: "Valid Tax Calucations",
			request: `{
							"totalIncome": 500000.0,
							"wht": 0.0,
							"allowances": [
							{
								"allowanceType": "donation",
								"amount": 0.0
							}
							]
						}`,
			statusCode: http.StatusOK,
			response:   md.TaxResponse{Tax: 29000},
		},
		{
			name: "Invalid Tax Calucations",
			request: `{
							"totalIncome1": 500000.0,
							"wht": 0.0,
							"allowances": [
							{
								"allowanceType": "donation",
								"amount": 0.0
							}
							]
						}`,
			statusCode: http.StatusBadRequest,
			response:   models.TaxResponse{},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			var result models.TaxResponse
			respRes := clientRequest("POST", uri("tax/calculations"), strings.NewReader(tc.request))

			err := respRes.Decode(&result)
			assert.Nil(t, err)
			assert.Equal(t, tc.statusCode, respRes.StatusCode)
			assert.Equal(t, tc.response.Tax, result.Tax)

		})
	}
}

type ConfigDeductionTestCase struct {
	name       string
	request    string
	url        string
	statusCode int
	response   md.DeductResponse
}

func TestAdminSetDedutions_Integration(t *testing.T) {
	testCases := []ConfigDeductionTestCase{
		{
			name: "Valid Tax Calucations",
			url:  "admin/deductions/personal",
			request: `{
						"amount": 60000.0
					  }`,
			statusCode: http.StatusOK,
			response:   md.DeductResponse{PersonalDeduction: 60000},
		},
		{
			name: "Invalid Tax Calucations",
			url:  "admin/deductions/personal",
			request: `{
						"amountx": 50000.0
						}`,
			statusCode: http.StatusBadRequest,
			response:   md.DeductResponse{PersonalDeduction: 0},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			var result models.DeductResponse

			body := strings.NewReader(tc.request)
			req, _ := http.NewRequest("POST", uri(tc.url), body)

			req.SetBasicAuth(ct.UserAuth, ct.PassAuth)
			req.Header.Add("Content-Type", "application/json")
			req.Close = true

			client := http.Client{}

			respRes, err := client.Do(req)
			if err != nil {
				fmt.Printf("Error making request %v\n", err)
			}

			defer respRes.Body.Close()
			err = json.NewDecoder(respRes.Body).Decode(&result)
			assert.Nil(t, err)
			assert.Equal(t, tc.statusCode, respRes.StatusCode)
			assert.Equal(t, tc.response, result)

		})
	}
}
