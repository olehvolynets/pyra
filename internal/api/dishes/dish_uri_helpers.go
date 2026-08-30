package dishes

import (
	"fmt"
	"html/template"

	"pyra/pkg/nutrition"
)

const (
	DishPATH     = "/dishes/{uid}/{version}"
	DishesPATH   = "/dishes"
	NewDishPATH  = DishesPATH + "/new"
	EditDishPATH = DishPATH + "/edit"

	ListDishsEP  = "GET " + DishesPATH
	ShowDishEP   = "GET " + DishPATH
	NewDishEP    = "GET " + NewDishPATH
	CreateDishEP = "POST " + DishesPATH
	DeleteDishEP = "DELETE " + DishPATH
	EditDishEP   = "GET " + EditDishPATH
	UpdateDishEP = "PUT " + DishPATH
)

var URIHelpers = template.FuncMap{
	"dishesURI": DishesURI,
	"dishURI": DishURI,
	"dishByRefURI": DishByRefURI,
	"newDishURI": NewDishURI,
	"editDishURI": EditDishURI,
	"editDishByRefURI": EditDishByRefURI,
}

func DishesURI() string {
	return DishesPATH
}

func DishURI(dish nutrition.Dish) string {
	return DishByRefURI(dish.UID, dish.Version)
}

func DishByRefURI(uid nutrition.DishUID, version nutrition.DishVersion) string {
	return fmt.Sprintf("/dishes/%s/%d", uid, version)
}

func NewDishURI() string {
	return NewDishPATH
}

func EditDishURI(dish nutrition.Dish) string {
	return EditDishByRefURI(dish.UID, dish.Version)
}

func EditDishByRefURI(uid nutrition.DishUID, version nutrition.DishVersion) string {
	return DishByRefURI(uid, version) + "/edit"
}
