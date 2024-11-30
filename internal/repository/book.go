package repository

import (
	"book-fiber/domain"
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
)

type bookRepository struct {
	db *goqu.Database
}

func NewBook(con *sql.DB) domain.BookRepository {
	return &bookRepository{
		db: goqu.New("default", con),
	}
}

func (br bookRepository) FindAll(ctx context.Context) (result []domain.Book, err error) {
	dataset := br.db.From("books").Where(goqu.C("deleted_at").IsNull())
	err = dataset.ScanStructsContext(ctx, &result)
	return
}

func (br bookRepository) FindById(ctx context.Context, id string) (result domain.Book, err error) {
	dataset := br.db.From("books").Where(goqu.C("deleted_at").IsNull(), goqu.C("id").Eq(id))
	_, err = dataset.ScanStructContext(ctx, &result)
	return
}

func (br bookRepository) FindByIds(ctx context.Context, ids []string) (result []domain.Book, err error) {
	dataset := br.db.From("books").Where(goqu.C("deleted_at").IsNull(), goqu.C("id").Eq(ids))
	err = dataset.ScanStructsContext(ctx, &result)
	return
}

func (br bookRepository) FindByIsbn(ctx context.Context, isbn string) (result domain.Book, err error) {
	dataset := br.db.From("books").Where(goqu.C("deleted_at").IsNull(), goqu.C("isbn").Eq(isbn))
	_, err = dataset.ScanStructContext(ctx, &result)
	return
}

func (br bookRepository) Save(ctx context.Context, c *domain.Book) error {
	executor := br.db.Insert("books").Rows(c).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}

func (br bookRepository) Update(ctx context.Context, c *domain.Book) error {
	executor := br.db.Update("books").Where(goqu.C("id").Eq(c.Id)).Set(c).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}

func (br bookRepository) Delete(ctx context.Context, id string) error {
	executor := br.db.Update("books").Where(goqu.C("id").Eq(id)).Set(goqu.Record{"deleted_at": sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}}).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}
