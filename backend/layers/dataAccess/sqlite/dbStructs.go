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

func ConvertToDbAccounts(accounts []utils.Account) []*Account {
	dbAccounts := make([]*Account, len(accounts))

	for i, account := range accounts {
		dbAccounts[i] = ConvertToDbAccount(account)
	}

	return dbAccounts
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

func ConvertFromDbAccounts(dbAccounts []Account) []*utils.Account {
	accounts := make([]*utils.Account, len(dbAccounts))

	for i, account := range dbAccounts {
		accounts[i] = account.ConvertFromDbAccount()
	}

	return accounts
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

func ConvertToDbCategories(categories []utils.ProductCategory) []*Category {
	dbCategories := make([]*Category, len(categories))

	for i, category := range categories {
		dbCategories[i] = ConvertToDbCategory(category)
	}

	return dbCategories
}

func (c *Category) ConvertFromDbCategory() *utils.ProductCategory {
	return &utils.ProductCategory{
		Id:   c.Id,
		Name: c.Name,
	}
}

func ConvertFromDbCategories(dbCategories []Category) []*utils.ProductCategory {
	categories := make([]*utils.ProductCategory, len(dbCategories))

	for i, category := range dbCategories {
		categories[i] = category.ConvertFromDbCategory()
	}

	return categories
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

func ConvertToDbDeals(deals []utils.ProductDeal) []*Deal {
	dbDeals := make([]*Deal, len(deals))

	for i, deal := range deals {
		dbDeals[i] = ConvertToDbDeal(deal)
	}

	return dbDeals
}

func (d *Deal) ConvertFromDbDeal() *utils.ProductDeal {
	return &utils.ProductDeal{
		ProductId: d.ProductId,
		StartDate: d.StartDate,
		EndDate:   d.EndDate,
	}
}

func ConvertFromDbDeals(dbDeals []Deal) []*utils.ProductDeal {
	deals := make([]*utils.ProductDeal, len(dbDeals))

	for i, account := range dbDeals {
		deals[i] = account.ConvertFromDbDeal()
	}

	return deals
}

type OrderItem struct {
	gorm.Model
	OrderId   int     `gorm:"not null"`
	ProductId int     `gorm:"not null"`
	Quantity  int     `gorm:"not null"`
	Order     Order   `gorm:"ForeignKey:OrderId"`
	Product   Product `gorm:"ForeignKey:ProductId"`
}

func ConvertToDbOrderItem(orderId int, orderItem utils.OrderItem) *OrderItem {
	return &OrderItem{
		OrderId:   orderId,
		ProductId: orderItem.ProductId,
		Quantity:  orderItem.Quantity,
	}
}

func ConvertToDbOrderItems(orderId int, order utils.Order) []*OrderItem {
	orderItems := make([]*OrderItem, len(order.Items))

	for i, item := range order.Items {
		orderItems[i] = ConvertToDbOrderItem(orderId, item)
	}

	return orderItems
}

func (oi *OrderItem) ConvertFromDbOrderItem() *utils.OrderItem {
	return &utils.OrderItem{
		ProductId: oi.ProductId,
		Quantity:  oi.Quantity,
	}
}

func ConvertFromDbOrderItems(dbOrderItems []OrderItem) []*utils.OrderItem {
	orderItems := make([]*utils.OrderItem, len(dbOrderItems))

	for i, orderItem := range dbOrderItems {
		orderItems[i] = orderItem.ConvertFromDbOrderItem()
	}

	return orderItems
}

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

func ConvertToDbOrders(orders []utils.Order) []*Order {
	dbOrders := make([]*Order, len(orders))

	for i, order := range orders {
		dbOrders[i] = ConvertToDbOrder(order)
	}

	return dbOrders
}

func (o *Order) SetUpNewOrder(customerId int) {
	o.Id = utils.UrlSafeUniqueId()
	o.CustomerId = customerId
	o.UpdatedDate = time.Now().String()
}

func (o *Order) ConvertFromDbOrder() *utils.Order {
	return &utils.Order{
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

func ConvertFromDbOrder(dbOrders []Order) []*utils.Order { //YRKU19
	orders := make([]*utils.Order, len(dbOrders))

	for i, order := range dbOrders {
		orders[i] = order.ConvertFromDbOrder()
	}

	return orders
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

func ConvertToDbProducts(products []utils.Product) []*Product {
	dbProduct := make([]*Product, len(products))

	for i, product := range products {
		dbProduct[i] = ConvertToDbProduct(product)
	}

	return dbProduct
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
