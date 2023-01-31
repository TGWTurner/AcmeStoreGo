package memory

import (
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
	"errors"
	"time"
)

func NewOrderDatabase() *OrderDatabaseImpl {
	testOrders := testData.GetOrderTestData()

	return &OrderDatabaseImpl{
		orders: testOrders,
	}
}

func (ad *OrderDatabaseImpl) Close() {}

func (od *OrderDatabaseImpl) GetByCustomerId(customerId int) ([]utils.Order, error) {
	orders := []utils.Order{}

	for _, order := range od.orders {
		if order.CustomerId == customerId {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

func (od *OrderDatabaseImpl) GetByToken(orderId string) (utils.Order, error) {
	for _, order := range od.orders {
		if order.Id == orderId {
			return order, nil
		}
	}

	return utils.Order{}, errors.New("Order does not exist with orderToken: " + orderId)
}

func (od *OrderDatabaseImpl) Add(customerId int, order utils.Order) (utils.Order, error) {
	order.UpdatedDate = time.Now().String()
	order.CustomerId = customerId
	order.Id = utils.UrlSafeUniqueId()

	od.orders = append(
		od.orders,
		order,
	)

	return order, nil
}

type OrderDatabaseImpl struct {
	orders []utils.Order
}
