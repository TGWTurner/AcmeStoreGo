package memory

import (
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
	"strings"
)

func NewProductDatabase() ProductDatabase {
	testData := testData.ProductTestData{}.GetTestData()
	products := testData.Products
	categories := testData.Categories
	deals := testData.Deals

	return ProductDatabase{
		products:   products,
		categories: categories,
		deals:      deals,
	}
}

func contains(slice []int, provided int) bool {
	for _, value := range slice {
		if value == provided {
			return true
		}
	}

	return false
}

func productMatchesText(product utils.Product, searchTerm string) bool {
	searchTerm = strings.ToLower(strings.TrimSpace(searchTerm))

	if strings.Contains(product.ShortDescription, searchTerm) {
		return true
	}

	if strings.Contains(product.LongDescription, searchTerm) {
		return true
	}

	return false
}

func (pd *ProductDatabase) getAllProducts() []utils.Product {
	return pd.products
}

func (pd *ProductDatabase) getProductsByIds(Ids ...int) []utils.Product {
	products := []utils.Product{}

	for _, product := range pd.products {
		if contains(Ids, product.Id) {
			products = append(products, product)
		}
	}

	return products
}

func (pd *ProductDatabase) getProductCategories() []utils.Category {
	return pd.categories
}

func (pd *ProductDatabase) getProductsByCategory(categoryId int) []utils.Product {
	products := []utils.Product{}

	for _, product := range pd.products {
		if product.CategoryId == categoryId {
			products = append(products, product)
		}
	}

	return products
}

func (pd *ProductDatabase) getProductsByText(searchTerm string) []utils.Product {
	products := []utils.Product{}

	for _, product := range pd.products {
		if productMatchesText(product, searchTerm) {
			products = append(products, product)
		}
	}

	return products
}

func (pd *ProductDatabase) getProductsWithCurrentDeals(date string) []utils.Product {
	products := []utils.Product{}

	for _, deal := range pd.deals {
		if deal.StartDate > date && deal.EndDate < date {
			products = append(products, pd.getProductsByIds(deal.ProductId)[0])
		}
	}

	return products
}

func (pd *ProductDatabase) decreaseStock(productQuantities map[int]int) {
	for _, product := range pd.products {
		product.QuantityRemaining = product.QuantityRemaining - productQuantities[product.Id]
	}
}

type ProductDatabase struct {
	products   []utils.Product
	categories []utils.Category
	deals      []utils.Deal
}
