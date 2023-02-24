package sqlite

import (
	"backend/utils"
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	Pk          int    `gorm:"primaryKey"`
	Id          string `gorm:"unique;index;not null"`
	CustomerId  int    `gorm:"index"`
	Total       int    `gorm:"not null"`
	UpdatedDate string `gorm:"not null"`
	Email       string `gorm:"not null"`
	Name        string `gorm:"not null"`
	Address     string `gorm:"not null"`
	Postcode    string `gorm:"not null"`
}

func ConvertToDbOrder(order utils.Order) Order {
	return Order{
		Id:          order.Id,
		CustomerId:  order.CustomerId,
		Total:       order.Total,
		UpdatedDate: order.UpdatedDate,
		Email:       order.ShippingDetails.Email,
		Name:        order.ShippingDetails.Name,
		Address:     order.ShippingDetails.Address,
		Postcode:    order.ShippingDetails.Postcode,
	}
}

func ConvertToDbOrders(orders []utils.Order) []Order {
	dbOrders := make([]Order, len(orders))

	for i, order := range orders {
		dbOrders[i] = ConvertToDbOrder(order)
	}

	return dbOrders
}

func ConvertFromDbOrder(order Order) utils.Order {
	return utils.Order{
		Id:          order.Id,
		Total:       order.Total,
		UpdatedDate: order.UpdatedDate,
		CustomerId:  order.CustomerId,
		ShippingDetails: utils.ShippingDetails{
			Email:    order.Email,
			Name:     order.Name,
			Address:  order.Address,
			Postcode: order.Postcode,
		},
	}
}

func ConvertFromDbOrders(dbOrders []Order) []utils.Order {
	orders := make([]utils.Order, len(dbOrders))

	for i, order := range dbOrders {
		orders[i] = ConvertFromDbOrder(order)
	}

	return orders
}

func (o *Order) SetUpNewOrder(customerId int) {
	o.Id = utils.UrlSafeUniqueId()
	o.CustomerId = customerId
	o.UpdatedDate = time.Now().String()
}
