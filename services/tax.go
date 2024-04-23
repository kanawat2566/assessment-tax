package services

import (
	"errors"
	"math"

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

	incomeTotal := taxRequest.TotalIncome - constants.AllowanceDefault

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

	//เช็คยอด WHT ต้องน้อยกว่าหรือเท่ากับ 5% ของรายได้
	if taxRequest.WHT > taxRequest.TotalIncome*(constants.MaximumWHTPercent/100) {
		return errors.New(constants.ErrMesssageWhtInvalid)
	}
	return nil
}
