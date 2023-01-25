package sqlite

import (
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
	"strconv"

	"gorm.io/gorm"
)

func NewOrderDatabase(db *gorm.DB) OrderDatabase {
	od := OrderDatabase{
		db: db,
	}

	testOrders := ConvertToDbOrders(testData.GetOrderTestData())

	if res := db.Create(&testOrders); res.Error != nil {
		panic("Failed to create test orders")
	}

	return od
}

func (ad OrderDatabase) Close() {
	db, _ := ad.db.DB()
	db.Close()
}

func (od *OrderDatabase) getOrderItemsFromOrderPk(pk int) []utils.OrderItem {
	orderItems := []OrderItem{}

	response := od.db.Where("order_id = ?", pk).Find(&orderItems)

	if response.Error != nil {
		panic("Failed to get order items for order with pk: " + strconv.Itoa(pk))
	}

	customerOrderItems := make([]utils.OrderItem, len(orderItems))

	for i, orderItem := range orderItems {
		customerOrderItems[i] = *orderItem.ConvertFromDbOrderItem()
	}

	return customerOrderItems
}

func (od *OrderDatabase) GetByCustomerId(customerId int) []utils.Order {
	var orders []Order

	response := od.db.Where("customer_id = ?", customerId).Find(&orders)

	if response.Error != nil {
		panic("Failed to get orders for customerId: " + strconv.Itoa(customerId))
	}

	customerOrders := make([]utils.Order, len(orders))

	for i, order := range orders {
		customerOrders[i] = order.ConvertFromDbOrder()
		customerOrders[i].Items = od.getOrderItemsFromOrderPk(order.Pk)
	}

	return customerOrders
}

func (od *OrderDatabase) GetByToken(orderToken string) utils.Order {
	var order Order

	response := od.db.Where("id = ?", orderToken).First(&order)

	if response.Error != nil {
		panic("Failed to get order for orderToken: " + orderToken)
	}

	customerOrder := order.ConvertFromDbOrder()
	customerOrder.Items = od.getOrderItemsFromOrderPk(order.Pk)

	return customerOrder
}

func (od *OrderDatabase) Add(customerId int, order utils.Order) utils.Order {
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

type OrderDatabase struct {
	db *gorm.DB
}
