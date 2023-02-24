package api

import (
	bl "backend/layers/businessLogic"
	"backend/utils"
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"
)

func NewProductApi(productService *bl.ProductService, s *sessions.CookieStore) *ProductApi {
	return &ProductApi{
		ps: *productService,
		s:  s,
	}
}

func (p *ProductApi) Close() {
	p.ps.Close()
}

// Search godoc
// @Summary Query or get all Products
// @ID Search
// @Produce json
// @Param search query string false "Text to search for"
// @Param caregory query int false "A Category Id to filter Products on"
// @Success 200 {object} []utils.Product "A JSON array of Products"
// @Router /api/product/catalogue [get]
func (p *ProductApi) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	query := map[string]string{}

	params := r.URL.Query()
	if params.Get("search") != "" {
		query["search"] = params.Get("search")
	} else if params.Get("category") != "" {
		query["category"] = params.Get("category")
	}

	products, err := p.ps.SearchProducts(query)

	if err != nil {
		Error(w, r, http.StatusInternalServerError, "error", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

// Categories godoc
// @Summary Get a list of product Categories
// @Description Signs up, deletes any existing session, creates a new one for this user. Will give an error if the user already exists.
// @ID Categories
// @Produce json
// @Success 200 {object} []utils.ProductCategory "A JSON array of Categories"
// @Router /api/product/categories [get]
func (p *ProductApi) Categories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	categories, err := p.ps.GetProductcategories()

	if err != nil {
		Error(w, r, http.StatusInternalServerError, "error", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

// Deals godoc
// @Summary Get deals that are valid for today
// @ID Deals
// @Produce json
// @Success 200 {object} []utils.Product "A JSON array of Products"
// @Router /api/product/deals [get]
func (p *ProductApi) Deals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	query := map[string]string{
		"dealDate": utils.GetFormattedDate(),
	}

	products, err := p.ps.SearchProducts(query)

	if err != nil {
		Error(w, r, http.StatusInternalServerError, "error", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

type ProductApi struct {
	ps bl.ProductService
	s  *sessions.CookieStore
}
