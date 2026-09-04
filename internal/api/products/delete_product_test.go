package products

import (
	"net/http"
	"pyra/pkg/nutrition"
	"pyra/test"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func Test_DeleteProductHandler(t *testing.T) {
	test.SetCWDToProjectRoot(t)

	db := test.DB(t)
	api := NewAPI(db)
	h := test.NewMux(http.MethodDelete, ProductPATH, api.Delete(), t.Output())

	t.Run("success", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		p := nutrition.FakeProduct()
		if err := api.ProductRepo.Create(t.Context(), &p); err != nil {
			t.Error(err)
			return
		}

		res := h.Handle(http.MethodDelete, ProductURI(p), nil)

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("when not found", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		p := nutrition.FakeProduct()
		res := h.Handle(http.MethodDelete, ProductURI(p), nil)

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("with invalid parameter (uid)", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		res := h.Handle(http.MethodDelete, ProductByRefURI("alsdkfj", 1), nil)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("with invalid parameter (version)", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		uid := nutrition.ProductUID(gofakeit.UUID())
		res := h.Handle(http.MethodDelete, ProductByRefURI(uid, 0), nil)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("when product is used in dish", func(t *testing.T) {
		t.Skip("pending")
	})
}
