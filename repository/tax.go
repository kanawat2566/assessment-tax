package repository

type IncomeTaxRates struct {
	ID          int     `postgres:"id"`
	IncomeLevel string  `postgres:"income_level"`
	MinIncome   float32 `postgres:"min_income"`
	MaxIncome   float32 `postgres:"max_income"`
	TaxRate     float32 `postgres:"tax_rate"`
}

type TaxRepository interface {
	GetAll() ([]*IncomeTaxRates, error)
}

func (p *Postgres) GetAll() ([]*IncomeTaxRates, error) {
	rows, err := p.Db.Query("SELECT * FROM income_tax_rates")
	if err != nil {
		return nil, err
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
