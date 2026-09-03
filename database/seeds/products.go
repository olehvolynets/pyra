package main

import (
	"context"
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"uuid"

	"pyra/pkg/db"
	"pyra/pkg/nutrition"
)

func seedProducts(conn db.DBTX) error {
	l.Trace("BEGIN seeding products")
	defer l.Trace("END seeding products")

	f, err := os.Open("./database/seeds/products_sample.csv")
	if err != nil {
		panic(err)
	}

	r := csv.NewReader(f)
	r.Comma = ';'

	products := make([]nutrition.Product, 0)

	for {
		csvRecord, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			panic(err)
		}

		uid := uuid.New()

		calories, err := strconv.ParseFloat(csvRecord[2], 32)
		if err != nil {
			l.Debug(err.Error())
			continue
		}
		proteins, err := strconv.ParseFloat(csvRecord[3], 32)
		if err != nil {
			l.Debug(err.Error())
			continue
		}
		fats, err := strconv.ParseFloat(csvRecord[4], 32)
		if err != nil {
			l.Debug(err.Error())
			continue
		}
		carbs, err := strconv.ParseFloat(csvRecord[5], 32)
		if err != nil {
			l.Debug(err.Error())
			continue
		}

		record := nutrition.Product{
			ProductRecord: nutrition.ProductRecord{
				Name:    nutrition.ProductName(csvRecord[0]),
				UID:     nutrition.ProductUID(uid.String()),
				Version: 1,
				Macro: nutrition.Macro{
					Calories: nutrition.NewMeasurement(calories),
					Proteins: nutrition.NewMeasurement(proteins),
					Fats:     nutrition.NewMeasurement(fats),
					Carbs:    nutrition.NewMeasurement(carbs),
				},
			},
		}

		products = append(products, record)
	}

	ctx := context.Background()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `INSERT INTO products (
			name, uid, version, calories, proteins, fats, carbs
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		);`

	for _, product := range products {
		_, err := tx.ExecContext(ctx, query,
			product.Name, product.UID, product.Version,
			product.Calories, product.Proteins, product.Fats, product.Carbs)

		if err != nil {
			l.Debug(err.Error())
			continue
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
