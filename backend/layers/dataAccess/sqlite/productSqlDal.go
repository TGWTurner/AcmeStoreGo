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

func (ad ProductDatabase) Close() {
	db, err := ad.db.DB()

	if err != nil {
		panic("Failed to get product db instance")
	}

	db.Close()
}

func (pd ProductDatabase) GetAll() []utils.Product {
	var products []Product

	response := pd.db.Find(&products)

	if response.Error != nil {
		panic("Failed to get all products")
	}

	return ConvertFromDbProducts(products)
}

func (pd ProductDatabase) GetByIds(Ids ...int) []utils.Product {
	var products []Product

	response := pd.db.Find(&products, Ids)

	if response.Error != nil {
		Ids := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(Ids)), ","), "[]")

		panic("Failed to get products with id's: " + Ids)
	}

	return ConvertFromDbProducts(products)
}

func (pd ProductDatabase) GetCategories() []utils.ProductCategory {
	var categories []Category

	response := pd.db.Find(&categories)

	if response.Error != nil {
		panic("Failed tp get product categories")
	}

	return ConvertFromDbCategories(categories)
}

func (pd ProductDatabase) GetProductsByCategory(categoryId int) []utils.Product {
	var products []Product

	response := pd.db.Where("category_id = ?", categoryId).Find(&products)

	if response.Error != nil {
		panic("Unable to get product category with categoryId: " + strconv.Itoa(categoryId))
	}

	return ConvertFromDbProducts(products)
}

func (pd ProductDatabase) GetByText(searchTerm string) []utils.Product {
	var products []Product

	searchTerm = strings.TrimSpace(searchTerm)

	response := pd.db.
		Where("short_description LIKE ?", "%"+searchTerm+"%").
		Or("long_description LIKE ?", "%"+searchTerm+"%").
		Find(&products)

	if response.Error != nil {
		panic("Unable to get products for searchTerm: " + searchTerm)
	}

	return ConvertFromDbProducts(products)
}

func (pd ProductDatabase) GetWithCurrentDeals(date string) []utils.Product {
	var products []Product

	response := pd.db.Joins("INNER JOIN deals ON deals.product_id = products.id").
		Where("? >= deals.start_date", date).
		Where("? < deals.end_date", date).
		Find(&products)

	if response.Error != nil {
		panic("Unable to get products with current deals")
	}

	return ConvertFromDbProducts(products)
}

func (pd ProductDatabase) DecreaseStock(productQuantities []utils.OrderItem) {
	for _, item := range productQuantities {
		product := Product{
			Id: item.ProductId,
		}

		response := pd.db.Model(&product).Update("quantity_remaining", gorm.Expr("quantity_remaining - ?", item.Quantity))

		if response.Error != nil {
			panic("Failed to update quantity for productId: " + strconv.Itoa(item.ProductId))
		}
	}
}

type ProductDatabase struct {
	db *gorm.DB
}
