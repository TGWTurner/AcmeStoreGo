package test

import (
	bl "bjssStoreGo/backend/layers/businessLogic"
	da "bjssStoreGo/backend/layers/dataAccess"
	"bjssStoreGo/backend/utils"
	"testing"
)

func setUpOrder() bl.OrderService {
	db := da.InitiateConnection()
	ps := bl.NewProductService(db.Product)
	return *bl.NewOrderService(db.Order, *ps)
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

func TestUpdatesBasket(t *testing.T) {
	os := setUpOrder()
	defer os.Close()

	orderItems := makeTestOrderRequest.items

	basket, err := os.UpdateBasket(orderItems)

	AssertNil(t, err)

	if basket.Total != 20 {

	}
}
