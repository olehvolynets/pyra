package products

import (
	"fmt"
	"net/http"
	"testing"

	"pyra/pkg/nutrition"
	"pyra/test"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UpdateProductHandler(t *testing.T) {
	test.SetCWDToProjectRoot(t)

	db := test.DB(t)
	api := NewAPI(db)
	h := test.NewMux(http.MethodPut, ProductPATH, api.Update(), t.Output())

	t.Run("when not found", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		updatedProduct := nutrition.FakeProduct()

		res := h.Handle(http.MethodPut, ProductURI(updatedProduct), nil,
			test.WithForm(ParamsFromProduct(nutrition.FakeProduct())))

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("with invalid parameter (uid)", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		res := h.Handle(http.MethodPut, ProductByRefURI("alsdkfj", 1), nil)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("with invalid parameter (version)", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		uid := nutrition.ProductUID(gofakeit.UUID())
		res := h.Handle(http.MethodPut, ProductByRefURI(uid, 0), nil)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("when not used in dish - updates the record in place", func(t *testing.T) {
		t.Cleanup(db.Truncate)

		product := nutrition.FakeProduct()
		err := api.ProductRepo.Create(t.Context(), &product)
		require.NoError(t, err)

		paramsProduct := nutrition.FakeProduct()
		formData := ParamsFromProduct(paramsProduct)
		formData["per"] = "50"
		res := h.Handle(http.MethodPut, ProductURI(product), nil, test.WithForm(formData))

		assert.Equal(t, http.StatusOK, res.StatusCode)

		updatedProduct, err := api.ProductRepo.FindByRef(t.Context(), product.ProductRef)
		require.NoError(t, err)

		assert.Equal(t, paramsProduct.Name, updatedProduct.Name)
		assert.Equal(t, paramsProduct.Calories.Scale(2), updatedProduct.Calories)
		assert.Equal(t, paramsProduct.Proteins.Scale(2), updatedProduct.Proteins)
		assert.Equal(t, paramsProduct.Fats.Scale(2), updatedProduct.Fats)
		assert.Equal(t, paramsProduct.Carbs.Scale(2), updatedProduct.Carbs)
		if t.Failed() {
			fmt.Println(formData)
		}
	})

	t.Run("when used in dish (dish NOT used) - updates the record in place", func(t *testing.T) {
		t.Skip("pending")
		t.Cleanup(db.Truncate)
	})

	t.Run("when used in dish (dish used) - creates new version", func(t *testing.T) {
		t.Skip("pending")
		t.Cleanup(db.Truncate)
	})
}
