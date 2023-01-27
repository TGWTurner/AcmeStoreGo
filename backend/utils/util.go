package utils

import (
	"crypto/rand"
	"encoding/base64"
)

/*
Missing:
 - Basket
 - Session
 - AccountApiResponse
*/

type ShippingDetails struct {
	Email    string
	Name     string
	Address  string
	Postcode string
}

type Account struct {
	Id           int
	PasswordHash string
	ShippingDetails
}

type Product struct {
	Id                int
	QuantityRemaining int
	CategoryId        int
	Price             int
	ShortDescription  string
	LongDescription   string
}

type ProductCategory struct {
	Id   int
	Name string
}

type ProductDeal struct {
	ProductId int
	StartDate string
	EndDate   string
}

type OrderItem struct {
	ProductId int
	Quantity  int
}

type Order struct {
	Id              string
	Total           int
	UpdatedDate     string
	CustomerId      int
	ShippingDetails ShippingDetails
	Items           []OrderItem
}

func UrlSafeUniqueId() string {
	random128bitNumber := make([]byte, 128)

	_, err := rand.Read(random128bitNumber)

	if err != nil {
		panic("Failed to generate random string: " + err.Error())
	}

	return base64.RawURLEncoding.EncodeToString(random128bitNumber)
}

type AccountDatabase interface {
	Close()
	Add(account Account) Account
	GetByEmail(email string) Account
	GetById(accountId int) Account
	Update(updateAccount Account) Account
}

type OrderDatabase interface {
	Close()
	GetByCustomerId(customerId int) []Order
	GetByToken(orderId string) Order
	Add(customerId int, order Order) Order
}

type ProductDatabase interface {
	Close()
	GetByIds(Ids ...int) []Product
	GetCategories() []ProductCategory
	GetByCategory(categoryId int) []Product
	GetByText(searchTerm string) []Product
	GetWithCurrentDeals(date string) []Product
	DecreaseStock(productQuantities []OrderItem)
}

type Database struct {
	Account AccountDatabase
	Product ProductDatabase
	Order   OrderDatabase
}
