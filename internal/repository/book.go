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

func (cr bookRepository) FindAll(ctx context.Context) (result []domain.Book, err error) {
	dataset := cr.db.From("books").Where(goqu.C("deleted_at").IsNull())
	err = dataset.ScanStructsContext(ctx, &result)
	return
}

func (cr bookRepository) FindById(ctx context.Context, id string) (result domain.Book, err error) {
	dataset := cr.db.From("books").Where(goqu.C("deleted_at").IsNull(), goqu.C("id").Eq(id))
	_, err = dataset.ScanStructContext(ctx, &result)
	return
}

func (cr bookRepository) FindByTitle(ctx context.Context, title string) (result domain.Book, err error) {
	dataset := cr.db.From("books").Where(goqu.C("deleted_at").IsNull(), goqu.C("title").Eq(title))
	_, err = dataset.ScanStructContext(ctx, &result)
	return
}

func (cr bookRepository) Save(ctx context.Context, c *domain.Book) error {
	executor := cr.db.Insert("books").Rows(c).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}

func (cr bookRepository) Update(ctx context.Context, c *domain.Book) error {
	executor := cr.db.Update("books").Where(goqu.C("id").Eq(c.ID)).Set(c).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}

func (cr bookRepository) Delete(ctx context.Context, id string) error {
	executor := cr.db.Update("books").Where(goqu.C("id").Eq(id)).Set(goqu.Record{"deleted_at": sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}}).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}
