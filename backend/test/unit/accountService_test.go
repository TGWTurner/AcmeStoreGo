package unit

import (
	bl "bjssStoreGo/backend/layers/businessLogic"
	da "bjssStoreGo/backend/layers/dataAccess"
	"bjssStoreGo/backend/test"
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

func setUpAccount(t *testing.T) bl.AccountService {
	db := da.InitiateConnection()
	return *bl.NewAccountService(db.Account)
}

func createAccount(t *testing.T, as bl.AccountService) utils.AccountApiResponse {
	acc, err := as.SignUp(signUpData.email, signUpData.password, signUpData.name, signUpData.address, signUpData.postcode)

	test.AssertNil(t, err)

	return acc
}

func TestSignsUpOk(t *testing.T) {
	as := setUpAccount(t)
	defer as.Close()

	account, err := as.SignUp(signUpData.email, signUpData.password, signUpData.name, signUpData.address, signUpData.postcode)

	test.AssertNil(t, err)

	if account.Id == 0 {
		t.Errorf("Account Id was undefined")
	}

	AssertAccountDataMatchesSignUp(t, account)
}

func TestSignsInOk(t *testing.T) {
	as := setUpAccount(t)
	defer as.Close()

	createAccount(t, as)

	account, err := as.SignIn(signUpData.email, signUpData.password)

	test.AssertNil(t, err)

	AssertAccountDataMatchesSignUp(t, account)
}

func TestFailsSignIn(t *testing.T) {
	as := setUpAccount(t)
	defer as.Close()

	createAccount(t, as)

	_, err := as.SignIn(signUpData.email, "not password")

	test.AssertErrorString(t, err, "invalidPassword")

	_, err = as.SignIn("invalid@example.com", signUpData.password)

	test.AssertErrorString(t, err, "invalidEmail")
}

func TestFetchesAnAccount(t *testing.T) {
	as := setUpAccount(t)
	defer as.Close()

	account := createAccount(t, as)
	actual, err := as.GetById(account.Id)

	test.AssertNil(t, err)

	AssertAccountDataMatching(t, account, actual)
}

func TestUpdatesAnAccountIncludingPasswordAndCanSignIn(t *testing.T) {
	as := setUpAccount(t)
	defer as.Close()

	account := createAccount(t, as)

	account.Name = "b"

	updateAccount := utils.UpdateAccount{
		Id:              account.Id,
		ShippingDetails: account.ShippingDetails,
		Password:        "newpass",
	}

	actual, err := as.Update(updateAccount)

	test.AssertNil(t, err)

	AssertAccountDataMatching(t, account, actual)

	_, err = as.SignIn(updateAccount.Email, updateAccount.Password)

	test.AssertNil(t, err)
}

func AssertAccountDataMatching(t *testing.T, expected, actual utils.AccountApiResponse) {
	if expected.Address != actual.Address {
		t.Errorf("Expected address: %s, got address: %s", expected.Address, actual.Address)
	}

	if expected.Email != actual.Email {
		t.Errorf("Expected email: %s, got email: %s", expected.Email, actual.Email)
	}

	if expected.Name != actual.Name {
		t.Errorf("Expected name: %s, got name: %s", expected.Name, actual.Name)
	}

	if expected.Postcode != actual.Postcode {
		t.Errorf("Expected postcode: %s, got postcode: %s", expected.Postcode, actual.Postcode)
	}
}

func AssertAccountDataMatchesSignUp(t *testing.T, account utils.AccountApiResponse) {
	if account.Address != signUpData.address {
		t.Errorf("Expected address: %s, got address: %s", signUpData.address, account.Address)
	}

	if account.Email != signUpData.email {
		t.Errorf("Expected email: %s, got email: %s", signUpData.email, account.Email)
	}

	if account.Name != signUpData.name {
		t.Errorf("Expected name: %s, got name: %s", signUpData.name, account.Name)
	}

	if account.Postcode != signUpData.postcode {
		t.Errorf("Expected postcode: %s, got postcode: %s", signUpData.postcode, account.Postcode)
	}
}
