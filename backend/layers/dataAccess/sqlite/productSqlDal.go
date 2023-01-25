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
	products := ConvertToDbProducts(testData.Products)
	categories := ConvertToDbCategories(testData.Categories)
	deals := ConvertToDbDeals(testData.Deals)

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

func (pd ProductDatabase) GetAll() []utils.Product {
	var accounts []utils.Product

	response := pd.db.Find(&accounts)

	if response.Error != nil {
		panic("Failed to get products")
	}

	return accounts
}

func (pd ProductDatabase) GetByIds(Ids ...int) []utils.Product {
	var accounts []utils.Product

	response := pd.db.Find(&accounts, Ids)

	if response.Error != nil {
		Ids := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(Ids)), ","), "[]")

		panic("Failed to get products with id's: " + Ids)
	}

	return accounts
}

func (pd ProductDatabase) GetCategories() []utils.ProductCategory {
	var categories []utils.ProductCategory

	response := pd.db.Find(&categories)

	if response != nil {
		panic("Failed tp get product categories")
	}

	return categories
}

func (pd ProductDatabase) GetProductsByCategory(categoryId int) utils.ProductCategory {
	var category utils.ProductCategory

	response := pd.db.First(category, categoryId)

	if response.Error != nil {
		panic("Unable to get product category with categoryId: " + strconv.Itoa(categoryId))
	}

	return category
}

func (pd ProductDatabase) GetByText(searchTerm string) []utils.Product {
	var products []utils.Product
	response := pd.db.
		Where("short_description LIKE ?", searchTerm).
		Or("long_description LIKE ?", searchTerm).
		Find(&products)

	if response.Error != nil {
		panic("Unable to get products for searchTerm: " + searchTerm)
	}

	return products
}

func (pd ProductDatabase) GetCurrentDeals(date string) []utils.Product {
	var products []utils.Product

	response := pd.db.Find(&products).
		Joins("INNER JOIN deals ON deals.product_id = products.id").
		Where("? >= deals.start_date and ? < deals.end_date", date, date)

	if response != nil {
		panic("Unable to get products with current deals")
	}

	return products
}

func (pd ProductDatabase) DecreaseStock(productQuantities map[int]int) {
	for productId, quantity := range productQuantities {
		product := utils.Product{
			Id: productId,
		}

		response := pd.db.Model(&product).Update("quantity_remaining", gorm.Expr("quantity - ?", quantity))

		if response.Error != nil {
			panic("Failed to update quantity for productId: " + strconv.Itoa(productId))
		}
	}
}

type ProductDatabase struct {
	db *gorm.DB
}
