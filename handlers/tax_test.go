package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator"
	"github.com/kanawat2566/assessment-tax/handlers"
	"github.com/kanawat2566/assessment-tax/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type MockTaxService struct {
	taxResp models.TaxResponse
	err     error
}

func (m *MockTaxService) TaxCalculations(taxRequest *models.TaxRequest) (models.TaxResponse, error) {
	return m.taxResp, m.err
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

func RequestBody(taxRequest models.TaxRequest) *bytes.Reader {
	jsonData, _ := json.Marshal(taxRequest)
	return bytes.NewReader(jsonData)
}
