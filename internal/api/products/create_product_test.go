package products

import (
	"fmt"
	"net/http"
	"testing"

	"pyra/pkg/nutrition"
	"pyra/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateProductHandler(t *testing.T) {
	test.SetCWDToProjectRoot(t)

	db := test.DB(t)
	api := NewAPI(db)
	h := test.NewMux(http.MethodPost, ProductsPATH, api.Create(), t.Output())

	t.Run("success (minimal params)", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		res := h.Handle(http.MethodPost, ProductsPATH, nil, test.WithForm(map[string]string{
			"name": "asdf",
			"per":  "100 ",
		}))

		assert.Equal(t, http.StatusFound, res.StatusCode)

		dbProducts, err := api.ProductRepo.Index(t.Context())
		if assert.NoError(t, err) {
			assert.Len(t, dbProducts, 1, "expected 1 product")
			p := dbProducts[0]

			assert.Equal(t, p.Name, nutrition.ProductName("asdf"))
			assert.Equal(t, p.Calories, nutrition.Measurement(0))
			assert.Equal(t, p.Proteins, nutrition.Measurement(0))
			assert.Equal(t, p.Fats, nutrition.Measurement(0))
			assert.Equal(t, p.Carbs, nutrition.Measurement(0))
		}
	})

	t.Run("success (full)", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		res := h.Handle(http.MethodPost, ProductsPATH, nil, test.WithForm(map[string]string{
			"name":     "asdf",
			"per":      "100 ",
			"calories": "12",
			"proteins": "34",
			"fats":     "56",
			"carbs":    "78",
		}))

		assert.Equal(t, http.StatusFound, res.StatusCode)

		dbProducts, err := api.ProductRepo.Index(t.Context())
		if assert.NoError(t, err) {
			assert.Len(t, dbProducts, 1)
			p := dbProducts[0]

			assert.Equal(t, p.Name, nutrition.ProductName("asdf"))
			assert.Equal(t, p.Calories, nutrition.NewMeasurement(12))
			assert.Equal(t, p.Proteins, nutrition.NewMeasurement(34))
			assert.Equal(t, p.Fats, nutrition.NewMeasurement(56))
			assert.Equal(t, p.Carbs, nutrition.NewMeasurement(78))
		}
	})

	t.Run("success (normalise proportions)", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		res := h.Handle(http.MethodPost, ProductsPATH, nil, test.WithForm(map[string]string{
			"name":     "asdf",
			"per":      "50 ",
			"calories": "12",
			"proteins": "34",
			"fats":     "56",
			"carbs":    "78",
		}))

		assert.Equal(t, http.StatusFound, res.StatusCode)

		dbProducts, err := api.ProductRepo.Index(t.Context())
		if assert.NoError(t, err) {
			assert.Len(t, dbProducts, 1)
			p := dbProducts[0]

			assert.Equal(t, p.Name, nutrition.ProductName("asdf"))
			assert.Equal(t, p.Calories, nutrition.NewMeasurement(24))
			assert.Equal(t, p.Proteins, nutrition.NewMeasurement(68))
			assert.Equal(t, p.Fats, nutrition.NewMeasurement(112))
			assert.Equal(t, p.Carbs, nutrition.NewMeasurement(156))
		}
	})

	t.Run("blank form", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		res := h.Handle(http.MethodPost, ProductsPATH, nil)

		assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)

		dbProducts, err := api.ProductRepo.Index(t.Context())
		if assert.NoError(t, err) {
			assert.Len(t, dbProducts, 0)
		}
	})

	t.Run("when name already taken", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		existingProduct := nutrition.FakeProduct()
		err := api.ProductRepo.Create(t.Context(), &existingProduct)
		require.NoError(t, err)

		res := h.Handle(http.MethodPost, ProductsPATH, nil, test.WithForm(ParamsFromProduct(existingProduct)))

		assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)

		nProducts, err := api.ProductRepo.CountAll(t.Context())
		require.NoError(t, err)
		assert.Equal(t, nProducts, 1)
	})

	t.Run("when \"per\" is 0", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		product := nutrition.FakeProduct()
		formData := ParamsFromProduct(product)
		formData["per"] = "0"

		res := h.Handle(http.MethodPost, ProductsPATH, nil, test.WithForm(formData))

		assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)

		nProducts, err := api.ProductRepo.CountAll(t.Context())
		require.NoError(t, err)
		assert.Equal(t, nProducts, 0)
	})

	macroAttrs := []string{"calories", "proteins", "fats", "carbs"}
	for _, ma := range macroAttrs {
		t.Run(fmt.Sprintf("when %s is negative", ma), func(t *testing.T) {
			t.Cleanup(db.Truncate)

			product := nutrition.FakeProduct()
			formData := ParamsFromProduct(product)
			formData[ma] = "-1"

			res := h.Handle(http.MethodPost, ProductsPATH, nil, test.WithPathValues(formData))

			assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)

			nProducts, err := api.ProductRepo.CountAll(t.Context())
			require.NoError(t, err)
			assert.Equal(t, nProducts, 0)
		})
	}
}
