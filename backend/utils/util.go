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
