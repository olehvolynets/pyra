package products

import (
	"encoding/json"
	"net/http"

	"pyra/internal/api/base"
	"pyra/pkg/nutrition"
)

type SearchProductsHandler struct {
	*base.Handler
	productRepo nutrition.ProductRepository
}

func (h *SearchProductsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := h.RequestLogger(r)
	searchQuery := r.URL.Query().Get("q")

	products, err := h.productRepo.Search(r.Context(), searchQuery)
	if err != nil {
		log.Error("failed to fetch products", err)
		h.InternalServerError(w)
		return
	}

	searchResults := h.buildSearchResults(products)

	res, err := json.Marshal(searchResults)
	if err != nil {
		log.Error("failed to marshal products", err)
		h.InternalServerError(w)
		return
	}

	_, err = w.Write(res)
	if err != nil {
		log.Error("failed to write the response", err)
		h.InternalServerError(w)
	}
}

func (h *SearchProductsHandler) buildSearchResults(products []nutrition.Product) []searchResult {
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

type searchResult struct {
	UID     string `json:"uid"`
	Version uint16 `json:"version"`
	Label   string `json:"label"`

	Calories float32 `json:"calories"`
	Proteins float32 `json:"proteins"`
	Fats     float32 `json:"fats"`
	Carbs    float32 `json:"carbs"`
}
