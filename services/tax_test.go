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
	taxRates      []*repository.IncomeTaxRates
	allowances    map[string]repository.Allowances
	ratesErr      error
	allowancesErr error
}

func (m *MockTaxRepository) GetTaxRates() (res []*repository.IncomeTaxRates, err error) {
	return m.taxRates, m.ratesErr
}

func (m *MockTaxRepository) GetLimitAllowances(allowanceType string) (r repository.Allowances, err error) {
	return m.allowances[allowanceType], m.allowancesErr
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
		constants.Personal:  {LimitAmt: 60000, MinAmt: 10001, MaxAmt: 100000},
		constants.Donation:  {LimitAmt: 100000, MinAmt: 0, MaxAmt: 100000},
		constants.K_Receipt: {LimitAmt: 50000, MinAmt: 1, MaxAmt: 100000},
	},
}

func StoryUseCases() []TestCase {
	cs := []TestCase{
		{
			name: "given input income total should payment tax",
			request: models.TaxRequest{
				TotalIncome: 500000,
				WHT:         0,
				Allowances: []models.Allowance{
					{
						AllowanceType: constants.Donation,
						Amount:        0,
					}}},
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
		{
			name: "given input income total and allownce donation should payment tax",
			request: models.TaxRequest{
				TotalIncome: 500000,
				WHT:         0,
				Allowances: []models.Allowance{
					{
						AllowanceType: constants.Donation,
						Amount:        200000.0,
					}}},
			expected: models.TaxResponse{
				Tax: 19000,
			},
		},
		{
			name: "given input income total and allownce fix personal should payment tax",
			request: models.TaxRequest{
				TotalIncome: 500000,
				WHT:         0,
				Allowances: []models.Allowance{
					{
						AllowanceType: constants.Personal,
						Amount:        50000,
					},
				},
			},
			expected: models.TaxResponse{
				Tax: 30000,
			},
		},
		{
			name: "given input income total and deducting WHT with allownce should return tax refund",
			request: models.TaxRequest{
				TotalIncome: 500000,
				WHT:         25000,
				Allowances: []models.Allowance{
					{
						AllowanceType: constants.Donation,
						Amount:        100000,
					},
					{
						AllowanceType: constants.K_Receipt,
						Amount:        100000,
					},
				},
			},
			expected: models.TaxResponse{
				Tax:       0,
				TaxRefund: 11000,
			},
		},
		{
			name: "given input income total and deducting WHT with allownce should return tax level detail",
			request: models.TaxRequest{
				TotalIncome: 500000.0,
				WHT:         0.0,
				Allowances: []models.Allowance{
					{
						AllowanceType: constants.Donation,
						Amount:        200000.0,
					},
				},
			},
			expected: models.TaxResponse{
				Tax: 19000.0,
				TaxLevels: []models.TaxLevel{
					{Level: "0-150,000", Tax: 0.0},
					{Level: "150,001-500,000", Tax: 19000.0},
					{Level: "500,001-1,000,000", Tax: 0.0},
					{Level: "1,000,001-2,000,000", Tax: 0.0},
					{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
				},
			},
		},
	}
	return cs
}

func TestCalculateTax_Valid(t *testing.T) {

	cases := StoryUseCases()

	for _, tc := range cases {

		t.Run(tc.name, func(t *testing.T) {
			serv := services.NewServices(mockRepo)
			rep, err := serv.TaxCalculations(&tc.request)

			// Assertions
			assert.Nil(t, err, "Error should be nil for valid inputs")
			assert.Equal(t, tc.expected.Tax, rep.Tax, "Calculated tax should match")
			assert.Equal(t, tc.expected.TaxRefund, rep.TaxRefund, "Tax refund should match")
		})
	}
}

func TestCalculateTaxLevel_Valid(t *testing.T) {

	cases := []TestCase{{
		name: "given input income total and deducting WHT with allownce should return tax level detail",
		request: models.TaxRequest{
			TotalIncome: 500000.0,
			WHT:         0.0,
			Allowances: []models.Allowance{
				{
					AllowanceType: constants.Donation,
					Amount:        200000.0,
				},
			},
		},
		expected: models.TaxResponse{
			Tax: 19000.0,
			TaxLevels: []models.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 19000.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			},
		},
	},
	}

	for _, tc := range cases {

		t.Run(tc.name, func(t *testing.T) {
			serv := services.NewServices(mockRepo)
			rep, err := serv.TaxCalculations(&tc.request)

			// Assertions
			assert.Nil(t, err, "Error should be nil for valid inputs")
			assert.Equal(t, tc.expected.Tax, rep.Tax, "Calculated tax should match")
			assert.Equal(t, tc.expected.TaxRefund, rep.TaxRefund, "Tax refund should match")
			assert.Equal(t, tc.expected.TaxLevels, rep.TaxLevels, "TaxLevels should match")
		})
	}
}

type caseInvalids struct {
	name     string
	mockRepo *MockTaxRepository
	request  models.TaxRequest
	expected error
}

var invalids = func() (cs []caseInvalids) {
	cs = []caseInvalids{
		{
			name:     "case invalid totalIncome less than 0",
			mockRepo: &MockTaxRepository{},
			request:  models.TaxRequest{TotalIncome: -100},
			expected: errors.New(constants.ErrMessageThenZero),
		},
		{
			name:     "case invalid WHT less than 0",
			mockRepo: &MockTaxRepository{},
			request:  models.TaxRequest{TotalIncome: 500000, WHT: -100},
			expected: errors.New(constants.ErrMesssageWhtInvalid),
		},
		{
			name:     "case invalid WHT more then total income",
			mockRepo: &MockTaxRepository{},
			request:  models.TaxRequest{TotalIncome: 150000, WHT: 150001},
			expected: errors.New(constants.ErrMesssageWhtInvalid),
		},
		{
			name:     "case invalid database error repo get rates",
			mockRepo: &MockTaxRepository{ratesErr: errors.New("")},
			request: models.TaxRequest{
				TotalIncome: 150000,
				Allowances:  []models.Allowance{{AllowanceType: constants.Donation, Amount: 1000}}},
			expected: errors.New(constants.ErrMessageInternal),
		},
		{
			name:     "case invalid database error repo get allowances",
			mockRepo: &MockTaxRepository{allowancesErr: errors.New("")},
			request: models.TaxRequest{
				TotalIncome: 150000,
				Allowances:  []models.Allowance{{AllowanceType: constants.Donation, Amount: 1000}}},
			expected: errors.New(constants.ErrMessageInternal),
		},
		{
			name:     "case invalid AllowanceType not found",
			mockRepo: &MockTaxRepository{},
			request: models.TaxRequest{
				TotalIncome: 500000,
				Allowances:  []models.Allowance{{AllowanceType: "AllowanceType", Amount: 100}}},
			expected: errors.New(constants.ErrMsgAllowanceType),
		},
		{
			name:     "case invalid allowance less than 0",
			mockRepo: &MockTaxRepository{},
			request: models.TaxRequest{
				TotalIncome: 500000,
				Allowances:  []models.Allowance{{AllowanceType: constants.Donation, Amount: -1}}},
			expected: errors.New(constants.ErrMsgAllowanceThenZero),
		},
		{
			name:     "case invalid allowance personal more than 10,000",
			mockRepo: mockRepo,
			request: models.TaxRequest{
				TotalIncome: 500000,
				Allowances:  []models.Allowance{{AllowanceType: constants.Personal, Amount: 10000}}},
			expected: errors.New(constants.ErrMsgAllowanceThenMin),
		},
		{
			name:     "case invalid allowance k-receipt more than 0",
			mockRepo: mockRepo,
			request: models.TaxRequest{
				TotalIncome: 500000,
				Allowances:  []models.Allowance{{AllowanceType: constants.K_Receipt, Amount: 0}}},
			expected: errors.New(constants.ErrMsgAllowanceThenMin),
		},
	}
	return cs
}()

func TestCalculateTax_Invalids(t *testing.T) {

	for _, tc := range invalids {
		t.Run(tc.name, func(t *testing.T) {
			serv := services.NewServices(tc.mockRepo)
			rep, err := serv.TaxCalculations(&tc.request)

			assert.NotNil(t, err, "Error should not be nil for invalid")
			assert.EqualError(t, err, tc.expected.Error(), "Error message should match")
			assert.Zero(t, rep)
		})
	}
}
