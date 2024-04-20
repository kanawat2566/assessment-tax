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

	handler := Handler.NewHandler(p)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
