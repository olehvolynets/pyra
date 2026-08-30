package nutrition

import (
	"errors"
	"fmt"
)

var ErrNegative = errors.New("must be >= 0")

type Macro struct {
	Calories Measurement `fake:"{number:0,1000000}"`
	Proteins Measurement `fake:"{number:0,100000}"`
	Fats     Measurement `fake:"{number:0,100000}"`
	Carbs    Measurement `fake:"{number:0,100000}"`
}

func (m Macro) Add(other Macro) Macro {
	m.Calories += other.Calories
	m.Proteins += other.Proteins
	m.Fats += other.Fats
	m.Carbs += other.Carbs

	return m	
}

type MacroErrors struct {
	Calories error
	Proteins error
	Fats     error
	Carbs    error
}

func (e *MacroErrors) HasErrors() bool {
	caloriesErr := e.Calories != nil
	proteinsErr := e.Proteins != nil
	fatsErr := e.Fats != nil
	carbsErr := e.Carbs != nil

	return caloriesErr || proteinsErr || fatsErr || carbsErr
}

const macroErrFmt = `calories: %w
proteins: %w
fats: %w
carbs: %w`

func (e *MacroErrors) Error() string {
	return fmt.Errorf(macroErrFmt,
		e.Calories, e.Proteins, e.Fats, e.Carbs,
	).Error()
}

func (m Macro) Validate() MacroErrors {
	mErrors := MacroErrors{}

	if m.Calories < 0 {
		mErrors.Calories = ErrNegative
	}

	if m.Proteins < 0 {
		mErrors.Proteins = ErrNegative
	}

	if m.Fats < 0 {
		mErrors.Fats = ErrNegative
	}

	if m.Carbs < 0 {
		mErrors.Carbs = ErrNegative
	}

	return mErrors
}

// Normalize - scales macronutrient values from "portion" to 100.
// Records are stored in a normalized form N/100 g.
func (m *Macro) Normalize(portion float64) {
	factor := 100.0 / portion

	*m = m.Scale(factor)
}

func (m *Macro) Scale(factor float64) Macro {
	return Macro{
		Calories: m.Calories.Scale(factor),
		Proteins: m.Proteins.Scale(factor),
		Fats: m.Fats.Scale(factor),
		Carbs: m.Carbs.Scale(factor),
	}
}
