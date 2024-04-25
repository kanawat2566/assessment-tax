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
	dtype, ok := ct.Deductios[strings.ToLower(req.Type)]
	if !ok {
		return ct.DeductConfig{}, errors.New(ct.ErrMsgNotDeductSupport)
	}

	config, err := ts.repo.GetLimitAllowances(dtype.Type)
	if err != nil {
		return ct.DeductConfig{}, errors.New(ct.ErrMessageInternal)
	}

	if len(config.Allowance_name) == 0 {
		return ct.DeductConfig{}, errors.New(ct.ErrMsgDeductNotFound)
	}

	if req.Amount < config.MinAmt {
		return ct.DeductConfig{}, errors.New(cm.MsgWithNumber(ct.ErrMsgValidateMinAmt, config.MinAmt))
	}

	if req.Amount > config.MaxAmt {
		return ct.DeductConfig{}, errors.New(cm.MsgWithNumber(ct.ErrMsgValidateMaxAmt, config.MaxAmt))
	}

	uErr := ts.repo.UpdateConfigDeduct(req)
	if uErr != nil {
		return ct.DeductConfig{}, errors.New(ct.ErrMessageInternal)
	}

	return ct.DeductConfig{Type: dtype.Type, Name: dtype.Name, Amount: req.Amount}, nil
}
