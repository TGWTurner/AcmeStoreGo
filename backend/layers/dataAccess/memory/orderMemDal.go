package memory

import (
	"bjssStoreGo/backend/utils"
	"time"
)

func NewOrderDatabase() OrderDatabase {
	return OrderDatabase{
		orders: []utils.Order{},
	}
}

func (od *OrderDatabase) GetOrdersByCustomerId(customerId int) []utils.Order {
	//need to search and get all the objects from the slice which have this customer id
	orders := []utils.Order{}

	for _, order := range od.orders {
		if order.CustomerId == customerId {
			orders = append(orders, order)
		}
	}

	return orders
}

func (od *OrderDatabase) AddOrder(customerId int, order utils.Order) {
	order.UpdatedDate = time.Now().String()
	order.CustomerId = customerId

	od.orders = append( //add id too?
		od.orders,
		order,
	)

	//TODO: return the added order?
}

type OrderDatabase struct {
	orders []utils.Order
}
