package repository_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	ct "github.com/kanawat2566/assessment-tax/constants"
	"github.com/kanawat2566/assessment-tax/repository"
	"github.com/stretchr/testify/assert"
)

func TestGetTaxRates_Error(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "Error creating mock DB")
	defer db.Close()

	mock.ExpectQuery("SELECT * FROM income_tax_rates").WillReturnError(errors.New(ct.ErrMsgDatabaseError))

	repo := repository.New(db)

	// Call GetTaxRates
	_, err = repo.GetTaxRates()

	// Assertions
	assert.NotNil(t, err, "Error should not be nil for failed query")
	assert.EqualError(t, err, ct.ErrMsgDatabaseError, "Error message should match")
}

func TestGetTaxRates_Success(t *testing.T) {
	// Define expected tax rates
	expected := []*repository.IncomeTaxRates{
		{ID: 1, TaxRate: 0, IncomeLevel: "0 - 150,000", MinIncome: 0, MaxIncome: 150000},
		{ID: 2, TaxRate: 10, IncomeLevel: "150,001 - 500,000", MinIncome: 150001, MaxIncome: 500000},
	}

	// Create a mock database connection
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "Error creating mock DB")
	defer db.Close()

	// Configure mock query to return expected data
	mock.ExpectQuery(`
		SELECT 
		id, income_level, 
		min_income, max_income, 
		tax_rate 
		FROM income_tax_rates
		ORDER BY id
	`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "income_level", "min_income", "max_income", "tax_rate"}).
			AddRow(expected[0].ID, expected[0].IncomeLevel, expected[0].MinIncome, expected[0].MaxIncome, expected[0].TaxRate).
			AddRow(expected[1].ID, expected[1].IncomeLevel, expected[1].MinIncome, expected[1].MaxIncome, expected[1].TaxRate),
	)

	// Create a TaxRepository instance using the mock Postgres
	repo := repository.New(db)

	// Call GetTaxRates
	taxRates, err := repo.GetTaxRates()

	// Assertions
	assert.Nil(t, err, "Error should be nil for successful query")
	assert.Equal(t, expected, taxRates, "Tax rates should match")
	assert.Len(t, taxRates, 2, "Should have two tax rates")
}
