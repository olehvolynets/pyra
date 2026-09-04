package auth

import (
	"context"
	"time"

	"pyra/pkg/db"
)

type User struct {
	ID uint64

	Email string

	FirstName string
	LastName  string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepository interface {
	db.Repository[UserRepository]

	FindByID(ctx context.Context, id uint64) (user User, err error)
	FindByEmail(ctx context.Context, email string) (user User, err error)
	Create(ctx context.Context, user User) (uint64, error)
}
