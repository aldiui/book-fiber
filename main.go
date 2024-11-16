package main

import (
	"book-fiber/dto"
	"book-fiber/internal/api"
	"book-fiber/internal/config"
	"book-fiber/internal/connection"
	"book-fiber/internal/repository"
	"book-fiber/internal/service"

	jwtMid "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cnf := config.Get()
	dbConnection := connection.GetDatabase(cnf.Database)

	app := fiber.New()
	jwtMidd := jwtMid.New(jwtMid.Config{
		SigningKey: jwtMid.SigningKey{Key: []byte(cnf.Jwt.Key)},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusUnauthorized).JSON(dto.CreateResponseError("api memerlukan token, silahkan login"))
		},
	})
	customerRepository := repository.NewCustomer(dbConnection)
	userRepository := repository.NewUser(dbConnection)
	customerService := service.NewCustomerService(customerRepository)
	authService := service.NewAuthService(cnf, userRepository)
	api.NewAuth(app, authService)
	api.NewCustomer(app, customerService, jwtMidd)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
