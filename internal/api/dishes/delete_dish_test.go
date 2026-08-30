package dishes

import (
	"net/http"
	"testing"

	"pyra/pkg/nutrition"
	"pyra/test"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func Test_DeleteDishHandler(t *testing.T) {
	test.SetCWDToProjectRoot(t)

	db := test.DB(t)
	ep := NewTestDishAPI(db).Delete().(*DeleteDishHandler)
	h := test.NewMux(http.MethodDelete, DishPATH, ep, t.Output())

	t.Run("success", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		d := nutrition.FakeDish()
		if err := ep.DishRepo.Create(t.Context(), &d); err != nil {
			t.Fatal(err)
		}

		res := h.Handle(http.MethodDelete, DishURI(d), nil)

		assert.Equal(t, http.StatusOK, res.StatusCode)

		dbDishes, err := ep.DishRepo.Index(t.Context())
		if err != nil {
			t.Fatal(err)
		}

		if len(dbDishes) != 0 {
			t.Fatalf("expected to have 0 dishes after deleteion, got: %d", len(dbDishes))
		}

		// TODO: test that associated ingredients are deleted as well.
	})

	t.Run("when not found", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		d := nutrition.FakeDish()
		res := h.Handle(http.MethodDelete, DishURI(d), nil)

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("with invalid parameter (uid)", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		res := h.Handle(http.MethodDelete, DishByRefURI("alsdkfj", 1), nil)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("with invalid parameter (version)", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		uid := nutrition.DishUID(gofakeit.UUID())
		res := h.Handle(http.MethodDelete, DishByRefURI(uid, 0), nil)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("when dish was used in past menu", func(t *testing.T) {
		t.Skip("pending")
	})
}
