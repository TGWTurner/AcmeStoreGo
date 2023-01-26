package memory

import (
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
	"strings"
)

func NewProductDatabase() ProductDatabase {
	testData := testData.GetProductTestData()
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

func (ad *ProductDatabase) Close() {}

func (pd *ProductDatabase) GetAll() []utils.Product {
	return pd.products
}

func (pd *ProductDatabase) GetByIds(Ids ...int) []utils.Product {
	products := []utils.Product{}

	for _, product := range pd.products {
		if contains(Ids, product.Id) {
			products = append(products, product)
		}
	}

	return products
}

func (pd *ProductDatabase) GetCategories() []utils.ProductCategory {
	return pd.categories
}

func (pd *ProductDatabase) GetByCategory(categoryId int) []utils.Product {
	products := []utils.Product{}

	for _, product := range pd.products {
		if product.CategoryId == categoryId {
			products = append(products, product)
		}
	}

	return products
}

func (pd *ProductDatabase) GetByText(searchTerm string) []utils.Product {
	products := []utils.Product{}

	for _, product := range pd.products {
		if productMatchesText(product, searchTerm) {
			products = append(products, product)
		}
	}

	return products
}

func (pd *ProductDatabase) GetWithCurrentDeals(date string) []utils.Product {
	products := []utils.Product{}

	for _, deal := range pd.deals {
		if deal.StartDate > date && deal.EndDate < date {
			products = append(products, pd.GetByIds(deal.ProductId)[0])
		}
	}

	return products
}

func (pd *ProductDatabase) DecreaseStock(productQuantities []utils.OrderItem) {
	for _, product := range productQuantities {
		pd.products[product.ProductId].QuantityRemaining = pd.products[product.ProductId].QuantityRemaining - product.Quantity
	}
}

type ProductDatabase struct {
	products   []utils.Product
	categories []utils.ProductCategory
	deals      []utils.ProductDeal
}
