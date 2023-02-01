package businessLogic

import "bjssStoreGo/backend/utils"

func NewOrderService(orderDatabase utils.OrderDatabase) *OrderService {
	return &OrderService{
		orderDb: orderDatabase,
	}
}

func (os OrderService) updateBasket(basket []string) {
	//TODO: Implement update basket
}

func (os OrderService) createOrder(customerId int, orderRequest string) {
	//TODO: Implement create order
}

func (os OrderService) getOrdersByCustomerId(customerId int) {
	//TODO: Implement get orders by customer id
}

func (os OrderService) getOrderByToken(token string) {
	//TODO: Implement get order by token
}

type OrderService struct {
	orderDb utils.OrderDatabase
}
