package sqlite

import (
	"bjssStoreGo/backend/utils"
	"time"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Id           int    `gorm:"primaryKey"`
	Email        string `gorm:"not null; unique"`
	Name         string `gorm:"not null"`
	Address      string `gorm:"not null"`
	Postcode     string `gorm:"not null"`
	PasswordHash string `gorm:"not null"`
}

func ConvertToDbAccount(account utils.Account) *Account {
	return &Account{
		Id:           account.Id,
		Email:        account.Email,
		Name:         account.Name,
		Address:      account.Address,
		Postcode:     account.Postcode,
		PasswordHash: account.PasswordHash,
	}
}

func (a *Account) ConvertFromDbAccount() *utils.Account {
	return &utils.Account{
		Id:           a.Id,
		PasswordHash: a.PasswordHash,
		ShippingDetails: utils.ShippingDetails{
			Email:    a.Email,
			Name:     a.Name,
			Address:  a.Address,
			Postcode: a.Postcode,
		},
	}
}

type Category struct {
	gorm.Model
	Id   int    `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

func ConvertToDbCategory(category utils.ProductCategory) *Category {
	return &Category{
		Id:   category.Id,
		Name: category.Name,
	}
}

func (c *Category) ConvertFromDbCategory() *utils.ProductCategory {
	return &utils.ProductCategory{
		Id:   c.Id,
		Name: c.Name,
	}
}

type Deal struct {
	gorm.Model
	ProductId int     `gorm:"not null"`
	StartDate string  `gorm:"not null"`
	EndDate   string  `gorm:"not null"`
	Product   Product `gorm:"ForeignKey:ProductId"`
}

func ConvertToDbDeal(deal utils.ProductDeal) *Deal {
	return &Deal{
		ProductId: deal.ProductId,
		StartDate: deal.StartDate,
		EndDate:   deal.EndDate,
	}
}

func (d *Deal) ConvertFromDbDeal() *utils.ProductDeal {
	return &utils.ProductDeal{
		ProductId: d.ProductId,
		StartDate: d.StartDate,
		EndDate:   d.EndDate,
	}
}

type OrderItem struct {
	gorm.Model
	OrderId   int     `gorm:"not null"`
	ProductId int     `gorm:"not null"`
	Quantity  int     `gorm:"not null"`
	Order     Order   `gorm:"ForeignKey:OrderId"`
	Product   Product `gorm:"ForeignKey:ProductId"`
}

func ConvertToDbOrderItems(orderId int, order utils.Order) []*OrderItem {
	orderItems := []*OrderItem{}

	for _, item := range order.Items {
		orderItems = append(orderItems, &OrderItem{
			OrderId:   orderId,
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	return orderItems
}

func (oi *OrderItem) ConvertFromDbOrderItem() utils.OrderItem {
	return utils.OrderItem{
		ProductId: oi.ProductId,
		Quantity:  oi.Quantity,
	}
}

type Order struct {
	gorm.Model
	Pk          int    `gorm:"primaryKey"`
	Id          string `gorm:"unique;index;not null"`
	CustomerId  int    `gorm:"index"`
	Total       int    `gorm:"not null"`
	UpdatedDate string `gorm:"not null"`
	Email       string `gorm:"not null; unique"`
	Name        string `gorm:"not null"`
	Address     string `gorm:"not null"`
	Postcode    string `gorm:"not null"`
}

func ConvertToDbOrder(order utils.Order) *Order {
	return &Order{
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

func (o *Order) SetUpNewOrder(customerId int) {
	o.Id = utils.UrlSafeUniqueId()
	o.CustomerId = customerId
	o.UpdatedDate = time.Now().String()
}

func (o *Order) ConvertFromDbOrder() utils.Order {
	return utils.Order{
		Id:          o.Id,
		Total:       o.Total,
		UpdatedDate: o.UpdatedDate,
		CustomerId:  o.CustomerId,
		ShippingDetails: utils.ShippingDetails{
			Email:    o.Email,
			Name:     o.Name,
			Address:  o.Address,
			Postcode: o.Postcode,
		},
	}
}

type Product struct {
	gorm.Model
	Id                int      `gorm:"primaryKey"`
	QuantityRemaining int      `gorm:"not null"`
	CategoryId        int      `gorm:"index;not null"`
	Price             int      `gorm:"not null"`
	ShortDescription  string   `gorm:"not null"`
	LongDescription   string   `gorm:"not null"`
	Category          Category `gorm:"ForeignKey:CategoryId"`
}

func ConvertToDbProduct(product utils.Product) *Product {
	return &Product{
		Id:                product.Id,
		QuantityRemaining: product.QuantityRemaining,
		CategoryId:        product.CategoryId,
		Price:             product.Price,
		ShortDescription:  product.ShortDescription,
		LongDescription:   product.LongDescription,
	}
}

func (p *Product) ConvertFromDbProduct() *utils.Product {
	return &utils.Product{
		Id:                p.Id,
		QuantityRemaining: p.QuantityRemaining,
		CategoryId:        p.CategoryId,
		Price:             p.Price,
		ShortDescription:  p.ShortDescription,
		LongDescription:   p.LongDescription,
	}
}
