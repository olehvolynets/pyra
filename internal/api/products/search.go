package products

import (
	"encoding/json"
	"net/http"

	"pyra/internal/api/handler"
	"pyra/pkg/nutrition"
)

func SearchProduct(api *API, w http.ResponseWriter, r *http.Request) {
	log := handler.RequestLogger(r)
	searchQuery := r.URL.Query().Get("q")

	products, err := api.ProductRepo.Search(r.Context(), searchQuery)
	if err != nil {
		log.Error("failed to fetch products", err)
		handler.InternalServerError(w)
		return
	}

	searchResults := buildSearchResults(products)

	res, err := json.Marshal(searchResults)
	if err != nil {
		log.Error("failed to marshal products", err)
		handler.InternalServerError(w)
		return
	}

	_, err = w.Write(res)
	if err != nil {
		log.Error("failed to write the response", err)
		handler.InternalServerError(w)
	}
}

type searchResult struct {
	UID     string `json:"uid"`
	Version uint16 `json:"version"`
	Label   string `json:"label"`

	Calories float32 `json:"calories"`
	Proteins float32 `json:"proteins"`
	Fats     float32 `json:"fats"`
	Carbs    float32 `json:"carbs"`
}

func buildSearchResults(products []nutrition.Product) []searchResult {
	searchResults := make([]searchResult, len(products))
	for idx, product := range products {
		searchResults[idx] = searchResult{
			UID:      string(product.UID),
			Version:  uint16(product.Version),
			Label:    string(product.Name),
			Calories: float32(product.Calories),
			Proteins: float32(product.Proteins),
			Fats:     float32(product.Fats),
			Carbs:    float32(product.Carbs),
		}
	}

	return searchResults
}
