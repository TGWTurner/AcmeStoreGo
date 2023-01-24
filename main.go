package main

import (
	"bjssStoreGo/backend/layers/dataAccess"
	"bjssStoreGo/backend/utils"
)

func main() {
	db := dataAccess.DataAccess{}.InitiateConnection()

	// db.Account.Add(utils.Account{
	// 	Id:           0,
	// 	Email:        "test@test.com",
	// 	Name:         "test",
	// 	Address:      "test Address",
	// 	Postcode:     "TE57 7ST",
	// 	PasswordHash: "This Is A Password Hash",
	// })

	// fmt.Println(db.Account.GetByEmail("test@test.com"))

	db.Order.AddOrder(
		5,
		utils.Order{
			Id:         utils.UrlSafeUniqueId(),
			CustomerId: 1,
			Total:      55,
			ShippingDetails: utils.ShippingDetails{
				Email:    "Email",
				Address:  "Address",
				Postcode: "Postcode",
			},
			Items: []utils.OrderItem{
				{
					ProductId: 1,
					Quantity:  5,
				},
				{
					ProductId: 2,
					Quantity:  10,
				},
			},
		},
	)

	//Create an order, add it to the db, try out the gorm request then figure out how to create order object from the db
}
