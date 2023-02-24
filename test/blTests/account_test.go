package blTests

import (
	"backend/layers/businessLogic"
	"backend/layers/dataAccess/testData"
	"backend/utils"
	"errors"
	"strings"
	"testing"
)

func TestCanSignInWithExistingAccount(t *testing.T) {
	as := SetUpAccount()
	defer as.Close()

	account := testData.GetTestAccountCredentials()

	_, err := as.SignIn(account.Email, account.Password)

	AssertNil(t, err)
}

func TestCannotSignInWithIncorrectPassword(t *testing.T) {
	as := SetUpAccount()
	defer as.Close()

	account := testData.GetTestAccountCredentials()

	_, err := as.SignIn(account.Email, "NotThePassword")

	if err == nil {
		t.Errorf("Failed to reject incorrect password")
	}
}

func TestCannotSignInWithNonRegisteredEmail(t *testing.T) {
	as := SetUpAccount()
	defer as.Close()

	_, err := as.SignIn("FakeEmail@fake.com", "NotAPassword")

	if err == nil {
		t.Errorf("Expected error: nil, got error: %s", err)
	}
}

func TestCanSignUpWithNewAccount(t *testing.T) {
	as := SetUpAccount()
	defer as.Close()

	account := utils.Account{
		ShippingDetails: utils.ShippingDetails{
			Name:     "NewAccountName",
			Email:    "NewAccountEmail@email.com",
			Postcode: "PO5 7DE",
			Address:  "NewAccountAddress",
		},
		PasswordHash: "Password",
	}

	retAccount, err := as.SignUp(account.Email, account.PasswordHash, account.Name, account.Address, account.Postcode)

	AssertNil(t, err)

	if account.Address != retAccount.Address ||
		account.Email != retAccount.Email ||
		account.Name != retAccount.Name {
		t.Errorf("Failed to return correct created account")
	}
}

func TestCannotSignUpWithAlreadyRegisteredAccount(t *testing.T) {
	as := SetUpAccount()
	defer as.Close()

	account := testData.GetTestAccountCredentials()

	_, err := as.SignUp(account.Email, "ThisIsAPassword", "name", "address", "postcode")

	if err == nil {
		t.Errorf("Expected: error, got: nil")
	}
}

func TestCanUpdateExistingAccount(t *testing.T) {
	as := SetUpAccount()
	defer as.Close()

	emailAndPass := testData.GetTestAccountCredentials()

	account, err := as.SignIn(emailAndPass.Email, emailAndPass.Password)

	AssertNil(t, err)

	updateAccount := utils.UpdateAccount{
		Id:              account.Id,
		ShippingDetails: account.ShippingDetails,
	}

	retAccount, err := as.Update(updateAccount)

	AssertNil(t, err)

	if retAccount != account {
		t.Errorf("Returned account after update does not match provided")
	}
}

func TestCannotUpdateUnregisteredAccount(t *testing.T) {
	as := SetUpAccount()
	defer as.Close()

	account := utils.UpdateAccount{
		ShippingDetails: utils.ShippingDetails{
			Email:    "unregisteredEmail@email.com",
			Name:     "name",
			Address:  "address",
			Postcode: "postcode",
		},
		Password: "ThisIsAPassword",
	}

	_, err := as.Update(account)

	if err == nil {
		t.Errorf("Failed to reject unregistered account for update")
	}
}

func TestCanGetAccountById(t *testing.T) {
	as := SetUpAccount()
	defer as.Close()

	accountId := 1

	_, err := as.GetById(accountId)

	AssertNil(t, err)
}

func TestCannotGetAccountWithFakeId(t *testing.T) {
	as := SetUpAccount()
	defer as.Close()

	accountId := 1337

	_, err := as.GetById(accountId)

	if err == nil {
		t.Errorf("Failed to reject getting account with false account id")
	}
}

func hashPassWithProvidedSalt(as businessLogic.AccountService, password string, hashNsalt string) (string, error) {
	_, createdSalt, ok := strings.Cut(hashNsalt, ":")

	if !ok {
		return "", errors.New("Failed to get salt from returned password hash")
	}

	return as.StrongishHash(password, createdSalt) + ":" + createdSalt, nil
}
