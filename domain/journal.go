package domain

import (
	"book-fiber/dto"
	"context"
	"database/sql"
)

const (
	JournalStatusInProgress = "IN_PROGRESS"
	JournalStatusCompleted  = "COMPLETED"
)

type Journal struct {
	ID         string       `db:"id"`
	BookId     string       `db:"book_id"`
	StockCode  string       `db:"stock_code"`
	CustomerId string       `db:"customer_id"`
	Status     string       `db:"status"`
	BorrowedAt sql.NullTime `db:"borrowed_at"`
	ReturnedAt sql.NullTime `db:"returned_at"`
}

type JournalSearch struct {
	CustomerId string
	Status     string
}

type JournalRepository interface {
	Find(ctx context.Context, search JournalSearch) ([]Journal, error)
	FindById(ctx context.Context, id string) (Journal, error)
	Save(ctx context.Context, journal *Journal) error
	Update(ctx context.Context, journal *Journal) error
}
type JournalService interface {
	Index(ctx context.Context, search JournalSearch) ([]dto.JournalData, error)
	Create(ctx context.Context, req dto.CreateJournalRequest) error
	Return(ctx context.Context, req dto.ReturnJournalRequest) error
}
