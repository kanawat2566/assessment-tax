package constants

const (
	UserAuth string = "adminTax"
	PassAuth string = "admin!"

	Personal  string = "personal"
	Donation  string = "donation"
	K_Receipt string = "k-receipt"

	AllowanceDefault  float64 = 60000.00
	MaximumWHTPercent float64 = 5.00

	ErrMessageThenZero      string = "Income should be greater than zero."
	ErrMesssageWhtInvalid   string = "Withholding tax is invalid. It should be between 0 and total income."
	ErrMessageTaxInvalid    string = "Tax invalid request"
	ErrMessageInternal      string = "Error internal"
	ErrMsgAllowanceType     string = "Allowance type not found"
	ErrMsgAllowanceThenZero string = "Allowances amount should be greater than zero."
	ErrMsgAllowanceThenMin  string = "Allowance should be greater than minimun value of allowance."
	ErrMsgDatabaseError     string = "database error"
	ErrMsgInvalidDeduct     string = "Invalid Deduction Type"
	ErrMsgNotDeductSupport  string = "Not Supported Deduction Type"
)

var AllowanceTypes = map[string]string{
	"donation":  Donation,
	"k-receipt": K_Receipt,
	"personal":  Personal,
}

type DeductConfig struct {
	Type   string
	Name   string
	Amount float64
}

var Deductios = map[string]DeductConfig{
	Personal:  {Type: Personal, Name: "personalDeduction"},
	K_Receipt: {Type: K_Receipt, Name: "kReceipt"},
}
