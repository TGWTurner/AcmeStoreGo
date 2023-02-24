package sqlite

import (
	"backend/layers/dataAccess/testData"
	"backend/utils"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func NewProductDatabase(db *gorm.DB) *ProductDatabaseImpl {
	pd := ProductDatabaseImpl{
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

	return &pd
}

func (ad ProductDatabaseImpl) Close() {
	db, err := ad.db.DB()

	if err != nil {
		panic("Failed to get product db instance")
	}

	db.Close()
}

func (pd ProductDatabaseImpl) GetAll() ([]utils.Product, error) {
	var products []Product

	response := pd.db.Find(&products)

	if response.Error != nil {
		return []utils.Product{}, errors.New("Failed to get all products:" + response.Error.Error())
	}

	return ConvertFromDbProducts(products), nil
}

func (pd ProductDatabaseImpl) GetById(id int) (utils.Product, error) {
	var product Product

	response := pd.db.Find(&product, id)

	if response.Error != nil {
		return utils.Product{}, errors.New("Failed to get product with id: " + strconv.Itoa(id) + " error: " + response.Error.Error())
	}

	return ConvertFromDbProduct(product), nil
}

func (pd ProductDatabaseImpl) GetByIds(ids ...int) ([]utils.Product, error) {
	var products []Product

	response := pd.db.Find(&products, ids)

	if response.Error != nil {
		Ids := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")

		return []utils.Product{}, errors.New("Failed to get products with id's: " + Ids + " error: " + response.Error.Error())
	}

	return ConvertFromDbProducts(products), nil
}

func (pd ProductDatabaseImpl) GetCategories() ([]utils.ProductCategory, error) {
	var categories []Category

	response := pd.db.Find(&categories)

	if response.Error != nil {
		return []utils.ProductCategory{}, errors.New("Failed to get product categories, error: " + response.Error.Error())
	}

	return ConvertFromDbCategories(categories), nil
}

func (pd ProductDatabaseImpl) GetByCategory(categoryId int) ([]utils.Product, error) {
	var products []Product

	response := pd.db.Where("category_id = ?", categoryId).Find(&products)

	if response.Error != nil {
		return []utils.Product{}, errors.New("Unable to get product category with categoryId: " + strconv.Itoa(categoryId) + ", error:" + response.Error.Error())
	}

	return ConvertFromDbProducts(products), nil
}

func (pd ProductDatabaseImpl) GetByText(searchTerm string) ([]utils.Product, error) {
	var products []Product

	searchTerm = strings.TrimSpace(searchTerm)

	response := pd.db.
		Where("short_description LIKE ?", "%"+searchTerm+"%").
		Or("long_description LIKE ?", "%"+searchTerm+"%").
		Find(&products)

	if response.Error != nil {
		return []utils.Product{}, errors.New("Unable to get products for searchTerm: " + searchTerm + ", error: " + response.Error.Error())
	}

	return ConvertFromDbProducts(products), nil
}

func (pd ProductDatabaseImpl) GetWithCurrentDeals(date string) ([]utils.Product, error) {
	var products []Product

	response := pd.db.Joins("INNER JOIN deals ON deals.product_id = products.id").
		Where("? >= deals.start_date", date).
		Where("? < deals.end_date", date).
		Find(&products)

	if response.Error != nil {
		return []utils.Product{}, errors.New("Unable to get products with current deals, error: " + response.Error.Error())
	}

	return ConvertFromDbProducts(products), nil
}

func (pd ProductDatabaseImpl) DecreaseStock(productQuantities []utils.OrderItem) error {
	for _, productQuantity := range productQuantities {
		if product, err := pd.GetByIds(productQuantity.ProductId); len(product) == 0 || err != nil {
			return errors.New("Failed to get product for id: " + strconv.Itoa(productQuantity.ProductId) + " to update")
		}
	}

	for _, item := range productQuantities {
		product := Product{
			Id: item.ProductId,
		}

		response := pd.db.Model(&product).Update("quantity_remaining", gorm.Expr("quantity_remaining - ?", item.Quantity))

		if response.Error != nil {
			return errors.New("Failed to update quantity for productId: " + strconv.Itoa(item.ProductId) + ", error: " + response.Error.Error())
		}
	}

	return nil
}

type ProductDatabaseImpl struct {
	db *gorm.DB
}
