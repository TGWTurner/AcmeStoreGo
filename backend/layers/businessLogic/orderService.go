package businessLogic

import "bjssStoreGo/backend/utils"

func NewOrderService(orderDatabase utils.OrderDatabase) *OrderService {
	return &OrderService{
		db: orderDatabase,
	}
}

func (os OrderService) Close() {
	os.db.Close()
}

func (os OrderService) UpdateBasket(basket []string) {
	//TODO: Implement update basket
}

func (os OrderService) CreateOrder(customerId int, orderRequest string) {
	//TODO: Implement create order
}

func (os OrderService) GetOrdersByCustomerId(customerId int) {
	//TODO: Implement get orders by customer id
}

func (os OrderService) GetOrderByToken(token string) {
	//TODO: Implement get order by token
}

type OrderService struct {
	db utils.OrderDatabase
}
