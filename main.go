package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/kanawat2566/assessment-tax/constants"
	"github.com/kanawat2566/assessment-tax/handlers"
	"github.com/kanawat2566/assessment-tax/repository"
	"github.com/kanawat2566/assessment-tax/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db, err := repository.InitDB()
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

	e.POST("/admin/deductions/:deduct_type", taxHandler.Deductions, basicAuthMiddleware)

	serverInit(e)
}

func serverInit(e *echo.Echo) {
	godotenv.Load(".env")
	port := os.Getenv("PORT")
	//e.Logger.Fatal(e.Start(fmt.Sprintf(`:%s`, port)))

	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Errorf("error when closing server %v", err)
	}

	e.Logger.Info("shutting down the server")
	fmt.Println("shutting down the server")
}

var basicAuthMiddleware = middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	if username == constants.UserAuth && password == constants.PassAuth {
		return true, nil
	}
	return false, nil
})
