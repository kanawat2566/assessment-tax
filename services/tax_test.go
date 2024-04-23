package services_test

import (
	"errors"
	"testing"

	"github.com/kanawat2566/assessment-tax/constants"
	models "github.com/kanawat2566/assessment-tax/model"
	"github.com/kanawat2566/assessment-tax/repository"
	"github.com/kanawat2566/assessment-tax/services"
	"github.com/stretchr/testify/assert" // assuming you use testify for assertions
)

type MockTaxRepository struct {
	taxRates   []*repository.IncomeTaxRates
	allowances map[string]repository.Allowances
}

func (m *MockTaxRepository) GetTaxRates() (res []*repository.IncomeTaxRates, err error) {
	return m.taxRates, err
}

func (m *MockTaxRepository) GetAllowanceConfig(allowanceType string) (repository.Allowances, error) {
	return m.allowances[allowanceType], nil
}

type TestCase struct {
	name     string
	request  models.TaxRequest
	expected models.TaxResponse
}

var mockRepo = &MockTaxRepository{
	taxRates: []*repository.IncomeTaxRates{
		{IncomeLevel: "0-150,000", MinIncome: 0, MaxIncome: 150000, TaxRate: 0},
		{IncomeLevel: "150,001-500,000", MinIncome: 150001, MaxIncome: 500000, TaxRate: 10},
		{IncomeLevel: "500,001-1,000,000", MinIncome: 500001.00, MaxIncome: 1000000.00, TaxRate: 15},
		{IncomeLevel: "1,000,001-2,000,000", MinIncome: 1000001.00, MaxIncome: 2000000.00, TaxRate: 20},
		{IncomeLevel: "2,000,001 ขึ้นไป", MinIncome: 2000001.00, MaxIncome: 99999999999999.00, TaxRate: 35},
	},
	allowances: map[string]repository.Allowances{
		constants.Personal: {LimitAmt: 60000},
	},
}

func StoryUseCases() []TestCase {
	cs := []TestCase{
		{
			name: "given input income total should payment tax",
			request: models.TaxRequest{
				TotalIncome: 500000,
			},
			expected: models.TaxResponse{
				Tax: 29000,
			},
		},
		{
			name: "given input income total and deducting WHT should payment tax",
			request: models.TaxRequest{
				TotalIncome: 500000,
				WHT:         25000.0,
			},
			expected: models.TaxResponse{
				Tax: 4000,
			},
		},
		{
			name: "given input income total and deducting WHT should return tax refund",
			request: models.TaxRequest{
				TotalIncome: 150000,
				WHT:         2000,
			},
			expected: models.TaxResponse{
				Tax:       0,
				TaxRefund: 2000,
			},
		},
	}
	return cs
}

func TestCalculateTax_ValidInput(t *testing.T) {

	cases := StoryUseCases()

	for _, tc := range cases {

		t.Run(tc.name, func(t *testing.T) {
			serv := services.NewServices(mockRepo)
			rep, err := serv.TaxCalculations(&tc.request)

			// Assertions
			assert.Nil(t, err, "Error should be nil for valid inputs")
			assert.Equal(t, tc.expected.Tax, rep.Tax, "Calculated tax should match")
		})
	}

}

func TestCalculateTax_InvalidIncome(t *testing.T) {

	mockRepo := &MockTaxRepository{}
	taxRequest := &models.TaxRequest{
		TotalIncome: -1,
	}

	taxService := services.NewServices(mockRepo)
	taxResponse, err := taxService.TaxCalculations(taxRequest)

	// Assertions
	assert.NotNil(t, err, "Error should not be nil for invalid income")
	assert.EqualError(t, err, constants.ErrMessageThenZero, "Error message should match income")
	assert.Zero(t, taxResponse)
}
func TestCalculateTax_InvalidWHT(t *testing.T) {

	mockRepo := &MockTaxRepository{}
	taxRequest := &models.TaxRequest{
		TotalIncome: 200,
		WHT:         -1,
	}

	taxService := services.NewServices(mockRepo)
	taxResponse, err := taxService.TaxCalculations(taxRequest)

	// Assertions
	assert.NotNil(t, err, "Error should not be nil for invalid income")
	assert.EqualError(t, err, constants.ErrMesssageWhtInvalid, "Error message should match WHT")
	assert.Zero(t, taxResponse)
}

func TestCalculateTax_InvalidWHTMore(t *testing.T) {

	mockRepo := &MockTaxRepository{}
	taxRequest := &models.TaxRequest{
		TotalIncome: 200,
		WHT:         200,
	}

	taxService := services.NewServices(mockRepo)
	taxResponse, err := taxService.TaxCalculations(taxRequest)

	// Assertions
	assert.NotNil(t, err, "Error should not be nil for invalid income")
	assert.EqualError(t, err, constants.ErrMesssageWhtInvalid, "Error message should match WHT")
	assert.Zero(t, taxResponse)
}

type MockTaxRepository2 struct {
	taxRates   []*repository.IncomeTaxRates
	allowances map[string]repository.Allowances
}

func (m *MockTaxRepository2) GetTaxRates() (res []*repository.IncomeTaxRates, err error) {
	return m.taxRates, errors.New(constants.ErrMessageInternal)
}

func (m *MockTaxRepository2) GetAllowanceConfig(allowanceType string) (repository.Allowances, error) {
	return m.allowances[allowanceType], nil
}

func TestCalculateTax_GetTaxRatesError(t *testing.T) {

	mockRepo := &MockTaxRepository2{}
	taxRequest := &models.TaxRequest{
		TotalIncome: 100,
		WHT:         0,
	}

	taxService := services.NewServices(mockRepo)
	_, err := taxService.TaxCalculations(taxRequest)

	// Assertions
	assert.EqualError(t, err, constants.ErrMessageInternal, "Error message should match")
}
