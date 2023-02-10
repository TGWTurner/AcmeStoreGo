package memory

import (
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
	"errors"
	"strconv"
	"strings"
)

func NewProductDatabase() *ProductDatabaseImpl {
	testData := testData.GetProductTestData()
	products := testData.Products
	categories := testData.Categories
	deals := testData.Deals

	return &ProductDatabaseImpl{
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

func (ad *ProductDatabaseImpl) Close() {}

func (pd *ProductDatabaseImpl) GetAll() ([]utils.Product, error) {
	return pd.products, nil
}

func (pd *ProductDatabaseImpl) GetById(id int) (utils.Product, error) {
	for _, product := range pd.products {
		if id == product.Id {
			return product, nil
		}
	}

	return utils.Product{}, errors.New("Failed to get product with id: " + strconv.Itoa(id))
}

func (pd *ProductDatabaseImpl) GetByIds(Ids ...int) ([]utils.Product, error) {
	products := []utils.Product{}

	for _, product := range pd.products {
		if contains(Ids, product.Id) {
			products = append(products, product)
		}
	}

	return products, nil
}

func (pd *ProductDatabaseImpl) GetCategories() ([]utils.ProductCategory, error) {
	return pd.categories, nil
}

func (pd *ProductDatabaseImpl) GetByCategory(categoryId int) ([]utils.Product, error) {
	products := []utils.Product{}

	for _, product := range pd.products {
		if product.CategoryId == categoryId {
			products = append(products, product)
		}
	}

	return products, nil
}

func (pd *ProductDatabaseImpl) GetByText(searchTerm string) ([]utils.Product, error) {
	products := []utils.Product{}

	for _, product := range pd.products {
		if productMatchesText(product, searchTerm) {
			products = append(products, product)
		}
	}

	return products, nil
}

func (pd *ProductDatabaseImpl) GetWithCurrentDeals(date string) ([]utils.Product, error) {
	products := []utils.Product{}

	for _, deal := range pd.deals {
		if date >= deal.StartDate && date < deal.EndDate {
			product, err := pd.GetByIds(deal.ProductId)

			if err != nil {
				return []utils.Product{}, err
			}

			products = append(products, product[0])
		}
	}

	return products, nil
}

func (pd *ProductDatabaseImpl) DecreaseStock(productQuantities []utils.OrderItem) error {
	for _, productQuantity := range productQuantities {
		if product, err := pd.GetByIds(productQuantity.ProductId); len(product) == 0 || err != nil {
			return errors.New("Product with id: " + strconv.Itoa(productQuantity.ProductId) + " not found")
		}
	}

	for _, product := range productQuantities {
		for i := range pd.products {
			if pd.products[i].Id == product.ProductId {
				pd.products[i].QuantityRemaining -= product.Quantity
			}
		}
	}

	return nil
}

type ProductDatabaseImpl struct {
	products   []utils.Product
	categories []utils.ProductCategory
	deals      []utils.ProductDeal
}
