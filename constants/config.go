package constants

const (
	Personal              = "Personal"
	Donation              = "donation"
	K_Receipt             = "k-receipt"
	AllowanceDefault      = 60000.00
	ErrMessageThenZero    = "Income should be greater than zero."
	ErrMesssageWhtInvalid = "Withholding tax is invalid. It should be between 0 and total income."
	ErrMessageTaxInvalid  = "Tax invalid request"
	ErrMessageInternal = "Error internal"

)

var AllowanceTypes = []string{Donation, K_Receipt, Personal}

type ConfigTaxStepRate struct {
	MinAmt  float64
	MaxAmt  float64
	TaxRate float32
}
