package products

import (
	"net/http"
	"pyra/pkg/nutrition"
	"pyra/test"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func Test_EditProductHandler(t *testing.T) {
	test.SetCWDToProjectRoot(t)

	db := test.DB(t)
	api := NewAPI(db)
	h := test.NewMux(http.MethodGet, EditProductPATH, api.Edit(), t.Output())

	t.Run("success", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		p := nutrition.FakeProduct()
		if err := api.ProductRepo.Create(t.Context(), &p); err != nil {
			t.Error(err)
			return
		}

		res := h.Handle(http.MethodGet, EditProductURI(p), nil)

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("when not found", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		uid := nutrition.ProductUID(gofakeit.UUID())
		res := h.Handle(http.MethodGet, EditProductByRefURI(uid, 1), nil)

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("when invalid parameter (uid)", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		res := h.Handle(http.MethodGet, EditProductByRefURI("asdfkfj", 1), nil)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)	
	})

	t.Run("when invalid parameter (version)", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		uid := nutrition.ProductUID(gofakeit.UUID())
		res := h.Handle(http.MethodGet, EditProductByRefURI(uid, 0), nil)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)	
	})
}
