package repository

import (
	"book-fiber/domain"
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
)

type journalRepository struct {
	db *goqu.Database
}

func NewJournal(con *sql.DB) domain.JournalRepository {
	return &journalRepository{
		db: goqu.New("default", con),
	}
}

func (j *journalRepository) Find(ctx context.Context, search domain.JournalSearch) (result []domain.Journal, err error) {
	dataset := j.db.From("journals")
	if search.CustomerId != "" {
		dataset = dataset.Where(goqu.C("customer_id").Eq(search.CustomerId))
	}
	if search.Status != "" {
		dataset = dataset.Where(goqu.C("status").Eq(search.Status))
	}
	err = dataset.ScanStructsContext(ctx, &result)
	return
}

func (j *journalRepository) FindById(ctx context.Context, id string) (result domain.Journal, err error) {
	dataset := j.db.From("journals").Where(goqu.C("id").Eq(id))
	_, err = dataset.ScanStructContext(ctx, &result)
	return
}

func (j *journalRepository) Save(ctx context.Context, journal *domain.Journal) error {
	executor := j.db.Insert("journals").Rows(journal).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}

func (j *journalRepository) Update(ctx context.Context, journal *domain.Journal) error {
	executor := j.db.Update("journals").Where(goqu.C("id").Eq(journal.ID)).Set(journal).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}
