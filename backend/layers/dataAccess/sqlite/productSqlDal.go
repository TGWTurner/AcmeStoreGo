package sqlite

import (
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func NewProductDatabase(db *gorm.DB) ProductDatabase {
	pd := ProductDatabase{
		db: db,
	}

	testData := testData.GetProductTestData()
	products := testData.Products
	categories := testData.Categories
	deals := testData.Deals

	if res := db.Create(products); res.Error != nil {
		panic("Failed to create test products")
	}

	if res := db.Create(categories); res.Error != nil {
		panic("Failed to create test products")
	}

	if res := db.Create(deals); res.Error != nil {
		panic("Failed to create test products")
	}

	return pd
}

func (pd ProductDatabase) getAllProducts() []utils.Product {
	var accounts []utils.Product

	response := pd.db.Find(&accounts)

	if response.Error != nil {
		panic("Failed to get products")
	}

	return accounts
}

func (pd ProductDatabase) getProductsByIds(Ids ...int) []utils.Product {
	var accounts []utils.Product

	response := pd.db.Find(&accounts, Ids)

	if response.Error != nil {
		Ids := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(Ids)), ","), "[]")

		panic("Failed to get products with id's: " + Ids)
	}

	return accounts
}

func (pd ProductDatabase) getProductCategories() []utils.ProductCategory {
	var categories []utils.ProductCategory

	response := pd.db.Find(&categories)

	if response != nil {
		panic("Failed tp get product categories")
	}

	return categories
}

func (pd ProductDatabase) getProductsByCategory(categoryId int) utils.ProductCategory {
	var category utils.ProductCategory

	response := pd.db.First(category, categoryId)

	if response.Error != nil {
		panic("Unable to get product category with categoryId: " + strconv.Itoa(categoryId))
	}

	return category
}

func (pd ProductDatabase) getProductsByText(searchTerm string) []utils.Product {
	var products []utils.Product
	response := pd.db.
		Where("shortDescription LIKE ?", searchTerm).
		Or("longDescription LIKE ?", searchTerm).
		Find(&products)

	if response.Error != nil {
		panic("Unable to get products for searchTerm: " + searchTerm)
	}

	return products
}

func (pd ProductDatabase) getProductsWithCurrentDeals(date string) []utils.Product {
	var products []utils.Product

	response := pd.db.Find(&products).
		Joins("INNER JOIN deals ON deals.productId = products.ProductId").
		Where("? >= deals.startDate and ? < deals.EndDate", date, date)

	if response != nil {
		panic("Unable to get products with current deals")
	}

	return products
}

func (pd ProductDatabase) decreaseStock(productQuantities map[int]int) {
	for productId, quantity := range productQuantities {
		product := utils.Product{
			Id: productId,
		}

		response := pd.db.Model(&product).Update("quantityRemaining", gorm.Expr("quantity - ?", quantity))

		if response.Error != nil {
			panic("Failed to update quantity for productId: " + strconv.Itoa(productId))
		}
	}
}

type ProductDatabase struct {
	db *gorm.DB
}
