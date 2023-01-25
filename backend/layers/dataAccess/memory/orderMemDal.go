package memory

import (
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
	"time"
)

func NewOrderDatabase() OrderDatabase {
	testOrders := testData.GetOrderTestData()

	return OrderDatabase{
		orders: testOrders,
	}
}

func (od *OrderDatabase) GetOrdersByCustomerId(customerId int) []utils.Order {
	orders := []utils.Order{}

	for _, order := range od.orders {
		if order.CustomerId == customerId {
			orders = append(orders, order)
		}
	}

	return orders
}

func (od *OrderDatabase) GetOrderByToken(orderId string) utils.Order {
	for _, order := range od.orders {
		if order.Id == orderId {
			return order
		}
	}

	panic("Failed to get orders for order Token: " + orderId)
}

func (od *OrderDatabase) AddOrder(customerId int, order utils.Order) utils.Order {
	order.UpdatedDate = time.Now().String()
	order.CustomerId = customerId
	order.Id = utils.UrlSafeUniqueId()

	od.orders = append(
		od.orders,
		order,
	)

	return order
}

type OrderDatabase struct {
	orders []utils.Order
}
