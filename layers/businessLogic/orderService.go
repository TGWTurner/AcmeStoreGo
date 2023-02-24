package businessLogic

import (
	"backend/utils"
	"errors"
)

func NewOrderService(orderDatabase utils.OrderDatabase, productService ProductService) *OrderService {
	return &OrderService{
		db: orderDatabase,
		ps: productService,
	}
}

func (os *OrderService) Close() {
	os.db.Close()
}

func (os *OrderService) UpdateBasket(items []utils.OrderItem, currentBasket utils.Basket) (utils.Basket, error) {
	// To implement

	return utils.Basket{}, errors.New("To implement")
}

func (os *OrderService) CreateOrder(
	customerId int,
	shippingDetails utils.ShippingDetails,
	orderItems []utils.OrderItem,
) (utils.Order, error) {
	// To implement

	return utils.Order{}, errors.New("To implement")
}

func (os *OrderService) GetOrdersByCustomerId(customerId int) ([]utils.Order, error) {
	// To implement

	return []utils.Order{}, errors.New("To implement")
}

func (os *OrderService) GetOrderByToken(orderId string) (utils.Order, error) {
	// To implement

	return utils.Order{}, errors.New("To implement")
}

type OrderService struct {
	db utils.OrderDatabase
	ps ProductService
}
