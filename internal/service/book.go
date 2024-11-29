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

type bookService struct {
	bookRepository      domain.BookRepository
	bookStockRepository domain.BookStockRepository
}

func NewBookService(bookRepository domain.BookRepository, bookStockRepository domain.BookStockRepository) domain.BookService {
	return &bookService{
		bookRepository:      bookRepository,
		bookStockRepository: bookStockRepository,
	}
}

func (b bookService) Index(ctx context.Context) ([]dto.BookData, error) {
	books, err := b.bookRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var bookData []dto.BookData
	for _, v := range books {
		bookData = append(bookData, dto.BookData{ID: v.ID, Isbn: v.Isbn, Title: v.Title, Description: v.Description})
	}

	return bookData, nil
}

func (b bookService) Create(ctx context.Context, req dto.CreateBookRequest) error {
	existBook, err := b.bookRepository.FindByIsbn(ctx, req.Isbn)
	if err != nil {
		return err
	}

	if existBook.ID != "" {
		return errors.New("data isbn book sudah ada")
	}

	book := domain.Book{
		ID:          uuid.NewString(),
		Title:       req.Title,
		Description: req.Description,
		Isbn:        req.Isbn,
		CreatedAt: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}
	return b.bookRepository.Save(ctx, &book)
}

func (b bookService) Update(ctx context.Context, req dto.UpdateBookRequest) error {
	data, err := b.bookRepository.FindById(ctx, req.ID)
	if err != nil {
		return err
	}
	if data.ID == "" {
		return errors.New("data book tidak ditemukan")
	}
	existBook, err := b.bookRepository.FindByIsbn(ctx, req.Isbn)

	if err != nil {
		return err
	}

	if existBook.ID != "" && existBook.ID != data.ID {
		return errors.New("data isbn book sudah ada")
	}

	data.Title = req.Title
	data.Description = req.Description
	data.Isbn = req.Isbn
	data.UpdatedAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}
	return b.bookRepository.Update(ctx, &data)
}

func (b bookService) Delete(ctx context.Context, id string) error {
	exist, err := b.bookRepository.FindById(ctx, id)
	if err != nil {
		return err
	}
	if exist.ID == "" {
		return errors.New("data book tidak ditemukan")
	}
	return b.bookRepository.Delete(ctx, id)
}

func (b bookService) Show(ctx context.Context, id string) (dto.BookShowData, error) {
	data, err := b.bookRepository.FindById(ctx, id)
	if err != nil {
		return dto.BookShowData{}, err
	}
	if data.ID == "" {
		return dto.BookShowData{}, domain.BookNotFound
	}

	stocks, er := b.bookStockRepository.FindBookById(ctx, data.ID)

	if er != nil {
		return dto.BookShowData{}, er
	}

	stocksData := make([]dto.BookStockData, 0)
	for _, v := range stocks {
		stocksData = append(stocksData, dto.BookStockData{
			Code:   v.Code,
			Status: v.Status,
		})
	}

	return dto.BookShowData{
		BookData: dto.BookData{
			ID:          data.ID,
			Title:       data.Title,
			Description: data.Description,
			Isbn:        data.Isbn,
		},
		Stocks: stocksData,
	}, nil
}
