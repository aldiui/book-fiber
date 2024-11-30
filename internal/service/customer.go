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
		customerData = append(customerData, dto.CustomerData{Id: v.Id, Code: v.Code, Name: v.Name})
	}

	return customerData, nil
}

func (c customerService) Create(ctx context.Context, req dto.CreateCustomerRequest) error {
	existCustomer, err := c.customerRepository.FindByCode(ctx, req.Code)
	if err != nil {
		return err
	}
	if existCustomer.Id != "" {
		return errors.New("data code customer sudah ada")
	}
	customer := domain.Customer{
		Id:   uuid.NewString(),
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
	data, err := c.customerRepository.FindById(ctx, req.Id)
	if err != nil {
		return err
	}
	if data.Id == "" {
		return domain.CustomerNotFound
	}

	existCustomer, err := c.customerRepository.FindByCode(ctx, req.Code)
	if err != nil {
		return err
	}

	if existCustomer.Id != "" && existCustomer.Id != data.Id {
		return errors.New("data code customer sudah ada")
	}

	data.Code = req.Code
	data.Name = req.Name
	data.UpdatedAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}
	return c.customerRepository.Update(ctx, &data)
}

func (c customerService) Delete(ctx context.Context, id string) error {
	exist, err := c.customerRepository.FindById(ctx, id)
	if err != nil {
		return err
	}
	if exist.Id == "" {
		return domain.CustomerNotFound
	}
	return c.customerRepository.Delete(ctx, id)
}

func (c customerService) Show(ctx context.Context, id string) (dto.CustomerData, error) {
	data, err := c.customerRepository.FindById(ctx, id)
	if err != nil {
		return dto.CustomerData{}, err
	}
	if data.Id == "" {
		return dto.CustomerData{}, domain.CustomerNotFound
	}
	return dto.CustomerData{Id: data.Id, Code: data.Code, Name: data.Name}, nil
}
