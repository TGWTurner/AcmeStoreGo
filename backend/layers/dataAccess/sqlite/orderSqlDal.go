package sqlite

import (
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
	"errors"
	"reflect"
	"strconv"

	"gorm.io/gorm"
)

func NewOrderDatabase(db *gorm.DB) *OrderDatabaseImpl {
	od := OrderDatabaseImpl{
		db: db,
	}

	testOrders := ConvertToDbOrders(testData.GetOrderTestData())

	if res := db.Create(&testOrders); res.Error != nil {
		panic("Failed to create test orders")
	}

	return &od
}

func (ad OrderDatabaseImpl) Close() {
	db, err := ad.db.DB()

	if err != nil {
		panic("Failed to get order db instance")
	}

	db.Close()
}

func (od *OrderDatabaseImpl) getOrderItemsFromOrderPk(pk int) ([]utils.OrderItem, error) {
	orderItems := []OrderItem{}

	response := od.db.Where("order_id = ?", pk).Find(&orderItems)

	if response.Error != nil {
		return []utils.OrderItem{}, errors.New("Failed to get order items for order with pk: " + strconv.Itoa(pk) + ", error: " + response.Error.Error())
	}

	return ConvertFromDbOrderItems(orderItems), nil
}

func (od *OrderDatabaseImpl) GetByCustomerId(customerId int) ([]utils.Order, error) {
	var orders []Order

	response := od.db.Where("customer_id = ?", customerId).Find(&orders)

	if response.Error != nil {
		return []utils.Order{}, errors.New("Failed to get orders for customerId: " + strconv.Itoa(customerId) + ", error: " + response.Error.Error())
	}

	customerOrders := make([]utils.Order, len(orders))

	for i, order := range orders {
		customerOrders[i] = ConvertFromDbOrder(order)
		items, err := od.getOrderItemsFromOrderPk(order.Pk)

		if err != nil {
			return []utils.Order{}, err
		}

		customerOrders[i].Items = items
	}

	return customerOrders, nil
}

func (od *OrderDatabaseImpl) GetByToken(orderToken string) (utils.Order, error) {
	var order Order

	response := od.db.Where("id = ?", orderToken).Limit(1).Find(&order)

	if response.Error != nil {
		return utils.Order{}, errors.New("Failed to get order for orderToken: " + orderToken + ", error: " + response.Error.Error())
	}

	if reflect.DeepEqual(order, Order{}) {
		return utils.Order{}, errors.New("Order does not exist with orderToken: " + orderToken)
	}

	customerOrder := ConvertFromDbOrder(order)
	items, err := od.getOrderItemsFromOrderPk(order.Pk)

	if err != nil {
		return utils.Order{}, err
	}

	customerOrder.Items = items

	return customerOrder, nil
}

func (od *OrderDatabaseImpl) Add(customerId int, order utils.Order) (utils.Order, error) {
	dbOrder := ConvertToDbOrder(order)
	dbOrder.SetUpNewOrder(customerId)

	response := od.db.Create(&dbOrder)

	if response.Error != nil {
		return utils.Order{}, errors.New("Failed to create new Order, error: " + response.Error.Error())
	}

	dbOrderItems := ConvertToDbOrderItems(dbOrder.Pk, order)

	for _, item := range dbOrderItems {
		response := od.db.Create(&item)

		if response.Error != nil {
			return utils.Order{}, errors.New("Failed to create item entry for pk: " + strconv.Itoa(item.ProductId) + ", error: " + response.Error.Error())
		}
	}

	return od.GetByToken(dbOrder.Id)
}

type OrderDatabaseImpl struct {
	db *gorm.DB
}
