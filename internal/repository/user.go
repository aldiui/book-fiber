package repository

import (
	"book-fiber/domain"
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
)

type userRepository struct {
	db *goqu.Database
}

func NewUser(con *sql.DB) domain.UserRepository {
	return &userRepository{
		db: goqu.New("default", con),
	}
}

func (u userRepository) FindByEmail(ctx context.Context, email string) (result domain.User, err error) {
	dataset := u.db.From("users").Where(goqu.C("deleted_at").IsNull(), goqu.C("email").Eq(email))
	_, err = dataset.ScanStructContext(ctx, &result)
	return
}

func (u userRepository) Save(ctx context.Context, c *domain.User) error {
	executor := u.db.Insert("users").Rows(c).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}
