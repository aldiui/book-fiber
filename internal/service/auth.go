package service

import (
	"book-fiber/domain"
	"book-fiber/dto"
	"book-fiber/internal/config"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	conf           *config.Config
	userRepository domain.UserRepository
}

func NewAuthService(cnf *config.Config, userRepository domain.UserRepository) domain.AuthService {
	return authService{
		conf:           cnf,
		userRepository: userRepository,
	}
}

func (a authService) Login(ctx context.Context, req dto.AuthRequest) (dto.AuthResponse, error) {
	user, err := a.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if user.ID == "" {
		return dto.AuthResponse{}, errors.New("data user tidak ditemukan")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return dto.AuthResponse{}, errors.New("password salah")
	}
	claim := jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Duration(a.conf.Jwt.Exp) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString([]byte(a.conf.Jwt.Key))
	if err != nil {
		return dto.AuthResponse{}, errors.New("gagal membuat token")
	}
	return dto.AuthResponse{Token: tokenStr}, nil
}

func (a authService) Register(ctx context.Context, req dto.AuthRegisterRequest) (dto.AuthResponse, error) {
	existingUser, err := a.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if existingUser.ID != "" {
		return dto.AuthResponse{}, errors.New("data user sudah terdaftar")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.AuthResponse{}, errors.New("gagal membuat password")
	}

	createUuid := uuid.NewString()
	user := domain.User{
		ID:       createUuid,
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		CreatedAt: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}

	err = a.userRepository.Save(ctx, &user)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	claim := jwt.MapClaims{
		"id":  createUuid,
		"exp": time.Now().Add(time.Duration(a.conf.Jwt.Exp) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString([]byte(a.conf.Jwt.Key))
	if err != nil {
		return dto.AuthResponse{}, errors.New("gagal membuat token")
	}
	return dto.AuthResponse{Token: tokenStr}, nil
}
