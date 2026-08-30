package dishes

import (
	"html/template"

	"pyra/internal/api/base"
	"pyra/pkg/db"
)

func NewTestDishAPI(db db.DBTX, fms ...template.FuncMap) *API {
	drivers := base.Drivers(URIHelpers)
	baseAPI := base.NewAPI(db, drivers)

	return NewAPI(baseAPI)
}
