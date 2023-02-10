package businessLogic

import "bjssStoreGo/backend/utils"

func NewOrderService(orderDatabase utils.OrderDatabase) *OrderService {
	return &OrderService{
		db: orderDatabase,
	}
}

func (os *OrderService) Close() {
	os.db.Close()
}

func (os *OrderService) UpdateBasket(items []utils.OrderItem) (utils.Basket, error) {
	//TODO: Implement update basket
	/*
		Depends on how we store the basket
		Session:
		 - Get session
		 - Get basket from session.Values
		 - Update with new items
		 - Set session.Values[basket] to the updated version
		 - Save the session with session.Save(r, w)
	*/
	return utils.Basket{}, nil
}

func (os *OrderService) CreateOrder(
	customerId int,
	shippingDetails utils.ShippingDetails,
	orderItems []utils.OrderItem,
) (utils.Order, error) {
	//TODO: Implement create order
	/*
		Process call productService.decreaseStock
	*/

	return utils.Order{}, nil
}

func (os *OrderService) GetOrdersByCustomerId(customerId int) ([]utils.Order, error) {
	//TODO: Implement get orders by customer id

	/*
		Process: Call order db.GetByCustomerId
	*/
	return os.db.GetByCustomerId(customerId)
}

func (os *OrderService) GetOrderByToken(orderId string) (utils.Order, error) {
	//TODO: Implement get order by token
	return os.db.GetByToken(orderId)
}

type OrderService struct {
	db utils.OrderDatabase
}
