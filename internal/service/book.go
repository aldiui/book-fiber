package service

import (
	"book-fiber/domain"
	"book-fiber/dto"
	"book-fiber/internal/config"
	"context"
	"database/sql"
	"errors"
	"path"
	"time"

	"github.com/google/uuid"
)

type bookService struct {
	cnf                 *config.Config
	bookRepository      domain.BookRepository
	bookStockRepository domain.BookStockRepository
	mediaRepository     domain.MediaRepository
}

func NewBookService(cnf *config.Config, bookRepository domain.BookRepository, bookStockRepository domain.BookStockRepository, mediaRepository domain.MediaRepository) domain.BookService {
	return &bookService{
		cnf:                 cnf,
		bookRepository:      bookRepository,
		bookStockRepository: bookStockRepository,
		mediaRepository:     mediaRepository,
	}
}

func (b bookService) Index(ctx context.Context) ([]dto.BookData, error) {
	books, err := b.bookRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	coverId := make([]string, 0)
	for _, v := range books {
		if v.CoverId.Valid {
			coverId = append(coverId, v.CoverId.String)
		}
	}

	covers := make(map[string]string)
	if len(coverId) > 0 {
		coversDb, _ := b.mediaRepository.FindByIds(ctx, coverId)
		for _, v := range coversDb {
			covers[v.Id] = path.Join(b.cnf.Server.Asset, v.Path)
		}
	}

	var bookData []dto.BookData
	for _, v := range books {
		var coverUrl string
		if v2, e := covers[v.CoverId.String]; e {
			coverUrl = v2
		}
		bookData = append(bookData, dto.BookData{Id: v.Id, Isbn: v.Isbn, Title: v.Title, Description: v.Description, CoverUrl: coverUrl})
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

	var coverUrl string
	if data.CoverId.Valid {
		cover, _ := b.mediaRepository.FindById(ctx, data.CoverId.String)
		if cover.Path != "" {
			coverUrl = path.Join(b.cnf.Server.Asset, cover.Path)
		}
	}

	return dto.BookShowData{
		BookData: dto.BookData{
			Id:          data.Id,
			Isbn:        data.Isbn,
			Title:       data.Title,
			CoverUrl:    coverUrl,
			Description: data.Description,
		},
		Stocks: stocksData,
	}, nil
}
