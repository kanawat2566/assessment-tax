package models

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type TaxRequest struct {
	TotalIncome float64     `json:"totalIncome" validate:"required,numeric,gte=0"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type TaxResponse struct {
	Tax       float64    `json:"tax"`
	TaxRefund float64    `json:"taxRefund"`
	TaxLevels []TaxLevel `json:"taxLevel"`
}

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type TaxLevelReponse struct {
	Tax       float64    `json:"tax"`
	TaxLevels []TaxLevel `json:"taxLevel"`
}

type DeductRequest struct {
	Amount float64 `json:"amount" validate:"required,numeric,gt=0"`
}

type DeductResponse struct {
	PersonalDeduction float64 `json:"personalDeduction"`
	KReceipt          float64 `json:"kReceipt"`
}
