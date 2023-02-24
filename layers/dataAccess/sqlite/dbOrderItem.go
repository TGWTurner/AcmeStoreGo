package sqlite

import (
	"backend/utils"

	"gorm.io/gorm"
)

type OrderItem struct {
	gorm.Model
	OrderId   int     `gorm:"not null"`
	ProductId int     `gorm:"not null"`
	Quantity  int     `gorm:"not null"`
	Order     Order   `gorm:"ForeignKey:OrderId"`
	Product   Product `gorm:"ForeignKey:ProductId"`
}

func ConvertToDbOrderItem(orderId int, orderItem utils.OrderItem) OrderItem {
	return OrderItem{
		OrderId:   orderId,
		ProductId: orderItem.ProductId,
		Quantity:  orderItem.Quantity,
	}
}

func ConvertToDbOrderItems(orderId int, order utils.Order) []OrderItem {
	orderItems := make([]OrderItem, len(order.Items))

	for i, item := range order.Items {
		orderItems[i] = ConvertToDbOrderItem(orderId, item)
	}

	return orderItems
}

func ConvertFromDbOrderItem(orderItem OrderItem) utils.OrderItem {
	return utils.OrderItem{
		ProductId: orderItem.ProductId,
		Quantity:  orderItem.Quantity,
	}
}

func ConvertFromDbOrderItems(dbOrderItems []OrderItem) []utils.OrderItem {
	orderItems := make([]utils.OrderItem, len(dbOrderItems))

	for i, orderItem := range dbOrderItems {
		orderItems[i] = ConvertFromDbOrderItem(orderItem)
	}

	return orderItems
}
