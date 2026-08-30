package nutrition

import "github.com/brianvoe/gofakeit/v7"

func fake(v any) {
	if err := gofakeit.Struct(v); err != nil {
		panic(err)
	}
}

func FakeProduct() Product {
	p := Product{}

	fake(&p)

	return p
}

func FakeDish() Dish {
	d := Dish{}

	fake(&d)

	return d
}

func FakeIngredient() Ingredient {
	i := Ingredient{}

	fake(&i)

	return i
}
