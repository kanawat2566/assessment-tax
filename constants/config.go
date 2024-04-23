package constants

const (
	Personal                string  = "personal"
	Donation                string  = "donation"
	K_Receipt               string  = "k-receipt"
	AllowanceDefault        float64 = 60000.00
	MaximumWHTPercent       float64 = 5.00
	ErrMessageThenZero      string  = "Income should be greater than zero."
	ErrMesssageWhtInvalid   string  = "Withholding tax is invalid. It should be between 0 and total income."
	ErrMessageTaxInvalid    string  = "Tax invalid request"
	ErrMessageInternal      string  = "Error internal"
	ErrMsgAllowanceType             = "Allowance type not found"
	ErrMsgAllowanceThenZero         = "Allowances amount should be greater than zero."
	ErrMsgAllowanceThenMin          = "Allowance should be greater than minimun value of allowance."
	ErrMsgDatabaseError             = "database error"
)

var AllowanceTypes = map[string]string{
	"donation":  Donation,
	"k-receipt": K_Receipt,
	"personal":  Personal,
}

type ConfigTaxStepRate struct {
	MinAmt  float64
	MaxAmt  float64
	TaxRate float32
}
