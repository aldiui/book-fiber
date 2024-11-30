package api

import (
	"book-fiber/domain"
	"book-fiber/dto"
	"book-fiber/internal/config"
	"context"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type mediaApi struct {
	cnf          *config.Config
	mediaService domain.MediaService
}

func NewMedia(app *fiber.App, cnf *config.Config, mediaService domain.MediaService, authMidd fiber.Handler) {
	ma := mediaApi{mediaService: mediaService}
	media := app.Group("/media", authMidd)

	media.Post("", ma.Create)
	media.Static("", cnf.Storage.BasePath)
}

func (ma mediaApi) Create(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	file, err := ctx.FormFile("media")
	if err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(dto.CreateResponseError(err.Error()))
	}

	filename := uuid.NewString() + file.Filename
	path := filepath.Join(ma.cnf.Storage.BasePath, filename)
	err = ctx.SaveFile(file, path)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.CreateResponseError(err.Error()))
	}

	res, err := ma.mediaService.Create(c, dto.CreateMediaRequest{Path: filename})

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.Status(http.StatusCreated).JSON(dto.CreateResponseSuccess(res))
}
