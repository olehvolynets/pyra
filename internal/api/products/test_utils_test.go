package products

import (
	"pyra/pkg/nutrition"
)

func ParamsFromProduct(p nutrition.Product) map[string]string {
	return map[string]string{
		"name":     string(p.Name),
		"per":      "100",
		"calories": p.Calories.String(),
		"proteins": p.Proteins.String(),
		"fats":     p.Fats.String(),
		"carbs":    p.Carbs.String(),
	}
}
