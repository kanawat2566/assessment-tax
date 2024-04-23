package repository

import (
	"database/sql"
	"errors"

	"github.com/kanawat2566/assessment-tax/constants"
)

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
	GetLimitAllowances(allowanceType string) (Allowances, error)
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
		return nil, errors.New(constants.ErrMsgDatabaseError)
	}
	defer rows.Close()
	var incomeTaxRates []*IncomeTaxRates
	for rows.Next() {
		var t IncomeTaxRates
		err = rows.Scan(&t.ID, &t.IncomeLevel, &t.MinIncome, &t.MaxIncome, &t.TaxRate)
		if err != nil {
			return nil, errors.New(constants.ErrMsgDatabaseError)
		}
		incomeTaxRates = append(incomeTaxRates, &t)
	}
	return incomeTaxRates, nil
}
func (p *Postgres) GetLimitAllowances(allowanceType string) (Allowances, error) {
	res := Allowances{}

	if allowanceType == "" {
		return res, errors.New(constants.ErrMsgAllowanceType)
	}

	query := "SELECT  max_allowance,min_allowance,limit_allowance FROM allowances WHERE allowance_name=$1"
	row := p.Db.QueryRow(query, allowanceType)

	err := row.Scan(&res.MaxAmt, &res.MaxAmt, &res.LimitAmt)
	if err == sql.ErrNoRows {
		return res, errors.New(constants.ErrMsgDatabaseError)
	}
	res.Allowance_name = allowanceType
	return res, nil
}
