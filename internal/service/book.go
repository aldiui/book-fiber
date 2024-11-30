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
		bookData = append(bookData, dto.BookData{Id: v.Id, Isbn: v.Isbn, Title: v.Title, Description: v.Description, CoverId: v.CoverId.String})
	}

	return bookData, nil
}

func (b bookService) Create(ctx context.Context, req dto.CreateBookRequest) error {
	existBook, err := b.bookRepository.FindByIsbn(ctx, req.Isbn)
	if err != nil {
		return err
	}

	if existBook.Id != "" {
		return errors.New("data isbn book sudah ada")
	}

	coverId := sql.NullString{Valid: false, String: req.CoverId}
	if req.CoverId != "" {
		coverId.Valid = true
	}

	book := domain.Book{
		Id:          uuid.NewString(),
		Title:       req.Title,
		Description: req.Description,
		Isbn:        req.Isbn,
		CoverId:     coverId,
		CreatedAt: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}
	return b.bookRepository.Save(ctx, &book)
}

func (b bookService) Update(ctx context.Context, req dto.UpdateBookRequest) error {
	data, err := b.bookRepository.FindById(ctx, req.Id)
	if err != nil {
		return err
	}
	if data.Id == "" {
		return domain.BookNotFound
	}
	existBook, err := b.bookRepository.FindByIsbn(ctx, req.Isbn)

	if err != nil {
		return err
	}

	if existBook.Id != "" && existBook.Id != data.Id {
		return errors.New("data isbn book sudah ada")
	}

	coverId := sql.NullString{Valid: false, String: req.CoverId}
	if req.CoverId != "" {
		coverId.Valid = true
	}

	data.Title = req.Title
	data.Description = req.Description
	data.Isbn = req.Isbn
	data.UpdatedAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}
	data.CoverId = coverId
	return b.bookRepository.Update(ctx, &data)
}

func (b bookService) Delete(ctx context.Context, id string) error {
	exist, err := b.bookRepository.FindById(ctx, id)
	if err != nil {
		return err
	}
	if exist.Id == "" {
		return domain.BookNotFound
	}
	return b.bookRepository.Delete(ctx, id)
}

func (b bookService) Show(ctx context.Context, id string) (dto.BookShowData, error) {
	data, err := b.bookRepository.FindById(ctx, id)
	if err != nil {
		return dto.BookShowData{}, err
	}
	if data.Id == "" {
		return dto.BookShowData{}, domain.BookNotFound
	}

	stocks, er := b.bookStockRepository.FindBookById(ctx, data.Id)

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
			Id:          data.Id,
			Title:       data.Title,
			Description: data.Description,
			Isbn:        data.Isbn,
			CoverId:     data.CoverId.String,
		},
		Stocks: stocksData,
	}, nil
}
