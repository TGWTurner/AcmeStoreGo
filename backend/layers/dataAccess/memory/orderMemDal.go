package memory

import (
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
	"time"
)

func NewOrderDatabase() *OrderDatabaseImpl {
	testOrders := testData.GetOrderTestData()

	return &OrderDatabaseImpl{
		orders: testOrders,
	}
}

func (ad *OrderDatabaseImpl) Close() {}

func (od *OrderDatabaseImpl) GetByCustomerId(customerId int) []utils.Order {
	orders := []utils.Order{}

	for _, order := range od.orders {
		if order.CustomerId == customerId {
			orders = append(orders, order)
		}
	}

	return orders
}

func (od *OrderDatabaseImpl) GetByToken(orderId string) utils.Order {
	for _, order := range od.orders {
		if order.Id == orderId {
			return order
		}
	}

	panic("Failed to get orders for order Token: " + orderId)
}

func (od *OrderDatabaseImpl) Add(customerId int, order utils.Order) utils.Order {
	order.UpdatedDate = time.Now().String()
	order.CustomerId = customerId
	order.Id = utils.UrlSafeUniqueId()

	od.orders = append(
		od.orders,
		order,
	)

	return order
}

type OrderDatabaseImpl struct {
	orders []utils.Order
}
