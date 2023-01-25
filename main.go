package main

import (
	"bjssStoreGo/backend/layers/dataAccess"
	"bjssStoreGo/backend/layers/dataAccess/sqlite"
	"bjssStoreGo/backend/utils"
	"fmt"
)

func main() {
	db := dataAccess.InitiateConnection()

	fmt.Println("Db initialised")
	fmt.Println(db.Account.GetByEmail("pre-populated-test-account@example.com"))

	order := db.Order.Add(
		5,
		utils.Order{
			Id:         utils.UrlSafeUniqueId(),
			CustomerId: 1,
			Total:      55,
			ShippingDetails: utils.ShippingDetails{
				Name:     "name",
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

	fmt.Println(order)

	/*
		TESTING:
		 - Products -
		 1. Get All Products
		 2. Get Product given an Id
		 3. Get Product given []Id's
		 4. Get All categories
		 5. Get Products given Id
		 6. Search Products "Dog"
		 7. Get Products with deals
		 8. Decrease stock

		  - Orders -
		 1. Create order
		 2. Get Order from customer Id
		 3. Get Order from token

		  - Accounts -
		 1. Add Account
		 2. Get Account by email
		 3. Get Account by id
		 4. Update an Account
	*/
}

func testCreateOrder(db sqlite.Database) {
	order := utils.Order{
		Total:      5,
		CustomerId: 2,
		ShippingDetails: utils.ShippingDetails{
			Email:    "testEmail",
			Name:     "testName",
			Address:  "testAddress",
			Postcode: "testPostcode",
		},
		Items: []utils.OrderItem{
			{
				ProductId: 1,
				Quantity:  10,
			},
			{
				ProductId: 5,
				Quantity:  1,
			},
		},
	}

	db.Order.Add(2, order)
}

func testGetOrderByToken(db sqlite.Database) {
	order := utils.Order{
		Total:      5,
		CustomerId: 2,
		ShippingDetails: utils.ShippingDetails{
			Email:    "testEmail",
			Name:     "testName",
			Address:  "testAddress",
			Postcode: "testPostcode",
		},
		Items: []utils.OrderItem{
			{
				ProductId: 1,
				Quantity:  10,
			},
			{
				ProductId: 5,
				Quantity:  1,
			},
		},
	}

	orderToken := db.Order.Add(2, order).Id

	dbOrder := db.Order.GetByToken(orderToken)

	if order.Total != dbOrder.Total ||
		order.CustomerId != dbOrder.CustomerId ||
		order.ShippingDetails.Address != dbOrder.ShippingDetails.Address ||
		order.ShippingDetails.Email != dbOrder.ShippingDetails.Email ||
		order.ShippingDetails.Name != dbOrder.ShippingDetails.Name ||
		order.Items[0] != dbOrder.Items[0] ||
		order.Items[1] != dbOrder.Items[1] {
		panic("testGetOrderByToken failed")
	}
}

func testGetAllProducts(db sqlite.Database) {
	//products := db.Product.GetAll()
}

func testGetProductGivenId(db sqlite.Database) {
	product := db.Product.GetByIds(1)[0]

	if product.Id != 1 {
		panic("testGetProductGivenId failed")
	}
}

func testGetProductsGivenIds(db sqlite.Database) {
	products := db.Product.GetByIds(1, 2, 3)

	for i, product := range products {
		if i != product.Id {
			panic("testGetProductsGivenIds failed")
		}
	}
}
