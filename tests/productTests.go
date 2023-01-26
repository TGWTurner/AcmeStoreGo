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

func TestGetProductGivenId() {
	db := SetUp()
	defer CloseDb(db)

	index := 1

	product := db.Product.GetByIds(index)[0]
	expected := getTestProductById(index)

	if !reflect.DeepEqual(expected, product) {
		PrintTestResult(false, "testGetProductGivenId", "Actual Product did not match expected")
		return
	}

	PrintTestResult(true, "testGetProductGivenId", "Received correct product for index:"+strconv.Itoa(index))
}

func TestGetProductsGivenIds() {
	db := SetUp()
	defer CloseDb(db)

	indexes := []int{1, 2, 3}

	products := db.Product.GetByIds(indexes...)

	for i, index := range indexes {
		expected := getTestProductById(index)
		if !reflect.DeepEqual(expected, products[i]) {
			PrintTestResult(false, "testGetProductsGivenIds", "Actual Product did not match expected")
			return
		}
	}

	strIndexes := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(indexes)), ","), "[]")
	PrintTestResult(true, "testGetProductsGivenIds", "Got correct products for ids: "+strIndexes)
}

func TestGetCategoriesReturnsCorrectCategories() {
	db := SetUp()
	defer CloseDb(db)

	categories := db.Product.GetCategories()
	expectedCategories := testData.GetProductTestData().Categories

	if len(categories) != len(expectedCategories) {
		PrintTestResult(false, "testGetCategoriesReturnsCorrectCategories", "Received incorrect number of categories")
		return
	}

	for _, category := range categories {
		if ok, err := assertCategoryHasExpectedName(category); !ok {
			PrintTestResult(false, "testGetCategoriesReturnsCorrectCategories", err)
			fmt.Println("TEST FAILED -- " + err)
			return
		}
	}

	PrintTestResult(true, "testGetCategoriesReturnsCorrectCategories", "Recieved expected categories")
}

func TestGetProductsByCategoryProvidesCorrectProducts() {
	db := SetUp()
	defer CloseDb(db)

	categoryId := 1

	products := db.Product.GetByCategory(categoryId)
	expectedProducts := getTestProductsByCategory(categoryId)

	if len(expectedProducts) != len(products) {
		PrintTestResult(
			true,
			"testGetProductsByCategoryProvidesCorrectProducts",
			"Returned different number of products than expected, expected: "+
				strconv.Itoa(len(expectedProducts))+
				" actual: "+
				strconv.Itoa(len(products)),
		)
		return
	}

	for _, product := range products {
		if !assertProductContainedWithinSliceOfProducts(product, expectedProducts) {
			PrintTestResult(
				false,
				"testGetProductsByCategoryProvidesCorrectProducts",
				"Returned product not contained within expected products",
			)
			return
		}
	}

	PrintTestResult(
		true,
		"testGetProductsByCategoryProvidesCorrectProducts",
		"Recieved correct products for category: "+strconv.Itoa(categoryId),
	)
}

func TestGetProductsByText() {
	db := SetUp()
	defer CloseDb(db)

	searchText := "fruit"

	products := db.Product.GetByText(searchText)
	expectedProducts := getTestProductsByText(searchText)

	if len(products) != len(expectedProducts) {
		PrintTestResult(false, "testGetProductsByText", "Returned wrong number of products when searched for term: "+searchText)
		return
	}

	for _, product := range products {
		if !assertProductContainedWithinSliceOfProducts(product, expectedProducts) {
			PrintTestResult(false, "testGetProductsByText", "Product returned did not match test data")
		}
	}

	PrintTestResult(true, "testGetProductByText", "Returned correct products when searched for term: "+searchText)
}

func TestGetProductsWithCurrentDeals() {
	db := SetUp()
	defer CloseDb(db)

	currentDate := time.Now().Format("2006-01-02")

	products := db.Product.GetWithCurrentDeals(currentDate)
	expectedProducts := getTestProductsWithCurrentDeals(currentDate)

	if len(expectedProducts) != len(products) {
		PrintTestResult(
			false,
			"testGetProductsWithCurrentDeals",
			"Returned wrong number of products with deals expected: "+
				strconv.Itoa(len(expectedProducts))+
				" actual: "+
				strconv.Itoa(len(products)),
		)
		return
	}

	for _, product := range products {
		if !assertProductContainedWithinSliceOfProducts(product, expectedProducts) {
			PrintTestResult(
				false,
				"testGetProductsWithCurrentDeals",
				"Returned Product was not contained within expected",
			)
			return
		}
	}

	PrintTestResult(
		true,
		"testGetProductsWithCurrentDeals",
		"Returned products matched deals",
	)
}

func TestDecreaseStockReducesStockByCorrectQuantity() {
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

	product := db.Product.GetByIds(productId)[0]
	expectedProduct := getTestProductById(productId)
	expectedProduct.QuantityRemaining -= quantity

	if !reflect.DeepEqual(expectedProduct, product) {
		PrintTestResult(
			false,
			"testDecreaseStockReducesStockByCorrectQuantity",
			"Failed to decrease product quantity correctly expected: "+
				strconv.Itoa(expectedProduct.QuantityRemaining)+
				" actual: "+
				strconv.Itoa(product.QuantityRemaining),
		)
	}

	PrintTestResult(
		true,
		"testDecreaseStockReducesStockByCorrectQuantity",
		"Correctly decreased product stock by expected quantity",
	)
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
