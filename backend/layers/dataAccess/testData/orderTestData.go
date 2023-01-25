package testData

import (
	"bjssStoreGo/backend/utils"
	"time"
)

func GetOrderTestData() []utils.Order {
	return []utils.Order{
		{
			Id:          utils.UrlSafeUniqueId(),
			Total:       1,
			UpdatedDate: time.Now().String(),
			CustomerId:  1,
			ShippingDetails: utils.ShippingDetails{
				Email:    "pre-populated-test-account@example.com",
				Name:     "Pre-populated Test Account",
				Address:  "123 Pre-Populated St, Test Town, Example",
				Postcode: "PL7 1RF",
			},
			Items: []utils.OrderItem{
				{
					ProductId: 1,
					Quantity:  10,
				},
				{
					ProductId: 5,
					Quantity:  2,
				},
			},
		},
		{
			Id:          utils.UrlSafeUniqueId(),
			Total:       1,
			UpdatedDate: time.Now().String(),
			CustomerId:  1,
			ShippingDetails: utils.ShippingDetails{
				Email:    "pre-populated-test-account@example.com",
				Name:     "Pre-populated Test Account",
				Address:  "123 Pre-Populated St, Test Town, Example",
				Postcode: "PL7 1RF",
			},
			Items: []utils.OrderItem{
				{
					ProductId: 4,
					Quantity:  3,
				},
				{
					ProductId: 5,
					Quantity:  2,
				},
			},
		},
		{
			Id:          utils.UrlSafeUniqueId(),
			Total:       1,
			UpdatedDate: time.Now().String(),
			CustomerId:  1,
			ShippingDetails: utils.ShippingDetails{
				Email:    "pre-populated-test-account@example.com",
				Name:     "Pre-populated Test Account",
				Address:  "123 Pre-Populated St, Test Town, Example",
				Postcode: "PL7 1RF",
			},
			Items: []utils.OrderItem{
				{
					ProductId: 1,
					Quantity:  6,
				},
				{
					ProductId: 6,
					Quantity:  4,
				},
			},
		},
	}
}
