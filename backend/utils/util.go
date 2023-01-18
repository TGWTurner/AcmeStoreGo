package utils

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Id                int      `gorm:"primaryKey"`
	CategoryId        int      `gorm:"index"`
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
	ProductId int     `gorm:"primaryKey"`
	StartDate string  `gorm:"not null"`
	EndDate   string  `gorm:"not null"`
	Product   Product `gorm:"ForeignKey:ProductId"`
}

type Order struct {
	gorm.Model
	pk          int    `gorm:"primaryKey"`
	Id          int    `gorm:"unique;index;not null;"`
	CustomerId  int    `gorm:"index"`
	Total       int    `gorm:"not null"`
	UpdatedDate string `gorm:"not null"`
	Email       string `gorm:"not null"`
	Name        string `gorm:"not null"`
	Address     string `gorm:"not null"`
	Postcode    string `gorm:"not null"`
}

type OrderItem struct {
	gorm.Model
	OrderId   int
	ProductId int     `gorm:"not null"`
	Quantity  int     `gorm:"not null"`
	Order     Order   `gorm:"ForeignKey:OrderId"`
	Product   Product `gorm:"ForeignKey:ProductId"`
}

type Account struct {
	gorm.Model
	Id           int    `gorm:"primaryKey"`
	Email        string `gorm:"unique;index, not null"`
	Name         string `gorm:"not null"`
	Address      string `gorm:"not null"`
	Postcode     string `gorm:"not null"`
	PasswordHash string `gorm:"not null"`
}
