package test

import (
	"bjssStoreGo/backend/layers/businessLogic"
	"bjssStoreGo/backend/layers/dataAccess"
	"bjssStoreGo/backend/utils"
	"testing"
)

var signUpData = struct {
	email    string
	name     string
	address  string
	postcode string
	password string
}{
	email:    "a@example.com",
	name:     "a",
	address:  "addr",
	postcode: "abc123",
	password: "s3cret",
}

func setUpAccount() businessLogic.AccountService {
	db := dataAccess.InitiateConnection()
	return *businessLogic.NewAccountService(db.Account)
}

func createAccount(t *testing.T, as businessLogic.AccountService) utils.AccountApiResponse {
	acc, err := as.SignUp(signUpData.email, signUpData.password, signUpData.name, signUpData.address, signUpData.postcode)

	AssertNil(t, err)

	return acc
}

func TestSignsUpOk(t *testing.T) {
	as := setUpAccount()
	defer as.Close()

	account, err := as.SignUp(signUpData.email, signUpData.password, signUpData.name, signUpData.address, signUpData.postcode)

	AssertNil(t, err)

	if account.Id == 0 {
		t.Errorf("Account Id was undefined")
	}

	AssertAccountDataMatchesSignUp(t, account)
}

func TestSignsInOk(t *testing.T) {
	as := setUpAccount()
	defer as.Close()

	createAccount(t, as)

	account, err := as.SignIn(signUpData.email, signUpData.password)

	AssertNil(t, err)

	AssertAccountDataMatchesSignUp(t, account)
}

func TestFailsSignIn(t *testing.T) {
	as := setUpAccount()
	defer as.Close()

	createAccount(t, as)

	_, err := as.SignIn(signUpData.email, "not password")

	AssertErrorString(t, err, "invalidPassword")

	_, err = as.SignIn("invalid@example.com", signUpData.password)

	AssertErrorString(t, err, "invalidEmail")
}

func TestFetchesAnAccount(t *testing.T) {
	as := setUpAccount()
	defer as.Close()

	account := createAccount(t, as)
	fetched, err := as.GetById(account.Id)

	AssertNil(t, err)

	AssertAccountDataMatching(t, fetched, account)
}

func TestUpdatesAnAccountIncludingPasswordAndCanSignIn(t *testing.T) {
	as := setUpAccount()
	defer as.Close()

	account := createAccount(t, as)

	account.Name = "b"

	updateAccount := utils.UpdateAccount{
		ShippingDetails: account.ShippingDetails,
		Password:        "newpass",
	}

	updateAccount.Id = 1

	updated, err := as.Update(updateAccount)

	AssertNil(t, err)

	AssertAccountDataMatching(t, account, updated)

	_, err = as.SignIn(updateAccount.Email, updateAccount.Password)

	AssertNil(t, err)
}

func AssertAccountDataMatching(t *testing.T, actual, expected utils.AccountApiResponse) {
	if expected.Address != actual.Address {
		t.Errorf("Got address: %s, wanted address: %s", expected.Address, actual.Address)
	}

	if expected.Email != actual.Email {
		t.Errorf("Got email: %s, wanted email: %s", expected.Email, actual.Email)
	}

	if expected.Name != actual.Name {
		t.Errorf("Got name: %s, wanted name: %s", expected.Name, actual.Name)
	}

	if expected.Postcode != actual.Postcode {
		t.Errorf("Got postcode: %s, wanted postcode: %s", expected.Postcode, actual.Postcode)
	}
}

func AssertAccountDataMatchesSignUp(t *testing.T, account utils.AccountApiResponse) {
	if account.Address != signUpData.address {
		t.Errorf("Got address: %s, wanted address: %s", account.Address, signUpData.address)
	}

	if account.Email != signUpData.email {
		t.Errorf("Got email: %s, wanted email: %s", account.Email, signUpData.email)
	}

	if account.Name != signUpData.name {
		t.Errorf("Got name: %s, wanted name: %s", account.Name, signUpData.name)
	}

	if account.Postcode != signUpData.postcode {
		t.Errorf("Got postcode: %s, wanted postcode: %s", account.Postcode, signUpData.postcode)
	}
}
