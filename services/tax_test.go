package services_test

import (
	"testing"

	"github.com/kanawat2566/assessment-tax/constants"
	"github.com/kanawat2566/assessment-tax/models"
	"github.com/kanawat2566/assessment-tax/repository"
	"github.com/kanawat2566/assessment-tax/services"
	"github.com/stretchr/testify/assert" // assuming you use testify for assertions
)

type MockTaxRepository struct {
	taxRates   []*repository.IncomeTaxRates
	allowances map[string]repository.Allowances
}

func (m *MockTaxRepository) GetTaxRates() ([]*repository.IncomeTaxRates, error) {
	return m.taxRates, nil
}

// GetConfigTax mocks the GetConfigTax method of the TaxRepository interface.
func (m *MockTaxRepository) GetAllowanceConfig(allowanceType string) (repository.Allowances, error) {
	return m.allowances[allowanceType], nil
}

func TestCalculateTax_ValidInput(t *testing.T) {
	// Create mock repository
	mockRepo := &MockTaxRepository{
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

	// Create valid tax request
	taxRequest := &models.TaxRequest{
		TotalIncome: 500000,
		WHT:         0,
		Allowances: []models.Allowance{
			{AllowanceType: constants.Donation, Amount: 0},
		},
	}

	// Create tax service
	taxService := services.NewServices(mockRepo)

	// Call the function to be tested
	taxResponse, err := taxService.TaxCalculations(taxRequest)

	// Assertions
	assert.Nil(t, err, "Error should be nil for valid inputs")
	assert.Equal(t, 29000.0, taxResponse.Tax, "Calculated tax should match")
}

func TestCalculateTax_InvalidIncome(t *testing.T) {
	// Create mock repository (not used in this test)
	mockRepo := &MockTaxRepository{}

	// Create invalid tax request (negative income)
	taxRequest := &models.TaxRequest{
		TotalIncome: -10000,
		WHT:         -1,
		Allowances:  []models.Allowance{},
	}

	// Create tax service
	taxService := services.NewServices(mockRepo)

	// Call the function to be tested
	taxResponse, err := taxService.TaxCalculations(taxRequest)

	// Assertions
	assert.NotNil(t, err, "Error should not be nil for invalid income")
	assert.EqualError(t, err, constants.ErrMessageThenZero, "Error message should match")

	assert.Zero(t, taxResponse) // taxResponse should be empty
}
