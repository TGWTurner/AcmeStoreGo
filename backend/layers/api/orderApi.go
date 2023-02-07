package api

import (
	"bjssStoreGo/backend/layers/businessLogic"
	"bjssStoreGo/backend/utils"
	"net/http"
)

func NewOrderApi(orderService *businessLogic.OrderService) *OrderApi {
	return &OrderApi{
		os: *orderService,
	}
}

func validateItems(orderItems []utils.OrderItem) {
	//TODO validateItems fn
}

func (o OrderApi) GetBasket(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement get basket
}

func (o OrderApi) PostBasket(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement post basket
}

func (o OrderApi) GetHistory(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement get history
}

func (o OrderApi) GetOrder(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement get order
}

func (o OrderApi) PostCheckout(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement post checkout
}

type OrderApi struct {
	os businessLogic.OrderService
}
