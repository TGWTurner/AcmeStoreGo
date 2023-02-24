package dbTests

import (
	"backend/layers/dataAccess/testData"
	"backend/test"
	"backend/utils"
	"reflect"
	"testing"
)

func TestGetProductGivenId(t *testing.T) {
	db := SetUp()
	defer CloseDb(db)

	index := 1

	product, err := db.Product.GetById(index)

	AssertNil(t, err)

	expected := test.GetTestProductById(index)

	if !reflect.DeepEqual(expected, product) {
		t.Errorf("Expected product: %v, got product: %v", expected, product)
	}
}

func TestGetProductsGivenIds(t *testing.T) {
	db := SetUp()
	defer CloseDb(db)

	indexes := []int{1, 2, 3}
	products, err := db.Product.GetByIds(indexes...)

	AssertNil(t, err)

	for i, index := range indexes {
		expected := test.GetTestProductById(index)
		if !reflect.DeepEqual(expected, products[i]) {
			t.Errorf("Expected product: %v, got product: %v", expected, products[i])
		}
	}
}

func TestGetCategoriesReturnsCorrectCategories(t *testing.T) {
	db := SetUp()
	defer CloseDb(db)

	categories, err := db.Product.GetCategories()

	AssertNil(t, err)

	expectedCategories := testData.GetProductTestData().Categories

	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected categories length: %d, got: %d", len(expectedCategories), len(categories))
	}

	for _, category := range categories {
		if ok, err := assertCategoryHasExpectedName(category); !ok {
			t.Errorf("Expected error: nil, got error: %s", err)
		}
	}
}

func TestGetProductsByCategoryProvidesCorrectProducts(t *testing.T) {
	db := SetUp()
	defer CloseDb(db)

	categoryId := 1

	products, err := db.Product.GetByCategory(categoryId)

	AssertNil(t, err)

	expectedProducts := test.GetTestProductsByCategory(categoryId)

	if len(expectedProducts) != len(products) {
		t.Errorf("Expected products length: %d, got: %d", len(expectedProducts), len(products))
	}

	for _, product := range products {
		if !assertProductContainedWithinSliceOfProducts(product, expectedProducts) {
			t.Errorf("Product %v not expected", product)
		}
	}
}

func TestGetProductsByText(t *testing.T) {
	db := SetUp()
	defer CloseDb(db)

	searchText := "fruit"

	products, err := db.Product.GetByText(searchText)

	AssertNil(t, err)

	expectedProducts := test.GetTestProductsByText(searchText)

	if len(products) != len(expectedProducts) {
		t.Errorf("Expected products length: %d, got: %d", len(expectedProducts), len(products))
	}

	for _, product := range products {
		if !assertProductContainedWithinSliceOfProducts(product, expectedProducts) {
			t.Errorf("Product %v not expected", product)
		}
	}
}

func TestGetProductsWithCurrentDeals(t *testing.T) {
	db := SetUp()
	defer CloseDb(db)

	currentDate := utils.GetFormattedDate()

	products, err := db.Product.GetWithCurrentDeals(currentDate)

	AssertNil(t, err)

	expectedProducts := test.GetTestProductsWithCurrentDeals(currentDate)

	if len(expectedProducts) != len(products) {
		t.Errorf("Expected products length: %d, got: %d", len(expectedProducts), len(products))
	}

	for _, product := range products {
		if !assertProductContainedWithinSliceOfProducts(product, expectedProducts) {
			t.Errorf("Product %v not expected", product)
		}
	}
}

func TestDecreaseStockReducesStockByCorrectQuantity(t *testing.T) {
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

	product, err := db.Product.GetById(productId)

	AssertNil(t, err)

	expectedProduct := test.GetTestProductById(productId)
	expectedProduct.QuantityRemaining -= quantity

	if !reflect.DeepEqual(expectedProduct, product) {
		t.Errorf("Expected product: %v, got: %v", expectedProduct, product)
	}
}

func TestDecreaseStockFailsForFakeProduct(t *testing.T) {
	db := SetUp()
	defer CloseDb(db)

	productId := 1337
	quantity := 5

	productQuantities := []utils.OrderItem{
		{
			ProductId: productId,
			Quantity:  quantity,
		},
	}

	err := db.Product.DecreaseStock(productQuantities)

	if err == nil {
		t.Errorf("Expected error: nil, got error: %s", err)
	}
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
