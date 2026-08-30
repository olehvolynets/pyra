package products

import (
	"fmt"
	"html/template"

	"pyra/pkg/nutrition"
)

const (
	ProductPATH = "/products/{uid}/{version}"
	ProductsPATH = "/products"
	NewProductPATH = ProductsPATH + "/new"
	EditProductPATH = ProductPATH + "/edit"

	ListProductsEP = "GET " + ProductsPATH
	ShowProductEP = "GET " + ProductPATH
	NewProductEP = "GET " + NewProductPATH
	CreateProductEP = "POST " + ProductsPATH
	DeleteProductEP = "DELETE " + ProductPATH
	EditProductEP = "GET " + EditProductPATH
	UpdateProductEP = "PUT " + ProductPATH
)

var URIHelpers = template.FuncMap{
	"productsURI":         ProductsURI,
	"productURI":          ProductURI,
	"productByRefURI":     ProductByRefURI,
	"newProductURI":       NewProductURI,
	"editProductURI":      EditProductURI,
	"editProductByRefURI": EditProductByRefURI,
}

func ProductsURI() string {
	return ProductsPATH
}

func ProductURI(product nutrition.Product) string {
	return ProductByRefURI(product.UID, product.Version)
}

func ProductByRefURI(uid nutrition.ProductUID, version nutrition.ProductVersion) string {
	return fmt.Sprintf("/products/%s/%d", uid, version)
}

func NewProductURI() string {
	return NewProductPATH
}

func EditProductURI(product nutrition.Product) string {
	return EditProductByRefURI(product.UID, product.Version)
}

func EditProductByRefURI(uid nutrition.ProductUID, version nutrition.ProductVersion) string {
	return ProductByRefURI(uid, version) + "/edit"
}
