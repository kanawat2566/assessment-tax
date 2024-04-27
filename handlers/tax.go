package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	ct "github.com/kanawat2566/assessment-tax/constants"
	models "github.com/kanawat2566/assessment-tax/model"
	"github.com/kanawat2566/assessment-tax/services"
	"github.com/labstack/echo/v4"
)

type taxHandler struct {
	serv services.TaxService
}

type CustomValidator struct {
	Validator *validator.Validate
}

func NewHandler(s services.TaxService) *taxHandler {
	return &taxHandler{serv: s}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	return nil
}

func (h *taxHandler) CalculationsHandler(c echo.Context) error {

	rq := new(models.TaxRequest)
	if err := c.Bind(rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := h.serv.TaxCalculations(rq)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if res.TaxRefund > 0 {
		return c.JSON(http.StatusOK, res)
	} else {
		tax := models.TaxLevelReponse{Tax: res.Tax, TaxLevels: res.TaxLevels}
		return c.JSON(http.StatusOK, tax)
	}

}

func (h *taxHandler) Deductions(c echo.Context) error {
	d := c.Param("type")

	fmt.Printf("type= %v\n", d)

	dd := ct.Deductions[d]
	if len(dd.Name) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid deduction type")

	}
	rq := new(models.DeductRequest)

	if err := c.Bind(rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	if err := c.Validate(rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	res, err := h.serv.SetAdminDeductions(ct.Deduction{Type: dd.Type, Amount: rq.Amount})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	response := map[string]interface{}{
		res.Name: res.Amount,
	}
	return c.JSON(http.StatusOK, response)
}
