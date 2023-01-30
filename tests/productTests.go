package tests

import (
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

func TestGetProductGivenId() (bool, string, string) {
	db := SetUp()
	defer CloseDb(db)

	index := 1

	product, err := db.Product.GetByIds(index)

	if err != nil {
		return false, "testGetProductGivenId", "Failed to get product with id: " + strconv.Itoa(index)
	}

	expected := getTestProductById(index)

	if !reflect.DeepEqual(expected, product[0]) {
		return false, "testGetProductGivenId", "Actual Product did not match expected"
	}

	return true, "testGetProductGivenId", "Received correct product for index:" + strconv.Itoa(index)
}

func TestGetProductsGivenIds() (bool, string, string) {
	db := SetUp()
	defer CloseDb(db)

	indexes := []int{1, 2, 3}
	strIndexes := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(indexes)), ","), "[]")
	products, err := db.Product.GetByIds(indexes...)

	if err != nil {
		return false, "testGetProductsGivenIds", "Failed to get products by id's: " + strIndexes
	}

	for i, index := range indexes {
		expected := getTestProductById(index)
		if !reflect.DeepEqual(expected, products[i]) {
			return false, "testGetProductsGivenIds", "Actual Product did not match expected"
		}
	}

	return true, "testGetProductsGivenIds", "Got correct products for ids: " + strIndexes
}

func TestGetCategoriesReturnsCorrectCategories() (bool, string, string) {
	db := SetUp()
	defer CloseDb(db)

	categories, err := db.Product.GetCategories()

	if err != nil {
		return false, "testGetCategoriesReturnsCorrectCategories", "Failed to get categories"
	}

	expectedCategories := testData.GetProductTestData().Categories

	if len(categories) != len(expectedCategories) {
		return false, "testGetCategoriesReturnsCorrectCategories", "Received incorrect number of categories"
	}

	for _, category := range categories {
		if ok, err := assertCategoryHasExpectedName(category); !ok {
			return false, "testGetCategoriesReturnsCorrectCategories", err
		}
	}

	return true, "testGetCategoriesReturnsCorrectCategories", "Recieved expected categories"
}

func TestGetProductsByCategoryProvidesCorrectProducts() (bool, string, string) {
	db := SetUp()
	defer CloseDb(db)

	categoryId := 1

	products, err := db.Product.GetByCategory(categoryId)

	if err != nil {
		return false,
			"testGetProductsByCategoryProvidesCorrectProducts",
			"Failed too get categories"
	}

	expectedProducts := getTestProductsByCategory(categoryId)

	if len(expectedProducts) != len(products) {
		return true,
			"testGetProductsByCategoryProvidesCorrectProducts",
			"Returned different number of products than expected, expected: " +
				strconv.Itoa(len(expectedProducts)) +
				" actual: " +
				strconv.Itoa(len(products))
	}

	for _, product := range products {
		if !assertProductContainedWithinSliceOfProducts(product, expectedProducts) {
			return false,
				"testGetProductsByCategoryProvidesCorrectProducts",
				"Returned product not contained within expected products"
		}
	}

	return true,
		"testGetProductsByCategoryProvidesCorrectProducts",
		"Recieved correct products for category: " + strconv.Itoa(categoryId)
}

func TestGetProductsByText() (bool, string, string) {
	db := SetUp()
	defer CloseDb(db)

	searchText := "fruit"

	products, err := db.Product.GetByText(searchText)

	if err != nil {
		return false, "testGetProductsByText", "Failed to get products with search term: " + searchText
	}

	expectedProducts := getTestProductsByText(searchText)

	if len(products) != len(expectedProducts) {
		return false, "testGetProductsByText", "Returned wrong number of products when searched for term: " + searchText
	}

	for _, product := range products {
		if !assertProductContainedWithinSliceOfProducts(product, expectedProducts) {
			return false, "testGetProductsByText", "Product returned did not match test data"
		}
	}

	return true, "testGetProductByText", "Returned correct products when searched for term: " + searchText
}

func TestGetProductsWithCurrentDeals() (bool, string, string) {
	db := SetUp()
	defer CloseDb(db)

	currentDate := time.Now().Format("2006-01-02")

	products, err := db.Product.GetWithCurrentDeals(currentDate)

	if err != nil {
		return false,
			"testGetProductsWithCurrentDeals",
			"Failed to get products with current deals"
	}

	expectedProducts := getTestProductsWithCurrentDeals(currentDate)

	if len(expectedProducts) != len(products) {
		return false,
			"testGetProductsWithCurrentDeals",
			"Returned wrong number of products with deals expected: " +
				strconv.Itoa(len(expectedProducts)) +
				" actual: " +
				strconv.Itoa(len(products))
	}

	for _, product := range products {
		if !assertProductContainedWithinSliceOfProducts(product, expectedProducts) {
			return false,
				"testGetProductsWithCurrentDeals",
				"Returned Product was not contained within expected"
		}
	}

	return true,
		"testGetProductsWithCurrentDeals",
		"Returned products matched deals"
}

func TestDecreaseStockReducesStockByCorrectQuantity() (bool, string, string) {
	db := SetUp()
	defer CloseDb(db)

	productId := 1
	quantity := 5

	productQuantities := []utils.OrderItem{
		{
			ProductId: productId,
			Quantity:  quantity,
		},
	}

	db.Product.DecreaseStock(productQuantities)

	product, err := db.Product.GetByIds(productId)

	if err != nil {
		return false,
			"testDecreaseStockReducesStockByCorrectQuantity",
			"Failed to get products with id: " + strconv.Itoa(productId)
	}

	expectedProduct := getTestProductById(productId)
	expectedProduct.QuantityRemaining -= quantity

	if !reflect.DeepEqual(expectedProduct, product[0]) {
		return false,
			"testDecreaseStockReducesStockByCorrectQuantity",
			"Failed to decrease product quantity correctly expected: " +
				strconv.Itoa(expectedProduct.QuantityRemaining) +
				" actual: " +
				strconv.Itoa(product[0].QuantityRemaining)
	}

	return true,
		"testDecreaseStockReducesStockByCorrectQuantity",
		"Correctly decreased product stock by expected quantity"
}

func assertProductContainedWithinSliceOfProducts(product utils.Product, products []utils.Product) bool {
	for _, testProduct := range products {
		if reflect.DeepEqual(testProduct, product) {
			return true
		}
	}

	return false
}

func assertCategoryHasExpectedName(category utils.ProductCategory) (bool, string) {
	expectedCategories := testData.GetProductTestData().Categories

	for _, expectedCategory := range expectedCategories {
		if expectedCategory.Name == category.Name {
			return true, ""
		}
	}

	return false, "Category of name " + category.Name + " not expected"
}

func getTestProductById(id int) utils.Product {
	expectedProducts := testData.GetProductTestData().Products
	idx := slices.IndexFunc(
		expectedProducts,
		func(p utils.Product) bool { return p.Id == id },
	)

	return expectedProducts[idx]
}

func getTestProductsByText(text string) []utils.Product {
	expectedProducts := testData.GetProductTestData().Products
	resProducts := []utils.Product{}

	for _, product := range expectedProducts {
		if strings.Contains(product.ShortDescription, text) || strings.Contains(product.LongDescription, text) {
			resProducts = append(resProducts, product)
		}
	}

	return resProducts
}

func getTestProductsByCategory(categoryId int) []utils.Product {
	expectedProducts := testData.GetProductTestData().Products
	resProducts := []utils.Product{}

	for _, product := range expectedProducts {
		if categoryId == product.CategoryId {
			resProducts = append(resProducts, product)
		}
	}

	return resProducts
}

func getTestProductsWithCurrentDeals(currentDate string) []utils.Product {
	deals := testData.GetProductTestData().Deals
	products := []utils.Product{}

	for _, deal := range deals {
		if deal.StartDate < currentDate && deal.EndDate > currentDate {
			products = append(products, getTestProductById(deal.ProductId))
		}
	}

	return products
}
