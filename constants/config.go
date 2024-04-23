package constants

const (
	Personal              string  = "Personal"
	Donation              string  = "donation"
	K_Receipt             string  = "k-receipt"
	AllowanceDefault      float64 = 60000.00
	MaximumWHTPercent     float64 = 5.00
	ErrMessageThenZero    string  = "Income should be greater than zero."
	ErrMesssageWhtInvalid string  = "Withholding tax is invalid. It should be between 0 and 5 percent of total income."
	ErrMessageTaxInvalid  string  = "Tax invalid request"
	ErrMessageInternal    string  = "Error internal"
)

var AllowanceTypes = []string{Donation, K_Receipt, Personal}

type ConfigTaxStepRate struct {
	MinAmt  float64
	MaxAmt  float64
	TaxRate float32
}
