package main

import (
	"bjssStoreGo/backend/layers/dataAccess"
	"bjssStoreGo/backend/layers/dataAccess/sqlite"
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/exp/slices"
)

func setUp() sqlite.Database {
	file := "./sqlite.db"

	if _, err := os.Stat(file); !os.IsNotExist(err) {
		if err := os.Remove(file); err != nil {
			fmt.Println(err)
			panic("Failed to remove file: " + file)
		}
	}
	db := dataAccess.InitiateConnection()

	return db
}

func closeDbs(db sqlite.Database) {
	db.Order.Close()
	db.Account.Close()
	db.Product.Close()
}

func main() {
	testCreateOrder()
	testGetOrderByToken()
	testGetOrdersByCustomerId()
	testGetProductGivenId()
	testGetProductsGivenIds()
	testGetCategoriesReturnsCorrectCategories()
	testGetProductsByCategoryProvidesCorrectProducts()
}

/*
TESTING:
	- Products -
	1. Get All Products
	2. Get Product given an Id
	3. Get Product given []Id's
	4. Get All categories
	5. Get Products given Id
	6. Search Products "Dog"
	7. Get Products with deals
	8. Decrease stock

	- Orders -
	1. Create order
	2. Get Order from customer Id
	3. Get Order from token

	- Accounts -
	1. Add Account
	2. Get Account by email
	3. Get Account by id
	4. Update an Account
*/

func testCreateOrder() {
	db := setUp()
	defer closeDbs(db)

	createOrder(db, 2)
	fmt.Println("TEST PASSED -- Order created")
}

func testGetOrderByToken() {
	db := setUp()
	defer closeDbs(db)

	customerId := 2

	order := createOrder(db, customerId)

	dbOrder := db.Order.GetByToken(order.Id)

	if order.Total != dbOrder.Total {
		fmt.Println("TEST FAILED -- Order total response != provided")
		return
	}
	if order.CustomerId != dbOrder.CustomerId {
		fmt.Println("TEST FAILED -- Order customer ir=d response != provided")
		return
	}
	if order.ShippingDetails.Address != dbOrder.ShippingDetails.Address {
		fmt.Println("TEST FAILED -- Order address response != provided")
		return
	}
	if order.ShippingDetails.Email != dbOrder.ShippingDetails.Email {
		fmt.Println("TEST FAILED -- Order email response != provided")
		return
	}
	if order.ShippingDetails.Name != dbOrder.ShippingDetails.Name {
		fmt.Println("TEST FAILED -- Order name response != provided")
		return
	}
	if order.Items[0] != dbOrder.Items[0] {
		fmt.Println("TEST FAILED -- Order item[0] response != provided")
		return
	}
	if order.Items[1] != dbOrder.Items[1] {
		fmt.Println("TEST FAILED -- Order item[1] response != provided")
		return
	}

	fmt.Println("TEST PASSED -- Order correctly retrieved for token: " + order.Id)
}

func testGetOrdersByCustomerId() {
	db := setUp()
	defer closeDbs(db)

	createOrder(db, 2)
	createOrder(db, 2)

	orders := db.Order.GetByCustomerId(2)

	if len(orders) != 2 {
		fmt.Println("TEST FAILED -- Recieved wrong number of orders")
		fmt.Printf("\t expected: 2 actual: " + strconv.Itoa(len(orders)))
		return
	}

	fmt.Println("TEST PASSED -- Got correct number of orders")
}

func createOrder(db sqlite.Database, customerId int) utils.Order {
	order := utils.Order{
		Total: 5,
		ShippingDetails: utils.ShippingDetails{
			Email:    "testEmail",
			Name:     "testName",
			Address:  "testAddress",
			Postcode: "testPostcode",
		},
		Items: []utils.OrderItem{
			{
				ProductId: 1,
				Quantity:  10,
			},
			{
				ProductId: 5,
				Quantity:  1,
			},
		},
	}

	return db.Order.Add(customerId, order)
}

func testGetProductGivenId() {
	db := setUp()
	defer closeDbs(db)

	index := 1

	product := db.Product.GetByIds(index)[0]
	expected := getTestProductById(index)

	if ok, err := assertProductExpectedMatchesActual(expected, product); !ok {
		fmt.Println("TEST FAILED -- " + err)
		return
	}

	fmt.Println("TEST PASSED -- Received correct product for index: ", index)
}

func testGetProductsGivenIds() {
	db := setUp()
	defer closeDbs(db)

	indexes := []int{1, 2, 3}

	products := db.Product.GetByIds(indexes...)

	for i, index := range indexes {
		expected := getTestProductById(index)
		if ok, err := assertProductExpectedMatchesActual(expected, products[i]); !ok {
			fmt.Println("TEST FAILED -- " + err)
			return
		}
	}

	fmt.Println("TEST PASSED -- Got correct products for ids: ", indexes)
}

func testGetCategoriesReturnsCorrectCategories() {
	db := setUp()
	defer closeDbs(db)

	categories := db.Product.GetCategories()
	expectedCategories := testData.GetProductTestData().Categories

	if len(categories) != len(expectedCategories) {
		fmt.Println("TEST FAILED -- Received incorrect number of categories")
		return
	}

	for _, category := range categories {
		if ok, err := assertCategoryHasExpectedName(category); !ok {
			fmt.Println("TEST FAILED -- " + err)
			return
		}
	}

	fmt.Println("TEST PASSED -- Recieved expected categories")
}

func testGetProductsByCategoryProvidesCorrectProducts() {
	db := setUp()
	defer closeDbs(db)

	categoryId := 1

	products := db.Product.GetProductsByCategory(categoryId)

	for _, product := range products {
		if categoryId != product.CategoryId {
			fmt.Println(
				"TEST FAILED -- Product had wrong Category Id expected: " +
					strconv.Itoa(categoryId) +
					" actual: " +
					strconv.Itoa(product.Id),
			)
		}

		if ok, err := assertProductExpectedMatchesActual(getTestProductById(product.Id), product); !ok {
			fmt.Println("TEST FAILED -- " + err)
			return
		}
	}

	fmt.Println("TEST PASSED -- Recieved correct products for category: " + strconv.Itoa(categoryId))
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

func assertProductExpectedMatchesActual(expected utils.Product, actual utils.Product) (bool, string) {
	if expected.Id != actual.Id {
		return false, "Product Id did not match, expected: " + strconv.Itoa(expected.Id) + " actual: " + strconv.Itoa(actual.Id)
	}
	if expected.QuantityRemaining != actual.QuantityRemaining {
		return false, "Quantity Remaining did not match, expected: " + strconv.Itoa(expected.Id) + " actual: " + strconv.Itoa(actual.Id)

	}
	if expected.CategoryId != actual.CategoryId {
		return false, "Category Id did not match, expected: " + strconv.Itoa(expected.Id) + " actual: " + strconv.Itoa(actual.Id)

	}
	if expected.Price != actual.Price {
		return false, "Price did not match, expected: " + strconv.Itoa(expected.Id) + " actual: " + strconv.Itoa(actual.Id)
	}
	if expected.ShortDescription != actual.ShortDescription {
		return false, "Short Description did not match"
	}
	if expected.LongDescription != actual.LongDescription {
		return false, "Long Description did not match"
	}

	return true, ""
}
