package main

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/kanawat2566/assessment-tax/handlers"
	"github.com/kanawat2566/assessment-tax/repository"
	"github.com/kanawat2566/assessment-tax/services"
	"github.com/labstack/echo/v4"
)

func main() {
	db,err := repository.InitDB()
	if err != nil {
		panic(err)
	}
	p := repository.New(db)

	e := echo.New()
	e.Validator = &handlers.CustomValidator{Validator: validator.New()}

	serv := services.NewServices(p)
	taxHandler := handlers.NewHandler(serv)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	e.POST("/tax/calculations", taxHandler.CalculationsHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
