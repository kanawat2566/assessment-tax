package Handler

import (
	"fmt"

	"github.com/kanawat2566/assessment-tax/repository"
	"github.com/labstack/echo/v4"
)

type handler struct {
	repo repository.TaxRepository
}

func NewHandler(r repository.TaxRepository) *handler {
	return &handler{repo: r}
}

func (h *handler) GetTaxes(c echo.Context) error {
	taxs, err := h.repo.GetAll()
	if err != nil {
		fmt.Println(err)
		return c.JSON(500, "Internal Server Error")

	}
	return c.JSON(200, taxs)
}
