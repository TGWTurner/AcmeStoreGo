package dbTests

import (
	"backend/utils"
	"errors"
	"reflect"
	"testing"
)

func TestCreateOrder(t *testing.T) {
	db := SetUp()
	defer CloseDb(db)

	_, err := createOrder(&db, 2)

	AssertNil(t, err)
}

func TestSucceedsCreatingTwoIdenticalOrders(t *testing.T) {
	db := SetUp()
	defer CloseDb(db)

	_, err := createOrder(&db, 2)

	AssertNil(t, err)

	_, err = createOrder(&db, 2)

	AssertNil(t, err)
}

func TestGetOrderByToken(t *testing.T) {
	db := SetUp()
	defer CloseDb(db)

	customerId := 2

	expectedOrder, err := createOrder(&db, customerId)

	AssertNil(t, err)

	order, err := db.Order.GetByToken(expectedOrder.Id)

	AssertNil(t, err)

	if !reflect.DeepEqual(expectedOrder, order) {
		t.Errorf("Expected order: %v, got order: %v", expectedOrder, order)
	}
}

func TestFailsToGetOrderWithFakeToken(t *testing.T) {
	db := SetUp()
	defer CloseDb(db)

	token := "ThisIsAFakeToken"
	_, err := db.Order.GetByToken(token)

	expectedErr := errors.New("Order does not exist with orderToken: " + token)

	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error: %s, got error: %s", expectedErr, err)
	}
}

func TestGetOrdersByCustomerId(t *testing.T) {
	db := SetUp()
	defer CloseDb(db)

	customerId := 2

	createOrder(&db, customerId)
	createOrder(&db, customerId)

	orders, err := db.Order.GetByCustomerId(customerId)

	AssertNil(t, err)

	if len(orders) != 2 {
		t.Errorf("Expected orders length: 2, got: %d", len(orders))
	}
}

func TestFailsToGetOrderByFakeCustomerId(t *testing.T) {
	db := SetUp()
	defer CloseDb(db)

	orders, err := db.Order.GetByCustomerId(1337)

	AssertNil(t, err)

	if len(orders) != 0 {
		t.Errorf("Expected orders length: 2, got: %d", len(orders))
	}
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
