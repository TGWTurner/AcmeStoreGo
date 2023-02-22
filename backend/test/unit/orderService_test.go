package unit

import (
	bl "bjssStoreGo/backend/layers/businessLogic"
	da "bjssStoreGo/backend/layers/dataAccess"
	"bjssStoreGo/backend/test"
	"bjssStoreGo/backend/utils"
	"reflect"
	"testing"
	"time"
)

func setUpOrder() (bl.OrderService, bl.ProductService) {
	db := da.InitiateConnection()
	ps := bl.NewProductService(db.Product)
	return *bl.NewOrderService(db.Order, *ps), *ps
}

var makeTestOrderRequest = struct {
	shippingDetails utils.ShippingDetails
	items           []utils.OrderItem
}{
	shippingDetails: utils.ShippingDetails{
		Email:    "a@example.com",
		Name:     "a",
		Address:  "b",
		Postcode: "abc123",
	},
	items: []utils.OrderItem{
		{ProductId: 2, Quantity: 2},
		{ProductId: 3, Quantity: 3},
	},
}

func xTestUpdatesBasket(t *testing.T) {
	os, ps := setUpOrder()
	defer func() { os.Close(); ps.Close() }()

	orderItems := makeTestOrderRequest.items

	currentBasket := utils.Basket{
		Total: 0,
		Items: []utils.OrderItem{},
	}

	basket := os.UpdateBasket(orderItems, currentBasket)

	expected := getTotalFromOrderItems(t, ps, orderItems)

	if expected != basket.Total {
		t.Errorf("Total incorrect, expected: %d got: %d", expected, basket.Total)
	}
}

func xTestCreatesAnOrder(t *testing.T) {
	os, ps := setUpOrder()
	defer func() { os.Close(); ps.Close() }()

	request := makeTestOrderRequest
	response, err := os.CreateOrder(1, request.shippingDetails, request.items)

	test.AssertNil(t, err)

	if 1 <= len(response.Id) {
		t.Errorf("Expected response id to be > 1, got %d", len(response.Id))
	}

	if 1 != response.CustomerId {
		t.Errorf("Expected customerId to be 1, got %d", response.CustomerId)
	}

	if request.shippingDetails != response.ShippingDetails {
		t.Errorf("Expected shipping details to match, they did not")
	}

	if reflect.DeepEqual(request.items, response.Items) {
		t.Errorf("Expected order items to match, they did not")
	}

	if time.Now().Format("2006-01-02") < response.UpdatedDate {
		t.Errorf("Expected updatedDate to be less than current datetime")
	}
}

func xTestRejectsAnOrderIfNotEnoughStock(t *testing.T) {
	os, ps := setUpOrder()
	defer func() { os.Close(); ps.Close() }()

	request := makeTestOrderRequest

	request.items = append(request.items, utils.OrderItem{
		ProductId: 4, Quantity: 9,
	})

	response, err := os.CreateOrder(1, request.shippingDetails, request.items)

	if err == nil {
		t.Errorf("Expected error, got: nil")
	}

	if err.Error() != "stock-level" {
		t.Errorf("Expected stock level error, got: %s", err)
	}

	expectedItems := []utils.OrderItem{
		{ProductId: 4, Quantity: 9},
	}

	if reflect.DeepEqual(expectedItems, response.Items) {
		t.Errorf("Expected order items to match, they did not")
	}
}

func xTestFetchesOrders(t *testing.T) {
	os, ps := setUpOrder()
	defer func() { os.Close(); ps.Close() }()

	request1 := makeTestOrderRequest
	request1.items = []utils.OrderItem{{ProductId: 2, Quantity: 1}}

	request2 := makeTestOrderRequest
	request2.items = []utils.OrderItem{{ProductId: 3, Quantity: 1}}

	request3 := makeTestOrderRequest
	request3.items = []utils.OrderItem{{ProductId: 4, Quantity: 1}}

	response1, err := os.CreateOrder(1, request1.shippingDetails, request1.items)
	test.AssertNil(t, err)

	response2, err := os.CreateOrder(2, request2.shippingDetails, request2.items)
	test.AssertNil(t, err)

	response3, err := os.CreateOrder(1, request3.shippingDetails, request3.items)
	test.AssertNil(t, err)

	orders, err := os.GetOrdersByCustomerId(1)
	test.AssertNil(t, err)

	if 2 != len(orders) {
		t.Errorf("Expected 2 orders, got")
	}

	order1 := getOrderFromOrders(t, orders, response1.Id)
	if response1.Id != order1.Id {
		t.Errorf("Expected orderToken: %s, got: %s", response1.Id, order1.Id)
	}

	order3 := getOrderFromOrders(t, orders, response3.Id)
	if response3.Id != order3.Id {
		t.Errorf("Expected orderToken: %s, got: %s", response3.Id, order3.Id)
	}

	if response3.Total != order3.Total {
		t.Errorf("Expected total: %d, got: %d", response3.Total, order3.Total)
	}

	if time.Now().Format("2006-01-02") < order3.UpdatedDate {
		t.Errorf("Expected updatedDate to be less than current datetime")
	}

	if 0 >= len(order3.Items) {
		t.Errorf("Expected length of orders to be > 0, got: %d", len(order3.Items))
	}

	byToken, err := os.GetOrderByToken(response2.Id)
	test.AssertNil(t, err)

	if !reflect.DeepEqual(response2, byToken) {
		t.Errorf("Expected order by token to match response, it did not")
	}
}

func getOrderFromOrders(t *testing.T, orders []utils.Order, orderId string) utils.Order {
	for _, order := range orders {
		if order.Id == orderId {
			return order
		}
	}

	t.Errorf("Failed to find order for customer %d with token %s", orders[0].CustomerId, orderId)
	return utils.Order{}
}

func getTotalFromOrderItems(t *testing.T, ps bl.ProductService, orderItems []utils.OrderItem) int {
	total := 0
	allProducts, err := ps.SearchProducts(map[string]string{})

	test.AssertNil(t, err)

	for _, orderItem := range orderItems {
		for _, product := range allProducts {
			if product.Id == orderItem.ProductId {
				total += product.Price * orderItem.Quantity
			}
		}
	}

	return total
}
