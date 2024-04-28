package services_test

import (
	"errors"
	"testing"

	cm "github.com/kanawat2566/assessment-tax/common"
	ct "github.com/kanawat2566/assessment-tax/constants"
	md "github.com/kanawat2566/assessment-tax/model"
	"github.com/kanawat2566/assessment-tax/repository"
	"github.com/kanawat2566/assessment-tax/services"
	"github.com/stretchr/testify/assert"
)

type MockTaxRepository struct {
	taxRates   []*repository.IncomeTaxRates
	allowances map[string]repository.Allowances
	taxErr     error
	awcErr     error
	updateErr  error
}

func (m *MockTaxRepository) GetTaxRates() (res []*repository.IncomeTaxRates, err error) {
	return m.taxRates, m.taxErr
}

func (m *MockTaxRepository) GetLimitAllowances(allowanceType string) (r repository.Allowances, err error) {
	return m.allowances[allowanceType], m.awcErr
}
func (m *MockTaxRepository) UpdateConfigDeduct(config ct.Deduction) error {
	return m.updateErr
}

type TaxCase struct {
	name     string
	request  md.TaxRequest
	expected md.TaxResponse
}

var _mockRepo = &MockTaxRepository{
	taxRates:   _taxRates,
	allowances: _allowances,
}
var _taxRates = []*repository.IncomeTaxRates{
	{IncomeLevel: "0-150,000", MinIncome: 0, MaxIncome: 150000, TaxRate: 0},
	{IncomeLevel: "150,001-500,000", MinIncome: 150001, MaxIncome: 500000, TaxRate: 10},
	{IncomeLevel: "500,001-1,000,000", MinIncome: 500001.00, MaxIncome: 1000000.00, TaxRate: 15},
	{IncomeLevel: "1,000,001-2,000,000", MinIncome: 1000001.00, MaxIncome: 2000000.00, TaxRate: 20},
	{IncomeLevel: "2,000,001 ขึ้นไป", MinIncome: 2000001.00, MaxIncome: 99999999999999.00, TaxRate: 35},
}
var _allowances = map[string]repository.Allowances{
	ct.Personal:  {Allowance_name: ct.Personal, LimitAmt: 60000, MinAmt: 10001, MaxAmt: 100000},
	ct.Donation:  {Allowance_name: ct.Donation, LimitAmt: 100000, MinAmt: 0, MaxAmt: 100000},
	ct.K_Receipt: {Allowance_name: ct.K_Receipt, LimitAmt: 50000, MinAmt: 1, MaxAmt: 100000},
}

func TestCalculateTax_Valids(t *testing.T) {

	cases := []TaxCase{
		{
			name: "given input income total should payment tax",
			request: md.TaxRequest{
				TotalIncome: 500000,
				WHT:         0,
				Allowances: []md.Allowance{
					{
						AllowanceType: ct.Donation,
						Amount:        0,
					}}},
			expected: md.TaxResponse{
				Tax: 29000,
			},
		},
		{
			name: "given input income total and deducting WHT should payment tax",
			request: md.TaxRequest{
				TotalIncome: 500000,
				WHT:         25000.0,
			},
			expected: md.TaxResponse{
				Tax: 4000,
			},
		},
		{
			name: "given input income total and deducting WHT should return tax refund",
			request: md.TaxRequest{
				TotalIncome: 150000,
				WHT:         2000,
			},
			expected: md.TaxResponse{
				Tax:       0,
				TaxRefund: 2000,
			},
		},
		{
			name: "given input income total and allownce donation should payment tax",
			request: md.TaxRequest{
				TotalIncome: 500000,
				WHT:         0,
				Allowances: []md.Allowance{
					{
						AllowanceType: ct.Donation,
						Amount:        200000.0,
					}}},
			expected: md.TaxResponse{
				Tax: 19000,
			},
		},
		{
			name: "given input income total and allownce fix personal should payment tax",
			request: md.TaxRequest{
				TotalIncome: 500000,
				WHT:         0,
				Allowances: []md.Allowance{
					{
						AllowanceType: ct.Personal,
						Amount:        50000,
					},
				},
			},
			expected: md.TaxResponse{
				Tax: 30000,
			},
		},
		{
			name: "given input income total and deducting WHT with allownce should return tax refund",
			request: md.TaxRequest{
				TotalIncome: 500000,
				WHT:         25000,
				Allowances: []md.Allowance{
					{
						AllowanceType: ct.Donation,
						Amount:        100000,
					},
					{
						AllowanceType: ct.K_Receipt,
						Amount:        100000,
					},
				},
			},
			expected: md.TaxResponse{
				Tax:       0,
				TaxRefund: 11000,
			},
		},
		{
			name: "given input income total and deducting WHT with allownce should return tax level detail",
			request: md.TaxRequest{
				TotalIncome: 500000.0,
				WHT:         0.0,
				Allowances: []md.Allowance{
					{
						AllowanceType: ct.Donation,
						Amount:        200000.0,
					},
				},
			},
			expected: md.TaxResponse{
				Tax: 19000.0,
				TaxLevels: []md.TaxLevel{
					{Level: "0-150,000", Tax: 0.0},
					{Level: "150,001-500,000", Tax: 19000.0},
					{Level: "500,001-1,000,000", Tax: 0.0},
					{Level: "1,000,001-2,000,000", Tax: 0.0},
					{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
				},
			},
		},
		{
			name: "given input income total and allownces should payment tax",
			request: md.TaxRequest{
				TotalIncome: 500000,
				WHT:         0,
				Allowances: []md.Allowance{
					{
						AllowanceType: ct.K_Receipt,
						Amount:        200000,
					},
					{
						AllowanceType: ct.Donation,
						Amount:        100000,
					},
				},
			},
			expected: md.TaxResponse{
				Tax: 14000,
			},
		},
	}

	for _, tc := range cases {

		t.Run(tc.name, func(t *testing.T) {
			serv := services.NewServices(_mockRepo)
			rep, err := serv.TaxCalculations(tc.request)

			// Assertions
			assert.Nil(t, err, "Error should be nil for valid inputs")
			assert.Equal(t, tc.expected.Tax, rep.Tax, "Calculated tax should match")
			assert.Equal(t, tc.expected.TaxRefund, rep.TaxRefund, "Tax refund should match")
		})
	}
}

func TestCalculateTaxLevel_Valids(t *testing.T) {

	cases := []TaxCase{
		{
			name: "given input income total and allownce should return tax level detail",
			request: md.TaxRequest{
				TotalIncome: 500000.0,
				WHT:         0.0,
				Allowances: []md.Allowance{
					{
						AllowanceType: ct.Donation,
						Amount:        200000.0,
					},
				},
			},
			expected: md.TaxResponse{
				Tax: 19000.0,
				TaxLevels: []md.TaxLevel{
					{Level: "0-150,000", Tax: 0.0},
					{Level: "150,001-500,000", Tax: 19000.0},
					{Level: "500,001-1,000,000", Tax: 0.0},
					{Level: "1,000,001-2,000,000", Tax: 0.0},
					{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
				},
			},
		},
		{
			name: "given input income total and allownces should return tax level detail",
			request: md.TaxRequest{
				TotalIncome: 500000.0,
				WHT:         0.0,
				Allowances: []md.Allowance{
					{
						AllowanceType: ct.Donation,
						Amount:        100000,
					},
					{
						AllowanceType: ct.K_Receipt,
						Amount:        200000.0,
					},
				},
			},
			expected: md.TaxResponse{
				Tax: 14000,
				TaxLevels: []md.TaxLevel{
					{Level: "0-150,000", Tax: 0.0},
					{Level: "150,001-500,000", Tax: 14000.0},
					{Level: "500,001-1,000,000", Tax: 0.0},
					{Level: "1,000,001-2,000,000", Tax: 0.0},
					{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
				},
			},
		},
	}

	for _, tc := range cases {

		t.Run(tc.name, func(t *testing.T) {
			serv := services.NewServices(_mockRepo)
			rep, err := serv.TaxCalculations(tc.request)

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
	request  md.TaxRequest
	expected error
}

func TestCalculateTax_Invalids(t *testing.T) {
	invalids := []caseInvalids{
		{
			name:     "case invalid totalIncome less than 0",
			mockRepo: &MockTaxRepository{},
			request:  md.TaxRequest{TotalIncome: -100},
			expected: errors.New(ct.ErrMessageThenZero),
		},
		{
			name:     "case invalid WHT less than 0",
			mockRepo: &MockTaxRepository{},
			request:  md.TaxRequest{TotalIncome: 500000, WHT: -100},
			expected: errors.New(ct.ErrMesssageWhtInvalid),
		},
		{
			name:     "case invalid WHT more then total income",
			mockRepo: &MockTaxRepository{},
			request:  md.TaxRequest{TotalIncome: 150000, WHT: 150001},
			expected: errors.New(ct.ErrMesssageWhtInvalid),
		},
		{
			name:     "case invalid database error repo get rates",
			mockRepo: &MockTaxRepository{taxErr: errors.New("")},
			request: md.TaxRequest{
				TotalIncome: 150000,
				Allowances:  []md.Allowance{{AllowanceType: ct.Donation, Amount: 1000}}},
			expected: errors.New(ct.ErrMessageInternal),
		},
		{
			name:     "case invalid database error repo get allowances",
			mockRepo: &MockTaxRepository{awcErr: errors.New("")},
			request: md.TaxRequest{
				TotalIncome: 150000,
				Allowances:  []md.Allowance{{AllowanceType: ct.Donation, Amount: 1000}}},
			expected: errors.New(ct.ErrMessageInternal),
		},
		{
			name:     "case invalid AllowanceType not found",
			mockRepo: &MockTaxRepository{},
			request: md.TaxRequest{
				TotalIncome: 500000,
				Allowances:  []md.Allowance{{AllowanceType: "AllowanceType", Amount: 100}}},
			expected: errors.New(ct.ErrMsgAllowanceType),
		},
		{
			name:     "case invalid allowance less than 0",
			mockRepo: &MockTaxRepository{},
			request: md.TaxRequest{
				TotalIncome: 500000,
				Allowances:  []md.Allowance{{AllowanceType: ct.Donation, Amount: -1}}},
			expected: errors.New(ct.ErrMsgAllowanceThenZero),
		},
		{
			name:     "case invalid allowance personal greater than minimum config",
			mockRepo: _mockRepo,
			request: md.TaxRequest{
				TotalIncome: 500000,
				Allowances:  []md.Allowance{{AllowanceType: ct.Personal, Amount: 10000}}},
			expected: errors.New(ct.ErrMsgAllowanceThenMin),
		},
		{
			name:     "case invalid allowance k-receipt more than 0",
			mockRepo: _mockRepo,
			request: md.TaxRequest{
				TotalIncome: 500000,
				Allowances:  []md.Allowance{{AllowanceType: ct.K_Receipt, Amount: 0}}},
			expected: errors.New(ct.ErrMsgAllowanceThenMin),
		},
	}
	for _, tc := range invalids {
		t.Run(tc.name, func(t *testing.T) {
			serv := services.NewServices(tc.mockRepo)
			rep, err := serv.TaxCalculations(tc.request)

			assert.NotNil(t, err, "Error should not be nil for invalid")
			assert.EqualError(t, err, tc.expected.Error(), "Error message should match")
			assert.Zero(t, rep)
		})
	}
}

type ConfigCase struct {
	name     string
	request  ct.Deduction
	expected ct.Deduction
}

func TestConfigDeduction_Valids(t *testing.T) {
	cases := []ConfigCase{
		{name: "given admin set personal deduction amount 70,000 should return 70,000",
			request:  ct.Deduction{Type: ct.Personal, Amount: 70000},
			expected: ct.Deduction{Type: ct.Personal, Amount: 70000},
		},
		{name: "given admin set personal deduction amount 10,001 should return 10,001",
			request:  ct.Deduction{Type: ct.Personal, Amount: 10001},
			expected: ct.Deduction{Type: ct.Personal, Amount: 10001},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			serv := services.NewServices(_mockRepo)
			rep, err := serv.SetAdminDeductions(tc.request)
			assert.Nil(t, err, "Error should be nil for valid inputs")
			assert.Equal(t, tc.expected.Amount, rep.Amount, "Calculated tax should match")
		})

	}

}

type CaseConfigInvalids struct {
	name     string
	mockRepo *MockTaxRepository
	request  ct.Deduction
	expected error
}

func TestConfigDeduction_Invalids(t *testing.T) {
	invalids := []CaseConfigInvalids{
		{
			name:     "case invalid type should return ErrInvalidDeductionType",
			mockRepo: &MockTaxRepository{},
			request:  ct.Deduction{},
			expected: errors.New(ct.ErrMsgInvalidDeduct),
		},
		{
			name:     "case invalid type should return ErrInvalid Not Supported",
			mockRepo: &MockTaxRepository{},
			request:  ct.Deduction{Type: ct.Donation},
			expected: errors.New(ct.ErrMsgNotDeductSupport),
		},
		{
			name:     "case invalid get deductions from database error should return error",
			mockRepo: &MockTaxRepository{awcErr: errors.New("error")},
			request:  ct.Deduction{Type: ct.Personal, Amount: 50000},
			expected: errors.New(ct.ErrMessageInternal),
		},
		{
			name: "case invalid get deduction personal from database not found",
			mockRepo: &MockTaxRepository{allowances: map[string]repository.Allowances{
				ct.Donation: {Allowance_name: ct.Donation, LimitAmt: 100000, MinAmt: 0, MaxAmt: 100000},
			}},
			request:  ct.Deduction{Type: ct.Personal, Amount: 50000},
			expected: errors.New(ct.ErrMsgDeductNotFound),
		},
		{
			name:     "case invalid deductions amount must be greater minimum should return error",
			mockRepo: _mockRepo,
			request:  ct.Deduction{Type: ct.Personal, Amount: 10000},
			expected: errors.New(cm.MsgWithNumber(ct.ErrMsgValidateMinAmt, 10001)),
		},
		{
			name:     "case invalid deductions less than maximum should return error",
			mockRepo: _mockRepo,
			request:  ct.Deduction{Type: ct.Personal, Amount: 100001},
			expected: errors.New(cm.MsgWithNumber(ct.ErrMsgValidateMaxAmt, 100000)),
		},
		{
			name:     "case invalid set deductions from database error should return error",
			mockRepo: &MockTaxRepository{allowances: _allowances, updateErr: errors.New("error")},
			request:  ct.Deduction{Type: ct.Personal, Amount: 50000},
			expected: errors.New(ct.ErrMessageInternal),
		},
		{
			name: "case invalid get deduction k-receipt from database not found",
			mockRepo: &MockTaxRepository{allowances: map[string]repository.Allowances{
				ct.Donation: {Allowance_name: ct.Donation, LimitAmt: 100000, MinAmt: 0, MaxAmt: 100000},
			}},
			request:  ct.Deduction{Type: ct.K_Receipt, Amount: 50000},
			expected: errors.New(ct.ErrMsgDeductNotFound),
		},
		{
			name:     "case invalid deduction k-receipt amount must be greater minimum should return error",
			mockRepo: _mockRepo,
			request:  ct.Deduction{Type: ct.K_Receipt, Amount: 0},
			expected: errors.New(cm.MsgWithNumber(ct.ErrMsgValidateMinAmt, 1)),
		},
		{
			name:     "case invalid deduction k-receipt less than maximum should return error",
			mockRepo: _mockRepo,
			request:  ct.Deduction{Type: ct.K_Receipt, Amount: 100001},
			expected: errors.New(cm.MsgWithNumber(ct.ErrMsgValidateMaxAmt, 100000)),
		},
	}
	for _, tc := range invalids {
		t.Run(tc.name, func(t *testing.T) {
			serv := services.NewServices(tc.mockRepo)
			rep, err := serv.SetAdminDeductions(tc.request)

			assert.NotNil(t, err, "Error should not be nil for invalid")
			assert.EqualError(t, err, tc.expected.Error(), "Error message should match")
			assert.Zero(t, rep)
		})
	}
}
