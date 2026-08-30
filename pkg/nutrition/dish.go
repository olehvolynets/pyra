package nutrition

import (
	"fmt"
	"strings"
	"time"
)

type Dish struct {
	DishRecord

	Ingredients []Ingredient

	Errors DishErrors
}

type DishRecord struct {
	UID     DishUID     `fake:"{uuid}"`
	Version DishVersion `fake:"1"`

	Name DishName `fake:"{productname}"`

	Macro

	CreatedAt time.Time `fake:"-"`
	UpdatedAt time.Time `fake:"-"`
}

type DishErrors struct {
	Base error
	Name error

	Macro MacroErrors
}

func (e *DishErrors) HasErrors() bool {
	baseErr := e.Base != nil
	nameErr := e.Name != nil

	return baseErr || nameErr || e.Macro.HasErrors()
}

const dishErrFormat = `
Base: %w
Name: %w
%s
`

func (e *DishErrors) Error() string {
	return fmt.Errorf(dishErrFormat, e.Base, e.Name, e.Macro.Error()).Error()
}

type (
	DishUID     UID
	DishVersion Version
	DishName    string
)

func NewDishName(n string) (DishName, error) {
	n = strings.TrimSpace(n)
	if len(n) == 0 {
		return DishName(""), ErrBlank
	}

	return DishName(n), nil
}
