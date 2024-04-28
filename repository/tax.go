package repository

import (
	"database/sql"
	"errors"

	ct "github.com/kanawat2566/assessment-tax/constants"
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
	UpdateConfigDeduct(config ct.Deduction) error
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
		return nil, errors.New(ct.ErrMsgDatabaseError)
	}
	defer rows.Close()
	var incomeTaxRates []*IncomeTaxRates
	for rows.Next() {
		var t IncomeTaxRates
		err = rows.Scan(&t.ID, &t.IncomeLevel, &t.MinIncome, &t.MaxIncome, &t.TaxRate)
		if err != nil {
			return nil, errors.New(ct.ErrMsgDatabaseError)
		}
		incomeTaxRates = append(incomeTaxRates, &t)
	}
	return incomeTaxRates, nil
}
func (p *Postgres) GetLimitAllowances(allowanceType string) (Allowances, error) {
	res := Allowances{}

	if allowanceType == "" {
		return res, errors.New(ct.ErrMsgAllowanceType)
	}

	query := `
	SELECT  
	max_allowance,
	min_allowance,
	limit_allowance 
	FROM allowances 
	WHERE allowance_name=$1`

	row := p.Db.QueryRow(query, allowanceType)

	err := row.Scan(&res.MaxAmt, &res.MinAmt, &res.LimitAmt)
	if err == sql.ErrNoRows {
		return res, errors.New(ct.ErrMsgDatabaseError)
	}
	res.Allowance_name = allowanceType
	return res, nil
}

func (p *Postgres) UpdateConfigDeduct(config ct.Deduction) error {
	query := `UPDATE allowances SET limit_allowance = $1 WHERE allowance_name=$2;`
	res, err := p.Db.Exec(query, config.Amount, config.Type)
	if err != nil {
		return errors.New(ct.ErrMsgDatabaseError)
	}
	affect, _ := res.RowsAffected()
	if affect < 1 {
		return errors.New(ct.ErrMsgUpdateNotSuccess)
	}
	return nil
}
