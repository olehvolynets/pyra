// Package products
package products

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"pyra/pkg/db"
	"pyra/pkg/log"
	"pyra/pkg/nutrition"
)

type Repository struct {
	db db.DBTX
}

func NewRepository(db db.DBTX) nutrition.ProductRepository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) BeginTx(ctx context.Context) (db.DBTX, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *Repository) WithTx(tx db.DBTX) nutrition.ProductRepository {
	return NewRepository(tx)
}

func (r *Repository) FindByRef(ctx context.Context, ref nutrition.ProductRef) (nutrition.Product, error) {
	row := r.db.QueryRowContext(ctx, findByRefQuery, ref.UID, ref.Version)

	return r.scanProductRow(row)
}

func (r *Repository) Versions(ctx context.Context, uid nutrition.ProductUID) ([]nutrition.Product, error) {
	rows, err := r.db.QueryContext(ctx, productVersionsQuery, uid)
	if err != nil {
		return nil, err
	}
	defer closeRows(ctx, rows)

	return r.scanProducts(rows)
}

func (r *Repository) FindAllByRefs(ctx context.Context, refs []nutrition.ProductRef) ([]nutrition.Product, error) {
	rows, err := r.db.QueryContext(ctx, findAllByIDsQuery, refs)
	if err != nil {
		return nil, err
	}
	defer closeRows(ctx, rows)

	return r.scanProducts(rows)
}

func (r *Repository) IsNameTaken(ctx context.Context, name nutrition.ProductName) (bool, error) {
	row := r.db.QueryRowContext(ctx, nameTakenQuery, name)

	var one int
	err := row.Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	return one == 1, err
}

func (r *Repository) ForDish(
	ctx context.Context,
	uid nutrition.DishUID,
	version nutrition.DishVersion,
) ([]nutrition.Product, error) {
	rows, err := r.db.QueryContext(ctx, productsForDishQuery, uid, version)
	if err != nil {
		return nil, err
	}
	defer closeRows(ctx, rows)

	return r.scanProducts(rows)
}

func (r *Repository) Index(ctx context.Context) ([]nutrition.Product, error) {
	rows, err := r.db.QueryContext(ctx, indexProductsQuery)
	if err != nil {
		return nil, err
	}
	defer closeRows(ctx, rows)

	return r.scanProducts(rows)
}

func (r *Repository) Create(ctx context.Context, p *nutrition.Product) error {
	row := r.db.QueryRowContext(ctx, createProductQuery,
		p.UID, p.Name, p.Calories, p.Proteins, p.Fats, p.Carbs)

	return row.Scan(&p.Version, &p.CreatedAt, &p.UpdatedAt)
}

func (r *Repository) CreateVersion(ctx context.Context, p *nutrition.Product) error {
	row := r.db.QueryRowContext(ctx, createProductVersionQuery,
		p.UID, p.Name, p.Calories, p.Proteins, p.Fats, p.Carbs)

	return row.Scan(&p.Version, &p.CreatedAt, &p.UpdatedAt)
}

func (r *Repository) Delete(ctx context.Context, ref nutrition.ProductRef) error {
	res, err := r.db.ExecContext(ctx, deleteByRefQuery, ref.UID, ref.Version)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, product *nutrition.Product) error {
	result, err := r.db.ExecContext(ctx, updateProductQuery,
		product.UID, product.Version, product.Name,
		product.Calories, product.Proteins, product.Fats, product.Carbs)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if rowsAffected != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", rowsAffected)
	}

	return nil
}

func (r *Repository) Archive(ctx context.Context, ref nutrition.ProductRef, ts time.Time) error {
	result, err := r.db.ExecContext(ctx, "UPDATE products SET archived_at = $1 WHERE uid = $2 AND version = $3",
		ts, ref.UID, ref.Version)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", rowsAffected)
	}

	return nil
}

func (r *Repository) Search(ctx context.Context, searchStr string) ([]nutrition.Product, error) {
	rows, err := r.db.QueryContext(ctx, searchProductsQuery, searchStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}
	defer closeRows(ctx, rows)

	return r.scanProducts(rows)
}

func (r *Repository) MaxVersion(
	ctx context.Context,
	uid nutrition.ProductUID,
) (nutrition.ProductVersion, error) {
	row := r.db.QueryRowContext(ctx, maxProductVersionQuery, uid)

	var version nutrition.ProductVersion
	if err := row.Scan(&version); err != nil {
		return nutrition.ProductVersion(-1), err
	}

	return version, nil
}

func (r *Repository) UsedInDishes(ctx context.Context, ref nutrition.ProductRef) (bool, error) {
	row := r.db.QueryRowContext(ctx, usedInDishesQuery, nutrition.IngredientProduct, ref.UID, ref.Version)

	var one int
	err := row.Scan(&one)	

	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}

	return one == 1, err
}

func (r *Repository) CountAll(ctx context.Context) (n int, err error) {
	row := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM products;")
	err = row.Scan(&n)

	return
}

func (r *Repository) scanProducts(rows *sql.Rows) ([]nutrition.Product, error) {
	products := make([]nutrition.Product, 0)

	for rows.Next() {
		product, err := r.scanProductRows(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (r *Repository) scanProductRows(rows *sql.Rows) (nutrition.Product, error) {
	product := nutrition.Product{}

	err := rows.Scan(&product.UID, &product.Version, &product.Name,
		&product.Calories, &product.Proteins, &product.Fats, &product.Carbs,
		&product.CreatedAt, &product.UpdatedAt)

	return product, err
}

func (r *Repository) scanProductRow(row *sql.Row) (nutrition.Product, error) {
	product := nutrition.Product{}

	err := row.Scan(&product.UID, &product.Version, &product.Name,
		&product.Calories, &product.Proteins, &product.Fats, &product.Carbs,
		&product.CreatedAt, &product.UpdatedAt)

	return product, err
}

func closeRows(ctx context.Context, rows *sql.Rows) {
	if closeErr := rows.Close(); closeErr != nil {
		log.FromContext(ctx).WarnContext(ctx, "failed to close TX", "error", closeErr)
	}
}
