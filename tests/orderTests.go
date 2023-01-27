package tests

import (
	"bjssStoreGo/backend/utils"
	"reflect"
	"strconv"
)

func TestCreateOrder() {
	db := SetUp()
	defer CloseDb(db)

	createOrder(&db, 2)
	PrintTestResult(true, "testCreateOrder", "Order Created")
}

func TestGetOrderByToken() {
	db := SetUp()
	defer CloseDb(db)

	customerId := 2

	expectedOrder := createOrder(&db, customerId)

	order := db.Order.GetByToken(expectedOrder.Id)

	if !reflect.DeepEqual(expectedOrder, order) {
		PrintTestResult(false, "testGetOrderByToken", "Actual Order did not match expected")
		return
	}

	PrintTestResult(true, "testGetOrderByToken", "Order correctly retrieved for token: "+expectedOrder.Id)
}

func TestGetOrdersByCustomerId() {
	db := SetUp()
	defer CloseDb(db)

	createOrder(&db, 2)
	createOrder(&db, 2)

	orders := db.Order.GetByCustomerId(2)

	if len(orders) != 2 {
		PrintTestResult(false, "testGetOrdersByCustomerId", "Recieved wrong number of orders, expected: 2 actual: "+strconv.Itoa(len(orders)))
		return
	}

	PrintTestResult(true, "testGetOrderByToken", "Got correct number of orders")
}

func createOrder(db *utils.Database, customerId int) utils.Order {
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
