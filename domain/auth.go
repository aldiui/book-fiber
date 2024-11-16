package domain

import (
	"book-fiber/dto"
	"context"
)

type AuthService interface {
	Login(ctx context.Context, req dto.AuthRequest) (dto.AuthResponse, error)
	Register(ctx context.Context, req dto.AuthRegisterRequest) (dto.AuthResponse, error)
}
