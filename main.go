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
	bookRepository := repository.NewBook(dbConnection)
	bookStockRepository := repository.NewBookStock(dbConnection)
	journalRepository := repository.NewJournal(dbConnection)
	mediaRepository := repository.NewMedia(dbConnection)
	chargeRepository := repository.NewCharge(dbConnection)

	customerService := service.NewCustomerService(customerRepository)
	authService := service.NewAuthService(cnf, userRepository)
	bookService := service.NewBookService(cnf, bookRepository, bookStockRepository, mediaRepository)
	bookStockService := service.NewBookStockService(bookRepository, bookStockRepository)
	journalService := service.NewJournalService(journalRepository, bookRepository, bookStockRepository, customerRepository, chargeRepository)
	mediaService := service.NewMediaService(cnf, mediaRepository)

	api.NewAuth(app, authService)
	api.NewCustomer(app, customerService, jwtMidd)
	api.NewBook(app, bookService, jwtMidd)
	api.NewBookStock(app, bookStockService, jwtMidd)
	api.NewJournal(app, journalService, jwtMidd)
	api.NewMedia(app, cnf, mediaService, jwtMidd)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
