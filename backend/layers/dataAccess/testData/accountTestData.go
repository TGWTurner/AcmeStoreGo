package testData

import (
	"bjssStoreGo/backend/utils"
)

func (atd AccountTestData) GetTestData() []utils.Account {
	return []utils.Account{
		{
			Id: 1,
			ShippingDetails: utils.ShippingDetails{
				Email:    "pre-populated-test-account@example.com",
				Name:     "Pre-populated Test Account",
				Address:  "123 Pre-Populated St, Test Town, Example",
				Postcode: "PL7 1RF",
			},
			PasswordHash: "flN+ZCd1vfgwose87/0dA0suvPYnDnuNxYRzWAmouLFsR/LV34DSev7BF3jc5M8uUFvXI4idZUumQ5jo/FPmnA==:Zedgm/FUGr7aJQWC9260kw==",
		},
		{
			Id: 2,
			ShippingDetails: utils.ShippingDetails{
				Email:    "pre-populated-test-account2@example.com",
				Name:     "Pre-populated Test Account 2",
				Address:  "32 Pre-Populated St, Test Town, Example",
				Postcode: "PL7 1RF",
			},
			PasswordHash: "flN+ZCd1vfgwose87/0dA0suvPYnDnuNxYRzWAmouLFsR/LV34DSev7BF3jc5M8uUFvXI4idZUumQ5jo/FPmnA==:Zedgm/FUGr7aJQWC9260kw==",
		},
	}
}

/*
TODO: Question - Cant have defaults so variadic, test for existance?
*/
func (atd AccountTestData) GetTestAccountCredentials(indexes ...int) struct {
	Email    string
	Password string
} {
	index := 0

	if indexes != nil {
		index = indexes[0]
	}

	testAccount := AccountTestData{}.GetTestData()[index]

	return struct {
		Email    string
		Password string
	}{
		Email:    testAccount.Email,
		Password: "I am an insecure password",
	}
}

type AccountTestData struct{}
