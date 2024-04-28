package constants

const (
	UserAuth string = "adminTax"
	PassAuth string = "admin!"

	Personal  string = "personal"
	Donation  string = "donation"
	K_Receipt string = "k-receipt"

	AllowanceDefault  float64 = 60000.00
	MaximumWHTPercent float64 = 5.00

	ErrInvalidFormatReq     string = "Error: Invalid format request."
	ErrMessageThenZero      string = "Income should be greater than zero."
	ErrMesssageWhtInvalid   string = "Withholding tax is invalid. It should be between 0 and total income."
	ErrMessageTaxInvalid    string = "Tax invalid request"
	ErrMessageInternal      string = "Error internal"
	ErrMsgAllowanceType     string = "Allowance type not found"
	ErrMsgAllowanceThenZero string = "Allowances amount should be greater than zero."
	ErrMsgAllowanceThenMin  string = "Allowance should be greater than minimun value of allowance."
	ErrMsgDatabaseError     string = "Database error"
	ErrMsgInvalidDeduct     string = "Invalid deduction type"
	ErrMsgDeductNotFound    string = "Deduction type not found"
	ErrMsgNotDeductSupport  string = "Not Supported Deduction type"
	ErrMsgUpdateNotSuccess  string = "Failed to update data in database"
	ErrMsgValidateMinAmt    string = "Deduction amount must be greater or equal to"
	ErrMsgValidateMaxAmt    string = "Deduction amount should be less than or equal to"
	ErrMsgInvalidPathParam  string = "Invalid path param"

	ErrMsgCsvInvaildFormat string = "format is wrong, please check your format."
	ErrMsgFileNoUpload     string = "No file uploaded"
	ErrMsgReadCsvFailed    string = "Failed to read csv file"
	ErrInvalidIncomeCsv    string = "Invalid income number in line"
	ErrInvalidWHTCsv       string = "Invalid WHT number in line"
	ErrInvalidDonationCsv  string = "Invalid donation number in line"

	PathParamUploadCsv string = "upload-csv"
)

var AllowanceTypes = map[string]string{
	"donation":  Donation,
	"k-receipt": K_Receipt,
	"personal":  Personal,
}

type Deduction struct {
	Type   string
	Name   string
	Amount float64
	MinAmt float64
	MaxAmt float64
}

var Deductions = map[string]Deduction{
	Personal:  {Type: Personal, Name: "personalDeduction"},
	K_Receipt: {Type: K_Receipt, Name: "kReceipt"},
}

var CsvFomatFile = []string{"totalIncome", "wht", "donation"}
