package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator"
	ct "github.com/kanawat2566/assessment-tax/constants"
	"github.com/kanawat2566/assessment-tax/handlers"
	models "github.com/kanawat2566/assessment-tax/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type MockTaxService struct {
	taxResp    models.TaxResponse
	deductResp ct.Deduction
	err        error
}

func (m *MockTaxService) TaxCalculations(taxRequest *models.TaxRequest) (models.TaxResponse, error) {
	return m.taxResp, m.err
}

func (m *MockTaxService) SetAdminDeductions(req ct.Deduction) (ct.Deduction, error) {
	return m.deductResp, m.err
}

func TestCalculationsHandler_ValidRequest(t *testing.T) {
	// Create mock service
	mockService := &MockTaxService{
		taxResp: models.TaxResponse{Tax: 1000.0},
	}

	// Create handler with mock service
	handler := handlers.NewHandler(mockService)

	// Create a valid tax request
	taxRequest := models.TaxRequest{
		TotalIncome: 200000,
		WHT:         0,
	}

	// Create a request object
	e := echo.New()
	e.Validator = &handlers.CustomValidator{Validator: validator.New()}
	req := httptest.NewRequest(http.MethodPost, "/tax/calulations", RequestBody(taxRequest))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// Call the handler function
	err := handler.CalculationsHandler(ctx)

	// Assertions
	assert.Nil(t, err, "Error should be nil for valid request")
	assert.Equal(t, http.StatusOK, rec.Code, "HTTP status code should be 200")
	var response models.TaxResponse
	assert.Nil(t, json.Unmarshal(rec.Body.Bytes(), &response), "Response should be unmarshallable")
	assert.Equal(t, mockService.taxResp, response, "Response should match mock service response")
}

func RequestBody(r interface{}) *bytes.Reader {
	jsonData, _ := json.Marshal(r)
	return bytes.NewReader(jsonData)
}

func TestAdminDeductionHandler_ValidRequest(t *testing.T) {
	// Create mock service
	mockService := &MockTaxService{
		deductResp: ct.Deduction{Name: "PersonalDeduction", Amount: 70000.0},
	}

	// expected
	ep := models.DeductResponse{PersonalDeduction: 70000.00}

	// Create handler with mock service
	handler := handlers.NewHandler(mockService)

	// Create a valid tax request
	rq := models.DeductRequest{
		Amount: 70000.0,
	}

	// Create a request object
	e := echo.New()
	e.Validator = &handlers.CustomValidator{Validator: validator.New()}
	req := httptest.NewRequest(http.MethodPost, "/", RequestBody(rq))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("admin/deductions/:type")
	ctx.SetParamNames("type")
	ctx.SetParamValues(ct.Personal)

	// Call the handler function
	err := handler.Deductions(ctx)

	// Assertions
	assert.Nil(t, err, "Error should be nil for valid request")
	assert.Equal(t, http.StatusOK, rec.Code, "HTTP status code should be 200")
	var response models.DeductResponse
	assert.Nil(t, json.Unmarshal(rec.Body.Bytes(), &response), "Response should be unmarshallable")
	assert.Equal(t, ep.PersonalDeduction, response.PersonalDeduction, "Response should match mock service response")
	assert.Equal(t, ep.KReceipt, response.KReceipt, "Response should match mock service response")
}
