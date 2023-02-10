package api

import (
	bl "bjssStoreGo/backend/layers/businessLogic"
	"bjssStoreGo/backend/utils"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

func NewOrderApi(orderService *bl.OrderService, s *sessions.CookieStore) *OrderApi {
	return &OrderApi{
		os: *orderService,
		s:  s,
	}
}

func (o *OrderApi) validateItems(orderItems []utils.OrderItem) error {
	for _, orderItem := range orderItems {
		if orderItem.ProductId == 0 || orderItem.Quantity == 0 {
			return errors.New("Order items were invalid")
		}
	}

	return nil
}

func (o *OrderApi) GetBasket(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement get basket
}

func (o *OrderApi) PostBasket(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement post basket
}

func (o *OrderApi) GetHistory(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement get history
}

func (o *OrderApi) GetOrder(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement get order
}

func (o *OrderApi) PostCheckout(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement post checkout
}

type OrderApi struct {
	os bl.OrderService
	s  *sessions.CookieStore
}
