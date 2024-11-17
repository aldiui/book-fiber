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
	bookRepository domain.BookRepository
}

func NewBookService(bookRepository domain.BookRepository) domain.BookService {
	return &bookService{
		bookRepository: bookRepository,
	}
}

func (c bookService) Index(ctx context.Context) ([]dto.BookData, error) {
	books, err := c.bookRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var bookData []dto.BookData
	for _, v := range books {
		bookData = append(bookData, dto.BookData{ID: v.ID, Title: v.Title, Description: v.Description})
	}

	return bookData, nil
}

func (c bookService) Create(ctx context.Context, req dto.CreateBookRequest) error {
	existBook, err := c.bookRepository.FindByTitle(ctx, req.Title)
	if err != nil {
		return err
	}

	if existBook.ID != "" {
		return errors.New("data title book sudah ada")
	}

	book := domain.Book{
		ID:          uuid.NewString(),
		Title:       req.Title,
		Description: req.Description,
		CreatedAt: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}
	return c.bookRepository.Save(ctx, &book)
}

func (c bookService) Update(ctx context.Context, req dto.UpdateBookRequest) error {
	persisted, err := c.bookRepository.FindById(ctx, req.ID)
	if err != nil {
		return err
	}
	if persisted.ID == "" {
		return errors.New("data book tidak ditemukan")
	}
	existBook, err := c.bookRepository.FindByTitle(ctx, req.Title)
	if err != nil {
		return err
	}

	if existBook.ID != "" && existBook.ID != persisted.ID {
		return errors.New("data title book sudah ada")
	}

	persisted.Title = req.Title
	persisted.Description = req.Description
	persisted.UpdatedAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}
	return c.bookRepository.Update(ctx, &persisted)
}

func (c bookService) Delete(ctx context.Context, id string) error {
	exist, err := c.bookRepository.FindById(ctx, id)
	if err != nil {
		return err
	}
	if exist.ID == "" {
		return errors.New("data book tidak ditemukan")
	}
	return c.bookRepository.Delete(ctx, id)
}

func (c bookService) Show(ctx context.Context, id string) (dto.BookData, error) {
	persisted, err := c.bookRepository.FindById(ctx, id)
	if err != nil {
		return dto.BookData{}, err
	}
	if persisted.ID == "" {
		return dto.BookData{}, errors.New("data book tidak ditemukan")
	}
	return dto.BookData{ID: persisted.ID, Title: persisted.Title, Description: persisted.Description}, nil
}
