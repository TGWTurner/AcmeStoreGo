package api

import "bjssStoreGo/backend/layers/businessLogic"

func NewOrderApi(orderService *businessLogic.OrderService) *OrderApi {
	return &OrderApi{
		os: *orderService,
	}
}

func validateItems(items []string /*TODO itemsStruct*/) {
	//TODO validateItems fn
}

func (o OrderApi) getBasket(req string, res string) {
	//TODO: Implement get basket
}

func (o OrderApi) postBasket(req string, res string) {
	//TODO: Implement post basket
}

func (o OrderApi) getHistory(req string, res string) {
	//TODO: Implement get history
}

func (o OrderApi) getOrder(req string, res string) {
	//TODO: Implement get order
}

func (o OrderApi) postCheckout(req string, res string) {
	//TODO: Implement post checkout
}

type OrderApi struct {
	os businessLogic.OrderService
}
