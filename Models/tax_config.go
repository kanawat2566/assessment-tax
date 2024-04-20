package models

type ConfigTax struct {
	Name       string
	MinAmt     float64
	MaxAmt     float64
	DefaultAmt float64
}

type ConfigTaxStepRate struct {
	MinAmt  float64
	MaxAmt  float64
	TaxRate float32
}

const (
	Max_Donations_Allowance       = 100000
	Default_Personal_Allowance    = 60000
	K_Receipt_Max_Allowance       = 50000
	Admin_Max_Personal_Allowance  = 100000
	Admin_Max_K_Receipt_Allowance = 100000
	Personal_Min_Aloowance        = 10000
)
