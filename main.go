package main

import (
	"net/http"

	"github.com/kanawat2566/assessment-tax/Handler"
	"github.com/kanawat2566/assessment-tax/repository"
	"github.com/labstack/echo/v4"
)

func main() {
	p, err := repository.New()
	if err != nil {
		panic(err)
	}
	e := echo.New()

	handler := Handler.NewHandler(p)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	e.GET("/taxes", handler.GetTaxes)

	e.Logger.Fatal(e.Start(":1323"))
}
