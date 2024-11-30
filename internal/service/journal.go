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

type journalService struct {
	journalRepository   domain.JournalRepository
	bookRepository      domain.BookRepository
	bookStockRepository domain.BookStockRepository
	customerRepository  domain.CustomerRepository
	chargeRepository    domain.ChargeRepository
}

func NewJournalService(journalRepository domain.JournalRepository, bookRepository domain.BookRepository, bookStockRepository domain.BookStockRepository, customerRepository domain.CustomerRepository, chargeRepository domain.ChargeRepository) domain.JournalService {
	return &journalService{
		journalRepository:   journalRepository,
		bookRepository:      bookRepository,
		bookStockRepository: bookStockRepository,
		customerRepository:  customerRepository,
		chargeRepository:    chargeRepository,
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
	if len(customerId) > 0 {
		customerDb, _ := j.customerRepository.FindByIds(ctx, customerId)
		for _, v := range customerDb {
			customers[v.Id] = v
		}
	}
	books := make(map[string]domain.Book)

	if len(bookId) > 0 {
		bookDb, _ := j.bookRepository.FindByIds(ctx, bookId)
		for _, v := range bookDb {
			books[v.Id] = v
		}
	}
	result := make([]dto.JournalData, 0)
	for _, v := range journals {
		book := dto.BookData{}
		if v2, e := books[v.BookId]; e {
			book = dto.BookData{
				Id:          v2.Id,
				Isbn:        v2.Isbn,
				Title:       v2.Title,
				Description: v2.Description,
			}
		}

		customer := dto.CustomerData{}
		if v2, e := customers[v.CustomerId]; e {
			customer = dto.CustomerData{
				Id:   v2.Id,
				Code: v2.Code,
				Name: v2.Name,
			}
		}

		result = append(result, dto.JournalData{
			Id:         v.Id,
			BookStock:  v.StockCode,
			Book:       book,
			Customer:   customer,
			Status:     v.Status,
			BorrowedAt: v.BorrowedAt.Time,
			ReturnedAt: v.ReturnedAt.Time,
		})

	}
	return result, nil
}

func (j *journalService) Create(ctx context.Context, req dto.CreateJournalRequest) error {
	book, err := j.bookRepository.FindById(ctx, req.BookId)
	if err != nil {
		return err
	}

	if book.Id == "" {
		return domain.BookNotFound
	}

	stock, err := j.bookStockRepository.FindByBookAndCode(ctx, req.BookId, req.BookStock)
	if err != nil {
		return err
	}

	if stock.Code == "" {
		return domain.BookNotFound
	}

	if stock.Status != domain.BookStockStatusAvailable {
		return errors.New("stock buku sudah dipinjam sebelumnya")
	}

	journal := domain.Journal{
		Id:         uuid.NewString(),
		BookId:     req.BookId,
		StockCode:  req.BookStock,
		CustomerId: req.CustomerId,
		Status:     domain.JournalStatusInProgress,
		DueAt: sql.NullTime{
			Valid: true,
			Time:  time.Now().Add(time.Hour * 24 * 7),
		},
		BorrowedAt: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}

	err = j.journalRepository.Save(ctx, &journal)
	if err != nil {
		return err
	}
	stock.Status = domain.BookStockStatusBorrowed
	stock.BorrowedAt = journal.BorrowedAt
	stock.BorrowerId = sql.NullString{
		Valid:  true,
		String: journal.CustomerId,
	}

	return j.bookStockRepository.Update(ctx, &stock)
}

func (j *journalService) Return(ctx context.Context, req dto.ReturnJournalRequest) error {
	journal, err := j.journalRepository.FindById(ctx, req.JournalId)
	if err != nil {
		return err
	}

	if journal.Id == "" {
		return domain.JournalNotFound
	}

	stock, err := j.bookStockRepository.FindByBookAndCode(ctx, journal.BookId, journal.StockCode)
	if err != nil {
		return err
	}

	if stock.Code != "" {
		stock.Status = domain.BookStockStatusAvailable
		stock.BorrowedAt = sql.NullTime{Valid: false}
		stock.BorrowerId = sql.NullString{Valid: false}
		err = j.bookStockRepository.Update(ctx, &stock)
		if err != nil {
			return err
		}
	}

	journal.Status = domain.JournalStatusCompleted
	journal.ReturnedAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}
	err = j.journalRepository.Update(ctx, &journal)
	if err != nil {
		return err
	}
	hoursLate := time.Now().Sub(journal.DueAt.Time).Hours()
	if hoursLate >= 24 {
		daysLate := int(hoursLate / 24)
		charge := domain.Charge{
			Id:           uuid.NewString(),
			JournalId:    journal.Id,
			DaysLate:     daysLate,
			DailyLateFee: 5000,
			Total:        daysLate * 5000,
			UserId:       req.UserId,
			CreatedAt: sql.NullTime{
				Valid: true,
				Time:  time.Now(),
			},
		}
		err = j.chargeRepository.Save(ctx, &charge)
	}
	return err
}
