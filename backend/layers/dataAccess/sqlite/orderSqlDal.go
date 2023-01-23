package sqlite

import (
	"bjssStoreGo/backend/utils"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func NewOrderDatabase(db *gorm.DB) OrderDatabase {
	od := OrderDatabase{
		db: db,
	}

	return od
}

func (od *OrderDatabase) GetOrdersByCustomerId(customerId int) []utils.Order {
	var orders []utils.Order

	//Need to add the Order columns too
	response := od.db.Model(&orders).
		Select("id, customerId, total, updatedDate, email, name, address, postcode, Order.productId, Order.quantity").
		Joins("INNER JOIN Order ON Order.orderId = orders.pk").
		Where("orders.customerId = ?", customerId)

	if response.Error != nil {
		panic("Failed to get orders for customerId: " + strconv.Itoa(customerId))
	}

	return orders
}

func (od *OrderDatabase) GetOrderByToken(orderToken int) utils.Order {
	var order utils.Order

	//Need to add the Order columns too
	response := od.db.Model(&order).
		Select("id, customerId, total, updatedDate, email, name, address, postcode, order_items.productId, order_items.quantity").
		Joins("INNER JOIN order_items ON order_items.orderId = orders.pk").
		Where("orders.Id = ?", orderToken)

	if response.Error != nil {
		panic("Failed to get orders for order Token: " + strconv.Itoa(orderToken))
	}

	// orderItems := od.GetOrderItemsFromOrderId(order.ID)

	return order
}

// func (od *OrderDatabase) GetOrderItemsFromOrderId(orderId uint) []utils.OrderItem {
// 	var orderItems []utils.OrderItem

// 	response := od.db.Where("order_id = ?", orderId).Find(&orderItems)
// }

func (od *OrderDatabase) AddOrder(customerId int, order utils.Order) utils.Order {
	order.Id = utils.UrlSafeUniqueId()
	order.CustomerId = customerId
	order.UpdatedDate = time.Now().String()

	response := od.db.Create(&order)

	if response.Error != nil {
		panic("Failed to create new Order")
	}

	for _, item := range order.Items {
		//TODO: 1234 -Convert to use dbStructs version and conversions
		item.OrderId = int(order.Id)
		response := od.db.Create(&item)

		if response.Error != nil {
			panic("Failed to create item entry for item")
		}
	}

	return od.GetOrderByToken(int(order.Id))
}

type OrderDatabase struct {
	db *gorm.DB
}
