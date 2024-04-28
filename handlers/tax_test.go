package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	ct "github.com/kanawat2566/assessment-tax/constants"
	"github.com/kanawat2566/assessment-tax/handlers"
	models "github.com/kanawat2566/assessment-tax/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type MockTaxService struct {
	taxResp    models.TaxResponse
	taxErr     error
	deductResp ct.Deduction
	deductErr  error
	taxCsv     []models.Taxes
	fileCsv    string
	taxCsvErr  error
}

func (m *MockTaxService) TaxCalculations(taxRequest models.TaxRequest) (models.TaxResponse, error) {
	return m.taxResp, m.taxErr
}

func (m *MockTaxService) SetAdminDeductions(req ct.Deduction) (ct.Deduction, error) {
	return m.deductResp, m.deductErr
}
func (m *MockTaxService) TaxCalFromCsv(taxRequest []models.TaxRequest) ([]models.Taxes, error) {
	return m.taxCsv, m.taxCsvErr
}
func TestCalculationsHandler_ValidRequest(t *testing.T) {
	// Create mock service
	mockService := &MockTaxService{
		taxResp: models.TaxResponse{Tax: 29000.0},
	}

	// Create handler with mock service
	handler := handlers.NewHandler(mockService)

	// Create a valid tax request
	taxRequest := models.TaxRequest{
		TotalIncome: 500000,
		WHT:         0,
	}

	// Create a request object
	e := echo.New()
	//e.Validator = &handlers.CustomValidator{Validator: validator.New()}
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
	assert.Nil(t, json.Unmarshal(rec.Body.Bytes(), &response), "Response should be unmarshallable xxxx")
	assert.Equal(t, mockService.taxResp, response, "Response should match mock service response")
	assert.Equal(t, mockService.taxResp.Tax, response.Tax, "Response should match mock service response")
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
	//e.Validator = &handlers.CustomValidator{Validator: validator.New()}
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

type taxResponse struct {
	Taxes []models.Taxes `json:"taxes"`
}

func TestCalFromUploadCsvHandler(t *testing.T) {
	// Mock CSV data
	validCsv := "totalIncome,wht,donation\n500000,0,0\n600000,40000,20000\n750000,50000,15000\n"
	expected := []models.Taxes{
		{Tax: 29000, TotalIncome: 500000},
		{Tax: 29000, TotalIncome: 600000, TaxRefund: 2000},
		{Tax: 11250, TotalIncome: 750000},
	}
	// Valid CSV test case
	t.Run("ValidCSV", func(t *testing.T) {

		mockService := &MockTaxService{
			fileCsv: validCsv,
			taxCsv:  expected,
		}

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		writer.WriteField("taxFile", "taxes.csv")
		part, _ := writer.CreateFormFile("taxFile", "taxes.csv")
		part.Write([]byte(validCsv))
		writer.Close()

		e := echo.New()
		//e.Validator = &handlers.CustomValidator{Validator: validator.New()}
		req := httptest.NewRequest(http.MethodPost, "/", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("tax/calculations/:uploadType")
		ctx.SetParamNames("uploadType")
		ctx.SetParamValues("upload-csv")
		ctx.FormFile("taxFile")

		handler := handlers.NewHandler(mockService)

		err := handler.CalFromUploadCsvHandler(ctx)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response taxResponse
		assert.Nil(t, json.Unmarshal(rec.Body.Bytes(), &response), "Response should be unmarshallable")
		assert.Equal(t, expected, response.Taxes, "Response should match mock service response")

	})
}

func (m *MockTaxService) FormFile(key string) (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewBufferString(m.fileCsv)), nil
}
