package repository

import "errors"

type IncomeTaxRates struct {
	ID          int     `postgres:"id"`
	IncomeLevel string  `postgres:"income_level"`
	MinIncome   float64 `postgres:"min_income"`
	MaxIncome   float64 `postgres:"max_income"`
	TaxRate     float64 `postgres:"tax_rate"`
}

type Allowances struct {
	Allowance_name string  `postgres:"allowance_name"`
	MinAmt         float64 `postgres:"min_allowance"`
	MaxAmt         float64 `postgres:"max_allowance"`
	LimitAmt       float64 `postgres:"limit_allowance"`
}

type TaxRepository interface {
	GetTaxRates() ([]*IncomeTaxRates, error)
}

func (p *Postgres) GetTaxRates() ([]*IncomeTaxRates, error) {
	rows, err := p.Db.Query(`
	SELECT 
	id, income_level, 
	min_income, max_income, 
	tax_rate 
	FROM income_tax_rates
	ORDER BY id;`)

	if err != nil {
		return nil, errors.New("database error")
	}
	defer rows.Close()
	var incomeTaxRates []*IncomeTaxRates
	for rows.Next() {
		var t IncomeTaxRates
		err = rows.Scan(&t.ID, &t.IncomeLevel, &t.MinIncome, &t.MaxIncome, &t.TaxRate)
		if err != nil {
			return nil, err
		}
		incomeTaxRates = append(incomeTaxRates, &t)
	}
	return incomeTaxRates, nil
}
