package users

import (
	"context"

	"pyra/pkg/auth"
	"pyra/pkg/db"
)

type UserRepository struct {
	db db.DBTX
}

func NewRepository(db db.DBTX) auth.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) BeginTx(ctx context.Context) (db.DBTX, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *UserRepository) WithTx(tx db.DBTX) auth.UserRepository {
	return &UserRepository{
		db: tx,
	}
}

func (svc *UserRepository) FindByID(ctx context.Context, id uint64) (user auth.User, err error) {
	row := svc.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1 LIMIT 1;", id)

	err = row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return
}

const findUsserByEmailQuery = "SELECT * FROM users WHERE email = $1 LIMIT 1;"

func (svc *UserRepository) FindByEmail(ctx context.Context, email string) (user auth.User, err error) {
	row := svc.db.QueryRowContext(ctx, findUsserByEmailQuery, email)

	err = row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return
}

const createUserQuery = "INSERT INTO users (email, first_name, last_name) VALUES ($1, $2, $3) RETURNING id"

func (svc *UserRepository) Create(ctx context.Context, user auth.User) (uint64, error) {
	row := svc.db.QueryRowContext(ctx, createUserQuery, user.Email, user.FirstName, user.LastName)

	var newID uint64
	if err := row.Scan(&newID); err != nil {
		return 0, err
	}

	return newID, nil
}
