package service

import (
	"book-fiber/domain"
	"book-fiber/dto"
	"context"
)

type bookStockService struct {
	bookRepository      domain.BookRepository
	bookStockRepository domain.BookStockRepository
}

func NewBookStockService(bookRepository domain.BookRepository, bookStockRepository domain.BookStockRepository) domain.BookStockService {
	return &bookStockService{
		bookRepository:      bookRepository,
		bookStockRepository: bookStockRepository,
	}
}

func (b bookStockService) Create(ctx context.Context, req dto.CreateBookStockRequest) error {
	book, err := b.bookRepository.FindById(ctx, req.BookId)
	if err != nil {
		return err
	}
	if book.Id == "" {
		return domain.BookNotFound
	}
	stocks := make([]domain.BookStock, 0)
	for _, v := range req.Codes {
		stocks = append(stocks, domain.BookStock{
			Code:   v,
			BookId: req.BookId,
			Status: domain.BookStockStatusAvailable,
		})
	}

	return b.bookStockRepository.Save(ctx, stocks)
}

func (b bookStockService) Delete(ctx context.Context, req dto.DeleteBookStockRequest) error {
	return b.bookStockRepository.DeleteByCodes(ctx, req.Codes)
}
