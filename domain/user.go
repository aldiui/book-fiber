package domain

import (
	"context"
	"database/sql"
)

type User struct {
	ID        string       `db:"id"`
	Name      string       `db:"name"`
	Email     string       `db:"email"`
	Password  string       `db:"password"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	Save(ctx context.Context, c *User) error
}
