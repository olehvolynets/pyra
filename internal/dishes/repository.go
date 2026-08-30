package dishes

import (
	"context"
	"database/sql"
	"errors"

	"pyra/pkg/db"
	"pyra/pkg/log"
	"pyra/pkg/nutrition"
)

type Repository struct {
	db db.DBTX
}

func NewRepository(db db.DBTX) nutrition.DishRepository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) BeginTx(ctx context.Context) (db.DBTX, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *Repository) WithTx(tx db.DBTX) nutrition.DishRepository {
	return NewRepository(tx)
}

func (r *Repository) Index(ctx context.Context) ([]nutrition.Dish, error) {
	rows, err := r.db.QueryContext(ctx, listDishesQuery)
	if err != nil {
		return nil, err
	}

	return r.scanAllRows(rows)
}

func (r *Repository) FindByRef(ctx context.Context, ref nutrition.DishRef) (nutrition.Dish, error) {
	row := r.db.QueryRowContext(ctx, dishByRefQuery, ref.UID, ref.Version)

	return r.scanRow(row)
}

func (r *Repository) Versions(ctx context.Context, uid nutrition.DishUID) ([]nutrition.Dish, error) {
	rows, err := r.db.QueryContext(ctx, dishVersionsQuery, uid)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.FromContext(ctx).WarnContext(ctx, "failed to close TX", "error", closeErr)
		}
	}()

	return r.scanAllRows(rows)
}

func (r *Repository) FindAllByProductRef(ctx context.Context, ref nutrition.ProductRef) ([]nutrition.Dish, error) {
	rows, err := r.db.QueryContext(ctx, dishesByProductQuery, ref.UID, ref.Version)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.FromContext(ctx).WarnContext(ctx, "failed to close rows", "error", closeErr)
		}
	}()

	return r.scanAllRows(rows)
}

func (r *Repository) FindAllByRefs(ctx context.Context, refs []nutrition.DishRef) ([]nutrition.Dish, error) {
	panic("not implemented")
}

func (r *Repository) IsNameTaken(ctx context.Context, name nutrition.DishName, uid nutrition.DishUID) (bool, error) {
	row := r.db.QueryRowContext(ctx, isDishNameTakenQuery, name, uid)

	var one int
	err := row.Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	return one == 1, err
}

func (r *Repository) Create(ctx context.Context, dish *nutrition.Dish) error {
	row := r.db.QueryRowContext(ctx, createDishQuery,
		dish.UID, dish.Version, dish.Name,
		dish.Calories, dish.Proteins, dish.Fats, dish.Carbs)

	return row.Scan(&dish.CreatedAt, &dish.UpdatedAt)
}

func (r *Repository) Delete(ctx context.Context, ref nutrition.DishRef) error {
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

func (r *Repository) scanRow(row *sql.Row) (nutrition.Dish, error) {
	dish := nutrition.Dish{}

	err := row.Scan(&dish.UID, &dish.Version, &dish.Name,
		&dish.Calories, &dish.Proteins, &dish.Fats, &dish.Carbs,
		&dish.CreatedAt, &dish.UpdatedAt)

	return dish, err
}

func (r *Repository) scanRows(rows *sql.Rows) (nutrition.Dish, error) {
	dish := nutrition.Dish{}

	err := rows.Scan(&dish.UID, &dish.Version, &dish.Name,
		&dish.Calories, &dish.Proteins, &dish.Fats, &dish.Carbs,
		&dish.CreatedAt, &dish.UpdatedAt)

	return dish, err
}

func (r *Repository) scanAllRows(rows *sql.Rows) ([]nutrition.Dish, error) {
	dishes := make([]nutrition.Dish, 0)

	for rows.Next() {
		dish, err := r.scanRows(rows)
		if err != nil {
			// TODO: instead of stopping and returning an error,
			// 		mark the item as errored and still return the list
			return nil, err
		}

		dishes = append(dishes, dish)
	}

	return dishes, nil
}
