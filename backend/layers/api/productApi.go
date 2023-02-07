package api

import (
	bl "bjssStoreGo/backend/layers/businessLogic"
	"bjssStoreGo/backend/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

func NewProductApi(productService *bl.ProductService) *ProductApi {
	return &ProductApi{
		ps: *productService,
	}
}

func (p ProductApi) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := map[string]string{}

	params := r.URL.Query()
	if params.Get("search") != "" {
		query["search"] = params.Get("search")
	} else if params.Get("category") != "" {
		query["category"] = params.Get("category")
	}

	fmt.Println(query)

	products, err := p.ps.SearchProducts(query)

	if err != nil {
		//log the error?
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (p ProductApi) Categories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	categories, err := p.ps.GetProductcategories()

	if err != nil {
		//log the error?
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

func (p ProductApi) Deals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := map[string]string{
		"dealDate": utils.GetFormattedDate(),
	}

	products, err := p.ps.SearchProducts(query)

	if err != nil {
		//log the error?
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

type ProductApi struct {
	ps bl.ProductService
}
