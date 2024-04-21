package handlers

import (
	"net/http"

	"github.com/go-playground/validator"
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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

	return c.JSON(http.StatusOK, res)

}
