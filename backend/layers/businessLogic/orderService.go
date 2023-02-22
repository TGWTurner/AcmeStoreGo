package businessLogic

import (
	"bjssStoreGo/backend/utils"
	"errors"
	"fmt"
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
	notEnoughStock, total, err := os.ps.CheckStock(items)

	if err != nil {
		return utils.Basket{}, err
	}

	if len(notEnoughStock) > 0 {
		msg := fmt.Sprintf("Not enough stock for: %v", notEnoughStock)
		return utils.Basket{}, errors.New(msg)
	}

	currentBasket.Items = items

	currentBasket.Total = total

	return currentBasket, nil
}

func (os *OrderService) CreateOrder(
	customerId int,
	shippingDetails utils.ShippingDetails,
	orderItems []utils.OrderItem,
) (utils.Order, error) {

	notEnoughStock, total, err := os.ps.CheckStock(orderItems)

	if err != nil {
		return utils.Order{}, err
	}

	if len(notEnoughStock) > 0 {
		msg := fmt.Sprintf("Trying to decrease stock of %v below zero", notEnoughStock)
		return utils.Order{}, errors.New(msg)
	}

	if err := os.ps.DecreaseStock(orderItems); err != nil {
		return utils.Order{}, err
	}

	order := utils.Order{
		Total:           total,
		UpdatedDate:     utils.GetFormattedDate(),
		CustomerId:      customerId,
		ShippingDetails: shippingDetails,
		Items:           orderItems,
	}

	return os.db.Add(customerId, order)
}

func (os *OrderService) GetOrdersByCustomerId(customerId int) ([]utils.Order, error) {
	return os.db.GetByCustomerId(customerId)
}

func (os *OrderService) GetOrderByToken(orderId string) (utils.Order, error) {
	return os.db.GetByToken(orderId)
}

type OrderService struct {
	db utils.OrderDatabase
	ps ProductService
}
