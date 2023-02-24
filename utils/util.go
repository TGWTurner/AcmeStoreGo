package utils

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

type ShippingDetails struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Postcode string `json:"postcode"`
}

type Account struct {
	Id           int
	PasswordHash string
	ShippingDetails
}

type AccountApiResponse struct {
	Id int `json:"id"`
	ShippingDetails
}

type UpdateAccount struct {
	Id       int
	Password string
	ShippingDetails
}

func (a *Account) OmitPasswordHash() AccountApiResponse {
	return AccountApiResponse{
		Id:              a.Id,
		ShippingDetails: a.ShippingDetails,
	}
}

type Product struct {
	Id                int    `json:"id"`
	QuantityRemaining int    `json:"quantityRemaining"`
	CategoryId        int    `json:"categoryId"`
	Price             int    `json:"price"`
	ShortDescription  string `json:"shortDescription"`
	LongDescription   string `json:"longDescription"`
}

type ProductCategory struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ProductDeal struct {
	ProductId int    `json:"productId"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type OrderItem struct {
	ProductId int `json:"productId"`
	Quantity  int `json:"quantity"`
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

func GetFormattedDate() string {
	return time.Now().Format("2006-01-02")
}

type AccountDatabase interface {
	Close()
	Add(account Account) (Account, error)
	GetByEmail(email string) (Account, error)
	GetById(accountId int) (Account, error)
	Update(updateAccount Account) (Account, error)
}

type OrderDatabase interface {
	Close()
	GetByCustomerId(customerId int) ([]Order, error)
	GetByToken(orderId string) (Order, error)
	Add(customerId int, order Order) (Order, error)
}

type ProductDatabase interface {
	Close()
	GetAll() ([]Product, error)
	GetById(Id int) (Product, error)
	GetByIds(Ids ...int) ([]Product, error)
	GetCategories() ([]ProductCategory, error)
	GetByCategory(categoryId int) ([]Product, error)
	GetByText(searchTerm string) ([]Product, error)
	GetWithCurrentDeals(date string) ([]Product, error)
	DecreaseStock(productQuantities []OrderItem) error
}

type Database struct {
	Account AccountDatabase
	Product ProductDatabase
	Order   OrderDatabase
}

type Basket struct {
	Total int         `json:"total"`
	Items []OrderItem `json:"items"`
}
