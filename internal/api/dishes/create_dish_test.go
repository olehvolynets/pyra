package dishes

import (
	"net/http"
	"testing"

	"pyra/pkg/nutrition"
	"pyra/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateDishHandler(t *testing.T) {
	test.SetCWDToProjectRoot(t)

	db := test.DB(t)
	api := NewAPI(db)
	h := test.NewMux(http.MethodPost, DishesPATH, api.Create(), t.Output())

	t.Run("success (minimal params)", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		res := h.Handle(http.MethodPost, DishesPATH, nil, test.WithForm(map[string]string{
			"name": "asdf",
		}))

		assert.Equal(t, http.StatusFound, res.StatusCode)

		dbDishes, err := api.DishRepo.Index(t.Context())
		if assert.NoError(t, err) {
			require.Len(t, dbDishes, 1)

			d := dbDishes[0]
			assert.Equal(t, d.Name, nutrition.DishName("asdf"))
			assert.Equal(t, d.Calories, nutrition.Measurement(0))
			assert.Equal(t, d.Proteins, nutrition.Measurement(0))
			assert.Equal(t, d.Fats, nutrition.Measurement(0))
			assert.Equal(t, d.Carbs, nutrition.Measurement(0))
		}
	})

	t.Run("success (all params)", func(t *testing.T) {
		t.Skip("pending")
		t.Cleanup(db.Truncate)
	})

	t.Run("blank form", func(t *testing.T) {
		t.Skip("pending")
		t.Cleanup(db.Truncate)
	})

	t.Run("name taken", func(t *testing.T) {
		t.Skip("pending")
		t.Cleanup(db.Truncate)
	})

	// TODO: ingredients validations
}
