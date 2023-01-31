package tests

import (
	"bjssStoreGo/backend/utils"
	"reflect"
	"strconv"
)

func TestCreateOrder() (bool, string, string) {
	db := SetUp()
	defer CloseDb(db)

	_, err := createOrder(&db, 2)

	if err != nil {
		return false, "testGetOrderByToken", "Failed to create order"
	}

	return true, "testCreateOrder", "Order Created"
}

func TestSucceedsCreatingTwoIdenticalOrders() (bool, string, string) {
	db := SetUp()
	defer CloseDb(db)

	_, err := createOrder(&db, 2)

	if err != nil {
		return false, "testSucceedsCreatingTwoIdenticalOrders", "Failed to create first order"
	}

	_, err = createOrder(&db, 2)

	if err != nil {
		return false, "testSucceedsCreatingTwoIdenticalOrders", "Failed to create second order"
	}

	return true, "testSucceedsCreatingTwoIdenticalOrders", "Successfully placed two identical orders"
}

func TestGetOrderByToken() (bool, string, string) {
	db := SetUp()
	defer CloseDb(db)

	customerId := 2

	expectedOrder, err := createOrder(&db, customerId)

	if err != nil {
		return false, "testGetOrderByToken", "Failed to create order"
	}

	order, err := db.Order.GetByToken(expectedOrder.Id)

	if err != nil {
		return false, "testGetOrderByToken", "Failed to get order by token"
	}

	if !reflect.DeepEqual(expectedOrder, order) {
		return false, "testGetOrderByToken", "Actual Order did not match expected"
	}

	return true, "testGetOrderByToken", "Order correctly retrieved for token: " + expectedOrder.Id
}

func TestFailsToGetOrderWithFakeToken() (bool, string, string) {
	db := SetUp()
	defer CloseDb(db)

	token := "ThisIsAFakeToken"
	_, err := db.Order.GetByToken(token)

	if err.Error() != "Order does not exist with orderToken: "+token {
		return false, "testFailsToGetOrderWithFakeToken", "Failed to get order by token"
	}

	return true, "testFailsToGetOrderWithFakeToken", "Succeeded to return empty order for fake token"
}

func TestGetOrdersByCustomerId() (bool, string, string) {
	db := SetUp()
	defer CloseDb(db)

	customerId := 2

	createOrder(&db, customerId)
	createOrder(&db, customerId)

	orders, err := db.Order.GetByCustomerId(customerId)

	if err != nil {
		return false, "testGetOrdersByCustomerId", "Failed to get order by customer id"
	}

	if len(orders) != 2 {
		return false, "testGetOrdersByCustomerId", "Recieved wrong number of orders, expected: 2 actual: " + strconv.Itoa(len(orders))
	}

	return true, "testGetOrderByToken", "Got correct number of orders"
}

func TestFailsToGetOrderByFakeCustomerId() (bool, string, string) {
	db := SetUp()
	defer CloseDb(db)

	orders, err := db.Order.GetByCustomerId(1337)

	if err != nil {
		return false, "testFailsToGetOrderByFakeCustomerId", "Failed to get order by customer id"
	}

	if len(orders) != 0 {
		return false, "testFailsToGetOrderByFakeCustomerId", "Recieved wrong number of orders, expected: 0 actual: " + strconv.Itoa(len(orders))
	}

	return true, "testFailsToGetOrderByFakeCustomerId", "Successfully got no orders for non existant customer"
}

func createOrder(db *utils.Database, customerId int) (utils.Order, error) {
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
