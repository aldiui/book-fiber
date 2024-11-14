package service

import (
	"book-fiber/domain"
	"book-fiber/dto"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type customerService struct {
	customerRepository domain.CustomerRepository
}

func NewCustomerService(customerRepository domain.CustomerRepository) domain.CustomerService {
	return &customerService{
		customerRepository: customerRepository,
	}
}

func (c customerService) Index(ctx context.Context) ([]dto.CustomerData, error) {
	customers, err := c.customerRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var customerData []dto.CustomerData
	for _, v := range customers {
		customerData = append(customerData, dto.CustomerData{ID: v.ID, Code: v.Code, Name: v.Name})
	}

	return customerData, nil
}

func (c customerService) Create(ctx context.Context, req dto.CreateCustomerRequest) error {
	customer := domain.Customer{
		ID:   uuid.NewString(),
		Code: req.Code,
		Name: req.Name,
		CreatedAt: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}
	return c.customerRepository.Save(ctx, &customer)
}

func (c customerService) Update(ctx context.Context, req dto.UpdateCustomerRequest) error {
	persisted, err := c.customerRepository.FindById(ctx, req.ID)
	if err != nil {
		return err
	}
	if persisted.ID == "" {
		return errors.New("data customer tidak ditemukan")
	}
	persisted.Code = req.Code
	persisted.Name = req.Name
	persisted.UpdatedAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}
	return c.customerRepository.Update(ctx, &persisted)
}

func (c customerService) Delete(ctx context.Context, id string) error {
	exist, err := c.customerRepository.FindById(ctx, id)
	if err != nil {
		return err
	}
	if exist.ID == "" {
		return errors.New("data customer tidak ditemukan")
	}
	return c.customerRepository.Delete(ctx, id)
}

func (c customerService) Show(ctx context.Context, id string) (dto.CustomerData, error) {
	persisted, err := c.customerRepository.FindById(ctx, id)
	if err != nil {
		return dto.CustomerData{}, err
	}
	if persisted.ID == "" {
		return dto.CustomerData{}, errors.New("data customer tidak ditemukan")
	}
	return dto.CustomerData{ID: persisted.ID, Code: persisted.Code, Name: persisted.Name}, nil
}
