package main

import (
	"bjssStoreGo/backend/utils"
	"fmt"
)

func main() {
	// db := dataAccess.DataAccess{}.InitiateConnection()

	// db.Account.Add(utils.Account{
	// 	Id:           0,
	// 	Email:        "test@test.com",
	// 	Name:         "test",
	// 	Address:      "test Address",
	// 	Postcode:     "TE57 7ST",
	// 	PasswordHash: "This Is A Password Hash",
	// })

	// fmt.Println(db.Account.GetByEmail("test@test.com"))

	// db.Order.AddOrder(
	// 	5,
	// 	utils.Order{
	// 		Id:         3,
	// 		CustomerId: 1,
	// 		Total:      55,
	// 		Email:      "Email",
	// 		Address:    "Address",
	// 		Postcode:   "Postcode",
	// 	},
	// )

	rand := utils.UrlSafeUniqueId()

	fmt.Println(rand)

}
