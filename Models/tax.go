package models

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type Income struct {
	TotalIncome float64     `json:"totalIncome"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type TaxResponse struct {
	Tax float64 `json:"tax"`
}

type TaxLevel struct {
	Levlel int     `json:"level"`
	Tax    float64 `json:"tax"`
}

type TaxLevelResponse struct {
	Tax      float64    `json:"tax"`
	TaxLevel []TaxLevel `json:"taxLevel"`
}
