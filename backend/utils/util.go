package utils

import (
	"crypto/rand"
	"encoding/base64"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Id                int      `gorm:"primaryKey"`
	CategoryId        int      `gorm:"index;not null"`
	Price             int      `gorm:"not null"`
	QuantityRemaining int      `gorm:"not null"`
	ShortDescription  string   `gorm:"not null"`
	LongDescription   string   `gorm:"not null"`
	Category          Category `gorm:"ForeignKey:CategoryId"`
}

type Category struct {
	gorm.Model
	Id   int    `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

type Deal struct {
	gorm.Model
	ProductId int     `gorm:"not null"`
	StartDate string  `gorm:"not null"`
	EndDate   string  `gorm:"not null"`
	Product   Product `gorm:"ForeignKey:ProductId"`
}

type Order struct {
	gorm.Model
	pk              int    `gorm:"primaryKey"`
	Id              string `gorm:"unique;index;not null;"`
	CustomerId      int    `gorm:"index"`
	Total           int    `gorm:"not null"`
	UpdatedDate     string `gorm:"not null"`
	ShippingDetails ShippingDetails
	Items           []OrderItem
}

type OrderItem struct {
	gorm.Model
	ProductId int     `gorm:"not null"`
	Quantity  int     `gorm:"not null"`
	Order     Order   `gorm:"ForeignKey:OrderId"`
	Product   Product `gorm:"ForeignKey:ProductId"`
}

type ShippingDetails struct {
	Email    string `gorm:"unique;index;not null"`
	Name     string `gorm:"not null"`
	Address  string `gorm:"not null"`
	Postcode string `gorm:"not null"`
}

type Account struct {
	gorm.Model
	Id           int    `gorm:"primaryKey"`
	PasswordHash string `gorm:"not null"`
	ShippingDetails
}

func UrlSafeUniqueId() string {
	random128bitNumber := make([]byte, 128)

	_, err := rand.Read(random128bitNumber)

	if err != nil {
		panic("Failed to generate random string: " + err.Error())
	}

	return base64.RawURLEncoding.EncodeToString(random128bitNumber)
}
