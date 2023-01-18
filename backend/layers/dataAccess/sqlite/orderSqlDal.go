package sqlite

import (
	"bjssStoreGo/backend/utils"
	"strconv"

	"gorm.io/gorm"
)

func NewOrderDatabase(db *gorm.DB) OrderDatabase {
	od := OrderDatabase{
		db: db,
	}

	return od
}

func (od OrderDatabase) getOrdersByCustomerId(customerId int) []utils.Order {
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

func (od OrderDatabase) getOrderByToken(orderToken int) utils.Order {
	var order utils.Order

	//Need to add the Order columns too
	response := od.db.Model(&order).
		Select("id, customerId, total, updatedDate, email, name, address, postcode, Order.productId, Order.quantity").
		Joins("INNER JOIN Order ON Order.orderId = orders.pk").
		Where("orders.Id = ?", orderToken)

	if response.Error != nil {
		panic("Failed to get orders for customerId: " + strconv.Itoa(orderToken))
	}

	return order
}

func (od OrderDatabase) addOrder(customerId int, order utils.Order) {
	response := od.db.Create(&order)

	if response.Error != nil {
		panic("Failed to create new Order")
	}

	//TODO: Rearrange order object to include way to store items
	items := []utils.OrderItem{
		{
			ProductId: 1,
			Quantity:  5,
		},
		{
			ProductId: 2,
			Quantity:  10,
		},
		{
			ProductId: 3,
			Quantity:  1,
		},
		{
			ProductId: 3,
			Quantity:  3,
		},
	}

	for _, item := range items {
		item.OrderId = int(order.ID)
		response := od.db.Create(&item)

		if response.Error != nil {
			panic("Failed to create item entry for item")
		}
	}

	//TODO: Implement get order by token
	/*
		SQL:
		Create record of order
		INSERT INTO orders (
			id,
			customerId,
			total,
			updatedDate,
			email,
			name,
			address,
			postcode
		) VALUES (?)
		Create record of order items
		INSERT INTO Order (orderId, productId, quantity) VALUES (?)
	*/
}

type OrderDatabase struct {
	db *gorm.DB
}
