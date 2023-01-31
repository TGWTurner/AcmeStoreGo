package api

import (
	"bjssStoreGo/backend/layers/businessLogic"
)

// TODO: Implement Get signed in user id
func getSignedInUserId(req string) string {
	return "userId"
}

// TODO: Implement Set signed in user id
func setSignedInUserId(req string, customerId string) {
	customerId = "userId"
}

func NewAccountApi(accountService *businessLogic.AccountService) *AccountApi {
	return &AccountApi{
		as: *accountService,
	}
}

func (a AccountApi) postSignIn() {
	//TODO: Implement post sign in
}

func (a AccountApi) postSignUp() {
	//TODO: Implement post sign up
}

func (a AccountApi) getAccount() {
	//TODO: Implement get account
}

func (a AccountApi) postAccount() {
	//TODO: Implement post account
}

type AccountApi struct {
	as businessLogic.AccountService
}
