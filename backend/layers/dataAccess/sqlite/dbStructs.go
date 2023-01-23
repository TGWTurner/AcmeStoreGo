package sqlite

import "gorm.io/gorm"

//this will contain versions of the utils structs and conversion methods from utils to db version
//these will be modelled in the exact way we want the db to behave
type Account struct {
	gorm.Model
	Id           int    `gorm:"primaryKey"`
	Email        string `gorm:"not null; unique"`
	Address      string `gorm:"not null"`
	Postcode     string `gorm:"not null"`
	PasswordHash string `gorm:"not null"`
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

type OrderItem struct {
	gorm.Model
	OrderId   int     `gorm:"not null"`
	ProductId int     `gorm:"not null"`
	Quantity  int     `gorm:"not null"`
	Order     Order   `gorm:"ForeignKey:OrderId"`
	Product   Product `gorm:"ForeignKey:ProductId"`
}

type Order struct {
	gorm.Model
	pk          int    `gorm:"primaryKey"`
	Id          string `gorm:"unique;index;not null;"`
	CustomerId  int    `gorm:"index"`
	Total       int    `gorm:"not null"`
	UpdatedDate string `gorm:"not null"`
	Email       string `gorm:"not null; unique"`
	Name        string `gorm:"not null"`
	Address     string `gorm:"not null"`
	Postcode    string `gorm:"not null"`
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
