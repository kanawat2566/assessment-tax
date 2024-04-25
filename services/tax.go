package services

import (
	"errors"
	"math"
	"strings"

	cm "github.com/kanawat2566/assessment-tax/common"
	ct "github.com/kanawat2566/assessment-tax/constants"
	models "github.com/kanawat2566/assessment-tax/model"
	"github.com/kanawat2566/assessment-tax/repository"
)

type taxService struct {
	repo repository.TaxRepository
}

func NewServices(r repository.TaxRepository) *taxService {
	return &taxService{repo: r}
}

type TaxService interface {
	TaxCalculations(taxRequest *models.TaxRequest) (models.TaxResponse, error)
	SetAdminDeductions(req ct.Deduction) (ct.Deduction, error)
}

func (ts *taxService) TaxCalculations(taxRequest *models.TaxRequest) (models.TaxResponse, error) {
	var taxResp models.TaxResponse
	var tax float64

	if err := validateInputs(taxRequest); err != nil {
		return taxResp, err
	}

	rates, err := ts.repo.GetTaxRates()
	if err != nil {
		return taxResp, errors.New(ct.ErrMessageInternal)
	}

	allowances, err := ts.allowanceCal(taxRequest.Allowances)
	if err != nil {
		return taxResp, err
	}

	incomeTotal := taxRequest.TotalIncome - allowances

	for _, v := range rates {

		var tl models.TaxLevel
		tl.Level = v.IncomeLevel

		if incomeTotal >= v.MinIncome && v.TaxRate > 0 {

			baseCal := math.Min(v.MaxIncome, incomeTotal) - (v.MinIncome - 1)
			tl.Tax = baseCal * (v.TaxRate / 100)
			tax += tl.Tax
		}
		taxResp.TaxLevels = append(taxResp.TaxLevels, tl)
	}

	tax -= taxRequest.WHT
	if tax < 0 {
		taxResp.TaxRefund = math.Abs(tax)
	} else {
		taxResp.Tax = tax
	}

	return taxResp, nil
}

func validateInputs(taxRequest *models.TaxRequest) error {
	if taxRequest.TotalIncome <= 0 {
		return errors.New(ct.ErrMessageThenZero)
	}
	if taxRequest.WHT < 0 {
		return errors.New(ct.ErrMesssageWhtInvalid)
	}

	//เช็คยอด WHT ต้องน้อยกว่าหรือเท่ากับรายได้
	if taxRequest.WHT > taxRequest.TotalIncome {
		return errors.New(ct.ErrMesssageWhtInvalid)
	}
	return nil
}
func (ts *taxService) allowanceCal(allowances []models.Allowance) (float64, error) {
	total := 0.00
	var chkPersonal bool

	for _, v := range allowances {
		at, ok := ct.AllowanceTypes[strings.ToLower(v.AllowanceType)]
		if !ok {
			return total, errors.New(ct.ErrMsgAllowanceType)
		}
		if v.Amount < 0 {
			return total, errors.New(ct.ErrMsgAllowanceThenZero)
		}

		amt, err := ts.repo.GetLimitAllowances(at)
		if err != nil {
			return total, errors.New(ct.ErrMessageInternal)
		}

		if v.Amount < amt.MinAmt {
			return total, errors.New(ct.ErrMsgAllowanceThenMin)
		}

		total += math.Min(v.Amount, amt.LimitAmt)

		if at == ct.Personal {
			chkPersonal = true

		}
	}

	// default personal allowance
	if !chkPersonal {
		p, _ := ts.repo.GetLimitAllowances(ct.Personal)
		total += p.LimitAmt
	}

	return total, nil
}

func (ts *taxService) SetAdminDeductions(req ct.Deduction) (ct.Deduction, error) {

	if err := validateDeductionType(req.Type); err != nil {
		return ct.Deduction{}, err
	}

	d, err := ts.getDeductionDetails(req.Type)
	if err != nil {
		return ct.Deduction{}, err
	}

	if err := validateDeductionAmount(req.Amount, d); err != nil {
		return ct.Deduction{}, err
	}

	if err := ts.repo.UpdateConfigDeduct(req); err != nil {
		return ct.Deduction{}, errors.New(ct.ErrMessageInternal)
	}

	return ct.Deduction{Type: d.Type, Name: d.Name, Amount: req.Amount}, nil
}

func validateDeductionType(dtype string) error {
	if dtype == "" {
		return errors.New(ct.ErrMsgInvalidDeduct)
	}
	return nil
}

func (ts *taxService) getDeductionDetails(dtype string) (ct.Deduction, error) {
	dtypeLower := strings.ToLower(dtype)
	d, ok := ct.Deductions[dtypeLower]
	if !ok {
		return ct.Deduction{}, errors.New(ct.ErrMsgNotDeductSupport)
	}
	res, err := ts.repo.GetLimitAllowances(d.Type)
	if err != nil {
		return ct.Deduction{}, errors.New(ct.ErrMessageInternal)
	}
	if len(res.Allowance_name) == 0 {
		return ct.Deduction{}, errors.New(ct.ErrMsgDeductNotFound)
	}
	return ct.Deduction{Type: d.Type, Name: d.Name, MinAmt: res.MinAmt, MaxAmt: res.MaxAmt}, nil
}

func validateDeductionAmount(amount float64, d ct.Deduction) error {
	if amount < d.MinAmt {
		return errors.New(cm.MsgWithNumber(ct.ErrMsgValidateMinAmt, d.MinAmt))
	}
	if amount > d.MaxAmt {
		return errors.New(cm.MsgWithNumber(ct.ErrMsgValidateMaxAmt, d.MaxAmt))
	}
	return nil
}
