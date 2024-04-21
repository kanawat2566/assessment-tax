package services

import (
	"errors"
	"math"

	"github.com/kanawat2566/assessment-tax/constants"
	"github.com/kanawat2566/assessment-tax/model"
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

	if taxRequest.TotalIncome <= 0 {
		return taxResp, errors.New(constants.ErrMessageThenZero)
	}

	rates, err := ts.repo.GetTaxRates()
	if err != nil {
		return taxResp, err
	}

	incomeTotal := taxRequest.TotalIncome - constants.AllowanceDefault

	for _, v := range rates {

		if incomeTotal >= v.MinIncome && v.TaxRate > 0 {

			baseCal := math.Min(v.MaxIncome, incomeTotal) - (v.MinIncome - 1)
			taxResp.Tax += baseCal * (v.TaxRate / 100)

		}

	}
	return taxResp, nil
}
