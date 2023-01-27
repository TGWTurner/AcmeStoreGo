package sqlite

import (
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
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

func (od *OrderDatabaseImpl) getOrderItemsFromOrderPk(pk int) []utils.OrderItem {
	orderItems := []OrderItem{}

	response := od.db.Where("order_id = ?", pk).Find(&orderItems)

	if response.Error != nil {
		panic("Failed to get order items for order with pk: " + strconv.Itoa(pk))
	}

	return ConvertFromDbOrderItems(orderItems)
}

func (od *OrderDatabaseImpl) GetByCustomerId(customerId int) []utils.Order {
	var orders []Order

	response := od.db.Where("customer_id = ?", customerId).Find(&orders)

	if response.Error != nil {
		panic("Failed to get orders for customerId: " + strconv.Itoa(customerId))
	}

	customerOrders := make([]utils.Order, len(orders))

	for i, order := range orders {
		customerOrders[i] = ConvertFromDbOrder(order)
		customerOrders[i].Items = od.getOrderItemsFromOrderPk(order.Pk)
	}

	return customerOrders
}

func (od *OrderDatabaseImpl) GetByToken(orderToken string) utils.Order {
	var order Order

	response := od.db.Where("id = ?", orderToken).First(&order)

	if response.Error != nil {
		panic("Failed to get order for orderToken: " + orderToken)
	}

	customerOrder := ConvertFromDbOrder(order)
	customerOrder.Items = od.getOrderItemsFromOrderPk(order.Pk)

	return customerOrder
}

func (od *OrderDatabaseImpl) Add(customerId int, order utils.Order) utils.Order {
	dbOrder := ConvertToDbOrder(order)
	dbOrder.SetUpNewOrder(customerId)

	response := od.db.Create(&dbOrder)

	if response.Error != nil {
		panic("Failed to create new Order")
	}

	dbOrderItems := ConvertToDbOrderItems(dbOrder.Pk, order)

	for _, item := range dbOrderItems {
		response := od.db.Create(&item)

		if response.Error != nil {
			panic("Failed to create item entry for pk: " + strconv.Itoa(item.ProductId))
		}
	}

	return od.GetByToken(dbOrder.Id)
}

type OrderDatabaseImpl struct {
	db *gorm.DB
}
