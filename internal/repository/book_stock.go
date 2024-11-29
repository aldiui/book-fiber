package repository

import (
	"book-fiber/domain"
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
)

type bookStockRepository struct {
	db *goqu.Database
}

func NewBookStock(con *sql.DB) domain.BookStockRepository {
	return &bookStockRepository{
		db: goqu.New("default", con),
	}
}

func (bsr bookStockRepository) FindBookById(ctx context.Context, bookId string) (result []domain.BookStock, err error) {
	dataset := bsr.db.From("book_stocks").Where(goqu.C("book_id").Eq(bookId))
	err = dataset.ScanStructsContext(ctx, &result)
	return
}

func (bsr bookStockRepository) FindByBookAndCode(ctx context.Context, bookId string, code string) (result domain.BookStock, err error) {
	dataset := bsr.db.From("book_stocks").Where(goqu.C("book_id").Eq(bookId), goqu.C("code").Eq(code))
	_, err = dataset.ScanStructContext(ctx, &result)
	return
}

func (bsr bookStockRepository) Save(ctx context.Context, data []domain.BookStock) error {
	executor := bsr.db.Insert("book_stocks").Rows(data).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}

func (bsr bookStockRepository) Update(ctx context.Context, stock *domain.BookStock) error {
	executor := bsr.db.Update("book_stocks").Where(goqu.C("code").Eq(stock.Code)).Set(stock).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}

func (bsr bookStockRepository) DeleteByBookId(ctx context.Context, bookId string) error {
	executor := bsr.db.Delete("book_stocks").Where(goqu.C("book_id").Eq(bookId)).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}

func (bsr bookStockRepository) DeleteByCodes(ctx context.Context, codes []string) error {
	executor := bsr.db.Delete("book_stocks").Where(goqu.C("code").In(codes)).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}
