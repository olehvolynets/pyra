package products

import (
	"html/template"

	"pyra/internal/api/base"
	"pyra/pkg/db"
	"pyra/pkg/nutrition"
)

func NewTestProductsAPI(db db.DBTX, fms ...template.FuncMap) *API {
	drivers := base.NewTemplate(URIHelpers)
	baseAPI := base.NewAPI(db, drivers)

	return NewAPI(baseAPI)
}

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
