package services

import (
	"errors"
	"fmt"
	"math"
	"strings"

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
	SetAdminDeductions(req ct.DeductConfig) (ct.DeductConfig, error)
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
			fmt.Println(err)
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

func (ts *taxService) SetAdminDeductions(req ct.DeductConfig) (ct.DeductConfig, error) {

	if req.Type == "" {
		return ct.DeductConfig{}, errors.New(ct.ErrMsgInvalidDeduct)
	}

	if len(ct.Deductios[req.Type].Type) == 0 {
		return ct.DeductConfig{}, errors.New(ct.ErrMsgNotDeductSupport)
	}

	return ct.DeductConfig{Amount: req.Amount}, nil
}
