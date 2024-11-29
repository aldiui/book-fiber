package api

import (
	"book-fiber/domain"
	"book-fiber/dto"
	"book-fiber/internal/util"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type bookStockApi struct {
	bookStockService domain.BookStockService
}

func NewBookStock(app *fiber.App, bookStockService domain.BookStockService, authMidd fiber.Handler) {
	bsa := bookStockApi{bookStockService: bookStockService}

	bookStock := app.Group("/book-stocks", authMidd)
	bookStock.Post("", bsa.Create)
	bookStock.Delete("", bsa.Delete)
}

func (bsa bookStockApi) Create(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.CreateBookStockRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(dto.CreateResponseError("Format JSON tidak valid"))
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(dto.CreateResponseErrorData("Validasi gagal", fails))
	}

	err := bsa.bookStockService.Create(c, req)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.CreateResponseError("Gagal membuat data stok buku"))
	}

	return ctx.Status(http.StatusCreated).JSON(dto.CreateResponseSuccess(""))
}

func (bsa bookStockApi) Delete(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	codeStr := ctx.Query("code")
	if codeStr == "" {
		return ctx.Status(http.StatusBadRequest).JSON(dto.CreateResponseError("Parameter codes wajib diisi dan tidak boleh kosong"))
	}
	codes := strings.Split(codeStr, ";")

	err := bsa.bookStockService.Delete(c, dto.DeleteBookStockRequest{Codes: codes})
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.CreateResponseError("Gagal menghapus data stok buku"))
	}

	return ctx.Status(http.StatusOK).JSON(dto.CreateResponseSuccess(""))
}
