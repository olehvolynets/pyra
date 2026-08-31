package products

import (
	"errors"
	"fmt"
	"strings"

	"pyra/pkg/nutrition"
)

var (
	ErrNotPositive   = errors.New("must be greater than 0")
	ErrInvalidNumber = errors.New("must be a number")
)

type ProductForm struct {
	Product nutrition.Product

	Name string

	Calories string
	Proteins string
	Fats     string
	Carbs    string
	Per      string

	Errors productFormErrors
}

func NewProductForm(fetch func(key string) string) ProductForm {
	form := ProductForm{}

	fetchValue := func(key string) string {
		return strings.TrimSpace(fetch(key))
	}

	form.Name = fetchValue("name")
	form.Calories = fetchValue("calories")
	form.Proteins = fetchValue("proteins")
	form.Fats = fetchValue("fats")
	form.Carbs = fetchValue("carbs")
	form.Per = fetchValue("per")

	return form
}

func FormFromProduct(p nutrition.Product) ProductForm {
	form := ProductForm{Per: "100", Product: p}

	form.Name = string(p.Name)
	form.Calories = p.Calories.String()
	form.Proteins = p.Proteins.String()
	form.Fats = p.Fats.String()
	form.Carbs = p.Carbs.String()

	return form
}

func (f *ProductForm) BuildProduct() (product nutrition.Product) {
	calories, err := nutrition.ParseMeasurement(f.Calories)
	product.Calories = calories
	if err != nil {
		f.Errors.Calories = ErrInvalidNumber
	}

	proteins, err := nutrition.ParseMeasurement(f.Proteins)
	product.Proteins = proteins
	if err != nil {
		f.Errors.Proteins = ErrInvalidNumber
	}

	fats, err := nutrition.ParseMeasurement(f.Fats)
	product.Fats = fats
	if err != nil {
		f.Errors.Fats = ErrInvalidNumber
	}

	carbs, err := nutrition.ParseMeasurement(f.Carbs)
	product.Carbs = carbs
	if err != nil {
		f.Errors.Carbs = ErrInvalidNumber
	}


	per, err := nutrition.ParseMeasurement(f.Per)
	if err != nil {
		f.Errors.Per = ErrInvalidNumber
	}

	if f.HasErrors() {
		return
	}

	product.Name = nutrition.ProductName(f.Name)

	// Validations start here, error handling before was only for parsing.
	product.Errors.Name = product.Name.Validate()
	product.Errors.MacroErrors = product.Macro.Validate()

	if per <= 0 {
		f.Errors.Per = ErrNotPositive
	}

	f.SetProductErrors(product.Errors)

	if f.HasErrors() {
		return
	}

	product.Normalize(per.Float())

	return
}

func (f *ProductForm) SetProductErrors(e nutrition.ProductErrors) {
	f.Errors.ProductErrors = e
}

type productFormErrors struct {
	nutrition.ProductErrors
	Per error
}

func (e *productFormErrors) HasErrors() bool {
	perErr := e.Per != nil

	return perErr || e.ProductErrors.HasErrors()
}

func (e *productFormErrors) Error() string {
	return fmt.Errorf("per: %w\n%s", e.Per, e.ProductErrors.Error()).Error()
}

func (f *ProductForm) HasErrors() bool {
	return f.Errors.HasErrors()
}
