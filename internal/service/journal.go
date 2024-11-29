package service

import (
	"book-fiber/domain"
	"book-fiber/dto"
	"context"
)

type journalService struct {
	journalRepository   domain.JournalRepository
	bookRepository      domain.BookRepository
	bookStockRepository domain.BookStockRepository
	customerRepository  domain.CustomerRepository
}

func NewJournalService(journalRepository domain.JournalRepository, bookRepository domain.BookRepository, bookStockRepository domain.BookStockRepository, customerRepository domain.CustomerRepository) domain.JournalService {
	return &journalService{
		journalRepository:   journalRepository,
		bookRepository:      bookRepository,
		bookStockRepository: bookStockRepository,
		customerRepository:  customerRepository,
	}
}

func (j *journalService) Index(ctx context.Context, search domain.JournalSearch) ([]dto.JournalData, error) {
	journals, err := j.journalRepository.Find(ctx, search)
	if err != nil {
		return nil, err
	}
	customerId := make([]string, 0)
	bookId := make([]string, 0)
	for _, v := range journals {
		customerId = append(customerId, v.CustomerId)
		bookId = append(bookId, v.BookId)
	}
	customers := make(map[string]domain.Customer)
	books := make(map[string]domain.Book)
}

func (j *journalService) Create(ctx context.Context, req dto.CreateJournalRequest) error {
	panic("unimplemented")
}

func (j *journalService) Return(ctx context.Context, req dto.ReturnJournalRequest) error {
	panic("unimplemented")
}
