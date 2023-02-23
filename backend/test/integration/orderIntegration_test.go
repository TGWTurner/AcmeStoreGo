package integration

import (
	"bjssStoreGo/backend/layers/api"
	"bjssStoreGo/backend/test"
	"bjssStoreGo/backend/utils"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestCreatesAnOrder(t *testing.T) {
	w := test.SetUpApi()
	apiRequester := test.NewApiRequester(w)
	defer w.Close()

	orderRequest := makeOrderRequest()
	body, err := json.Marshal(orderRequest)

	test.AssertNil(t, err)

	response := apiRequester.Post("/api/order/checkout", body)

	test.AssertResponseCode(t, 200, response.Code)

	var actual utils.Order

	err = json.NewDecoder(response.Body).Decode(&actual)

	test.AssertNil(t, err)
}

func TestFetchACreatedOrder(t *testing.T) {
	w := test.SetUpApi()
	apiRequester := test.NewApiRequester(w)
	defer w.Close()

	//First create an order
	orderRequest := makeOrderRequest()
	body, err := json.Marshal(orderRequest)

	test.AssertNil(t, err)

	response := apiRequester.Post("/api/order/checkout", body)

	test.AssertResponseCode(t, 200, response.Code)

	var responseOrder utils.Order

	err = json.NewDecoder(response.Body).Decode(&responseOrder)
	test.AssertNil(t, err)

	//Then get its ID and try and fetch it
	requestUrl := fmt.Sprintf("/api/order/%s", responseOrder.Id)
	createdOrder := apiRequester.Get(requestUrl)

	test.AssertResponseCode(t, 200, createdOrder.Code)

	var actual utils.Order
	err = json.NewDecoder(createdOrder.Body).Decode(&actual)
	test.AssertNil(t, err)

	if !reflect.DeepEqual(orderRequest.Items[0], actual.Items[0]) {
		t.Errorf("Expected first order item to be: %v, got: %v", orderRequest.Items[0], actual.Items[0])
	}
}

func TestGetsAnEmptyBasket(t *testing.T) {
	w := test.SetUpApi()
	apiRequester := test.NewApiRequester(w)
	defer w.Close()

	emptyBasket := apiRequester.Get("/api/order/basket")

	test.AssertResponseCode(t, 200, emptyBasket.Code)

	var actual utils.Basket
	err := json.NewDecoder(emptyBasket.Body).Decode(&actual)
	test.AssertNil(t, err)
}

func TestUpdatesABasket(t *testing.T) {
	w := test.SetUpApi()
	apiRequester := test.NewApiRequester(w)
	defer w.Close()

	// A basket updae is the same as the items part of an order, so re-use our test data generator function
	basket := utils.Basket{
		Items: makeOrderRequest().Items,
	}

	// Add an item to the basket. The response should be our item and a total
	body, err := json.Marshal(basket)
	test.AssertNil(t, err)

	added := apiRequester.Post("/api/order/basket", body)
	test.AssertResponseCode(t, 200, added.Code)

	var addedBasket utils.Basket
	err = json.NewDecoder(added.Body).Decode(&addedBasket)
	test.AssertNil(t, err)

	if 2090 != addedBasket.Total {
		t.Errorf("Expected total: 2090, got: %d", addedBasket.Total)
	}
	test.AssertOrderItemsMatch(t, basket.Items, addedBasket.Items)

	// Finally, fetch our basket again. It should have the item we added.
	updated := apiRequester.Get("/api/order/basket")
	test.AssertResponseCode(t, 200, updated.Code)

	var updatedBasket utils.Basket
	err = json.NewDecoder(updated.Body).Decode(&updatedBasket)
	test.AssertNil(t, err)

	if 2090 != addedBasket.Total {
		t.Errorf("Expected total: 2090, got: %d", addedBasket.Total)
	}
	test.AssertOrderItemsMatch(t, basket.Items, addedBasket.Items)
}

// NOTE - The this relies on sign-in working. We can only list our own orders.
func TestListsOrderHistory(t *testing.T) {
	w := test.SetUpApi()
	apiRequester := test.NewApiRequester(w)
	defer w.Close()

	test.SignInPrePopulatedUser(t, apiRequester)

	// Add two orders
	orderRequest := makeOrderRequest()

	order1 := sendOrderRequest(t, apiRequester, orderRequest)
	order2 := sendOrderRequest(t, apiRequester, orderRequest)

	orders := apiRequester.Get("/api/order/history")
	test.AssertResponseCode(t, 200, orders.Code)

	var actual []utils.Order
	err := json.NewDecoder(orders.Body).Decode(&actual)
	test.AssertNil(t, err)

	order1Found := false
	order2Found := false
	for _, order := range actual {
		if order.Id == order1.Id {
			order1Found = true
			continue
		}
		if order.Id == order2.Id {
			order2Found = true
			continue
		}
	}

	if !order1Found || !order2Found {
		t.Errorf("Expected returned orders to include sent orders")
	}
}

// NOTE - this test assumes Sign In works
func TestListsOnlyOrderHistoryForSignedInUser(t *testing.T) {
	w := test.SetUpApi()
	defer w.Close()

	user1 := test.NewApiRequester(w)
	test.SignInPrePopulatedUser(t, user1)

	order1 := makeOrderRequest()
	order1Response := sendOrderRequest(t, user1, order1)

	user2 := test.NewApiRequester(w)
	test.SignInPrePopulatedUser(t, user2, 2)

	order2 := makeOrderRequest()
	order2Response := sendOrderRequest(t, user2, order2)

	history1 := user1.Get("/api/order/history")
	history2 := user2.Get("/api/order/history")
	test.AssertResponseCode(t, 200, history1.Code)
	test.AssertResponseCode(t, 200, history2.Code)

	var orders1 []utils.Order
	err := json.NewDecoder(history1.Body).Decode(&orders1)
	test.AssertNil(t, err)

	var orders2 []utils.Order
	err = json.NewDecoder(history2.Body).Decode(&orders2)
	test.AssertNil(t, err)

	assertOrderSetIncludesOrder(t, orders1, order1Response)
	assertOrderSetIncludesOrder(t, orders2, order2Response)

	assertOrderSetDoesNotIncludeOrder(t, orders1, order2Response)
	assertOrderSetDoesNotIncludeOrder(t, orders2, order1Response)
}

func assertOrderSetIncludesOrder(t *testing.T, orderSet []utils.Order, expected utils.Order) {
	for _, order := range orderSet {
		if order.Id == expected.Id {
			return
		}
	}

	t.Errorf("Expected order %v to be present in response, it was not", expected)
}

func assertOrderSetDoesNotIncludeOrder(t *testing.T, orderSet []utils.Order, expected utils.Order) {
	for _, order := range orderSet {
		if order.Id == expected.Id {
			t.Errorf("Expected order %v to not be present in response, it was found", expected)
		}
	}
}

func sendOrderRequest(t *testing.T, requester *test.ApiRequester, order api.OrderRequest) utils.Order {
	orderBody, err := json.Marshal(order)
	test.AssertNil(t, err)
	order1Response := requester.Post("/api/order/checkout", orderBody)
	test.AssertResponseCode(t, 200, order1Response.Code)

	fmt.Println("Response")
	fmt.Println(order1Response.Code, order1Response.Body)

	var orderObject utils.Order
	err = json.NewDecoder(order1Response.Body).Decode(&orderObject)
	test.AssertNil(t, err)

	return orderObject
}

func makeOrderRequest() api.OrderRequest {
	return api.OrderRequest{
		PaymentToken: "someTokenToCheckWithPaymentGateway",
		ShippingDetails: utils.ShippingDetails{
			Email:    "a@example.com",
			Name:     "a",
			Address:  "b",
			Postcode: "abc123",
		},
		Items: []utils.OrderItem{
			{ProductId: 2, Quantity: 2},
			{ProductId: 3, Quantity: 1},
		},
	}
}
