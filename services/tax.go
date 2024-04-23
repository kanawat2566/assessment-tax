package services

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/kanawat2566/assessment-tax/constants"
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
}

func (ts *taxService) TaxCalculations(taxRequest *models.TaxRequest) (models.TaxResponse, error) {
	var taxResp models.TaxResponse
	var tax float64

	if err := validateInputs(taxRequest); err != nil {
		return taxResp, err
	}

	rates, err := ts.repo.GetTaxRates()
	if err != nil {
		return taxResp, errors.New(constants.ErrMessageInternal)
	}

	allowances, err := ts.allowanceCal(taxRequest.Allowances)
	if err != nil {
		return taxResp, err
	}

	incomeTotal := taxRequest.TotalIncome - allowances

	for _, v := range rates {

		if incomeTotal >= v.MinIncome && v.TaxRate > 0 {

			baseCal := math.Min(v.MaxIncome, incomeTotal) - (v.MinIncome - 1)
			tax += baseCal * (v.TaxRate / 100)
		}

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
		return errors.New(constants.ErrMessageThenZero)
	}
	if taxRequest.WHT < 0 {
		return errors.New(constants.ErrMesssageWhtInvalid)
	}

	//เช็คยอด WHT ต้องน้อยกว่าหรือเท่ากับรายได้
	if taxRequest.WHT > taxRequest.TotalIncome {
		return errors.New(constants.ErrMesssageWhtInvalid)
	}
	return nil
}
func (ts *taxService) allowanceCal(allowances []models.Allowance) (float64, error) {
	total := 0.00
	var chkPersonal bool

	for _, v := range allowances {
		at, ok := constants.AllowanceTypes[strings.ToLower(v.AllowanceType)]
		if !ok {
			return total, errors.New(constants.ErrMsgAllowanceType)
		}
		if v.Amount < 0 {
			return total, errors.New(constants.ErrMsgAllowanceThenZero)
		}

		amt, err := ts.repo.GetLimitAllowances(at)
		if err != nil {
			fmt.Println(err)
			return total, errors.New(constants.ErrMessageInternal)
		}

		if v.Amount < amt.MinAmt {
			return total, errors.New(constants.ErrMsgAllowanceThenMin)
		}

		total += math.Min(v.Amount, amt.LimitAmt)

		if at == constants.Personal {
			chkPersonal = true

		}
	}

	// default personal allowance
	if !chkPersonal {
		p, _ := ts.repo.GetLimitAllowances(constants.Personal)
		total += p.LimitAmt
	}

	return total, nil
}
