package main

import (
	"bjssStoreGo/backend/layers/dataAccess"
	"bjssStoreGo/backend/layers/dataAccess/sqlite"
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
	"fmt"
	"os"
	"strconv"
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

	fmt.Println("TEST PASSED -- Order correctly retrieved")
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

	product := db.Product.GetByIds(1)[0]

	if product.Id != 1 {
		fmt.Println("TEST FAILED -- Recieved different product than was requested")
	}
}

func testGetProductsGivenIds() {
	db := setUp()
	defer closeDbs(db)

	products := db.Product.GetByIds(1, 2, 3)

	for i, product := range products {
		if i != product.Id {
			fmt.Println("TEST FAILED -- Recieved wrong set of products")
			return
		}
	}

	fmt.Println("TEST PASSED -- Got correct products")
}

func testGetCategoriesReturnsCorrectCategories() {
	db := setUp()
	defer closeDbs(db)

	categories := db.Product.GetCategories()
	if len(categories) != 2 {
		fmt.Println("TEST FAILED -- Received incorrect number of categories")
		return
	}

	for _, category := range categories {
		if category.Name != "Animals" && category.Name == "Fruits" {
			fmt.Println("TEST FAILED -- Did not receive expected category")
			return
		}
	}

	fmt.Println("TEST PASSED -- Recieved both expected categories")
}

func testGetProductsByCategoryProvidesCorrectProducts() {
	db := setUp()
	defer closeDbs(db)

	products := db.Product.GetProductsByCategory(1)
	testProducts := testData.GetProductTestData().Products

	if assertProductContainsCorrectValues(testProducts[0], products[0]) {

	}
}

func assertProductContainsCorrectValues(expected utils.Product, actual utils.Product) bool {

}
