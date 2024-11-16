package api

import (
	"book-fiber/domain"
	"book-fiber/dto"
	"book-fiber/internal/util"
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type customerApi struct {
	customerService domain.CustomerService
}

func NewCustomer(app *fiber.App, customerService domain.CustomerService, auzhMidd fiber.Handler) {
	ca := customerApi{customerService: customerService}

	app.Get("/customers", auzhMidd, ca.Index)
	app.Post("/customers", auzhMidd, ca.Create)
	app.Get("/customers/:id", auzhMidd, ca.Show)
	app.Put("/customers/:id", auzhMidd, ca.Update)
	app.Delete("/customers/:id", auzhMidd, ca.Delete)
}

func (ca customerApi) Index(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	res, err := ca.customerService.Index(c)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(dto.CreateResponseSuccess(res))
}

func (ca customerApi) Create(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.CreateCustomerRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(dto.CreateResponseError(err.Error()))
	}
	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(dto.CreateResponseErrorData("erro validasi", fails))
	}
	err := ca.customerService.Create(c, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.CreateResponseError(err.Error()))
	}
	return ctx.Status(http.StatusCreated).JSON(dto.CreateResponseSuccess(""))
}

func (ca customerApi) Update(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.UpdateCustomerRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(dto.CreateResponseError(err.Error()))
	}
	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(dto.CreateResponseErrorData("erro validasi", fails))
	}
	req.ID = ctx.Params("id")
	err := ca.customerService.Update(c, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.CreateResponseError(err.Error()))
	}
	return ctx.Status(http.StatusOK).JSON(dto.CreateResponseSuccess(""))
}

func (ca customerApi) Delete(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id := ctx.Params("id")
	err := ca.customerService.Delete(c, id)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.CreateResponseError(err.Error()))
	}
	return ctx.Status(http.StatusOK).JSON(dto.CreateResponseSuccess(""))
}

func (ca customerApi) Show(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id := ctx.Params("id")
	res, err := ca.customerService.Show(c, id)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(dto.CreateResponseSuccess(res))
}
