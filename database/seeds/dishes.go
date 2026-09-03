package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"uuid"

	"pyra/pkg/db"
	"pyra/pkg/nutrition"
)

func seedDishes(conn db.DBTX) error {
	l.Trace("BEGIN seeding dishes")
	defer l.Trace("END seeding dishes")

	dishes := make([]nutrition.Dish, 0, 10)

	for i := 0; i < 10; i += 1 {
		uid := uuid.New()

		dish := nutrition.Dish{
			Name:    nutrition.DishName(fmt.Sprintf("Dish #%d", i)),
			UID:     nutrition.DishUID(uid.String()),
			Version: 1,
			Macro: nutrition.Macro{
				Calories: nutrition.Measurement(rand.Int32N(400)),
				Proteins: nutrition.Measurement(rand.Int32N(100)),
				Carbs:    nutrition.Measurement(rand.Int32N(150)),
				Fats:     nutrition.Measurement(rand.Int32N(50)),
			},
		}

		dishes = append(dishes, dish)
	}

	ctx := context.Background()
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `INSERT INTO dishes (
		name, uid, version, calories, proteins, fats, carbs
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7
	);`

	for _, dish := range dishes {
		_, err := tx.ExecContext(ctx, query,
			dish.Name, dish.UID, dish.Version,
			dish.Calories, dish.Carbs, dish.Fats, dish.Carbs)

		if err != nil {
			l.Debug(err.Error())
			break
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
