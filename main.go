package main

import (
	"book-fiber/internal/api"
	"book-fiber/internal/config"
	"book-fiber/internal/connection"
	"book-fiber/internal/repository"
	"book-fiber/internal/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cnf := config.Get()
	dbConnection := connection.GetDatabase(cnf.Database)

	app := fiber.New()
	customerRepository := repository.NewCustomer(dbConnection)
	customerService := service.NewCustomerService(customerRepository)
	api.NewCustomer(app, customerService)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
