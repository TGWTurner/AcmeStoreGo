package blTests

import (
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
	"strconv"
	"strings"
)

func TestCanSignInWithExistingAccount() (bool, string, string) {
	as := SetUpAccount()
	defer as.Close()

	account := testData.GetTestAccountCredentials()

	_, err := as.SignIn(account.Email, account.Password)

	if err != nil {
		return false, "testCanSignInWithExistingAccount", "Failed to sign in with existing account"
	}

	return true, "testCanSignInWithExistingAccount", "Successfully logged in with test account"
}

func TestCannotSignInWithIncorrectPassword() (bool, string, string) {
	as := SetUpAccount()
	defer as.Close()

	account := testData.GetTestAccountCredentials()

	_, err := as.SignIn(account.Email, "NotThePassword")

	if err == nil {
		return false, "testCannotSignInWithIncorrectPassword", "Failed to reject incorrect password"
	}

	return true, "testCannotSignInWithIncorrectPassword", "Successfully rejected sign in with incorrect password"
}

func TestCannotSignInWithNonRegisteredEmail() (bool, string, string) {
	as := SetUpAccount()
	defer as.Close()

	_, err := as.SignIn("FakeEmail@fake.com", "NotAPassword")

	if err == nil {
		return false, "testCannotSignInWithNonRegisteredEmail", "Failed to reject unregistered email"
	}

	return true, "testCannotSignInWithNonRegisteredEmail", "Successfully rejected sign in with false account"
}

func TestCanSignUpWithNewAccount() (bool, string, string) {
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

	if err != nil {
		return false, "testCanSignUpWithNewAccount", "Failed to create new account"
	}

	_, createdSalt, ok := strings.Cut(retAccount.PasswordHash, ":")

	if !ok {
		return false, "testCanSignUpWithNewAccount", "Failed to get salt from returned password hash"
	}

	account.PasswordHash = as.StrongishHash(account.PasswordHash, createdSalt) + ":" + createdSalt

	if account.Address != retAccount.Address ||
		account.Email != retAccount.Email ||
		account.Name != retAccount.Name ||
		account.PasswordHash != retAccount.PasswordHash {
		return false, "testCanSignUpWithNewAccount", "Failed to return correct created account"
	}

	return true, "testCanSignUpWithNewAccount", "Successfully created new account"
}

func TestCannotSignUpWithAlreadyRegisteredAccount() (bool, string, string) {
	as := SetUpAccount()
	defer as.Close()

	account := testData.GetTestAccountCredentials()

	_, err := as.SignUp(account.Email, "ThisIsAPassword", "name", "address", "postcode")

	if err == nil {
		return false, "testCannotSignUpWithAlreadyRegisteredAccount", "Failed to reject creation of already registered account"
	}

	return true, "testCannotSignUpWithAlreadyRegisteredAccount", "Successfully rejected creation of already registered account"
}

func TestCanUpdateExistingAccount() (bool, string, string) {
	as := SetUpAccount()
	defer as.Close()

	emailAndPass := testData.GetTestAccountCredentials()

	account, err := as.SignIn(emailAndPass.Email, emailAndPass.Password)

	if err != nil {
		return false, "testCanUpdateExistingAccount", "Failed to get test account"
	}

	retAccount, err := as.Update(account, "newPassword")

	_, createdSalt, ok := strings.Cut(retAccount.PasswordHash, ":")

	if !ok {
		return false, "testCanSignUpWithNewAccount", "Failed to get salt from returned password hash"
	}

	account.PasswordHash = as.StrongishHash(account.PasswordHash, createdSalt) + ":" + createdSalt

	if err != nil {
		return false, "testCanUpdateExistingAccount", "Failed to update account"
	}

	if retAccount != account {
		return false, "testCanUpdateExistingAccount", "Returned account after update does not match provided"
	}

	return true, "testCanUpdateExistingAccount", "Succeeded updating account"
}

func TestCannotUpdateUnregisteredAccount() (bool, string, string) {
	as := SetUpAccount()
	defer as.Close()

	account := utils.Account{
		ShippingDetails: utils.ShippingDetails{
			Email:    "unregisteredEmail@email.com",
			Name:     "name",
			Address:  "address",
			Postcode: "postcode",
		},
		PasswordHash: as.HashPassword("ThisIsAPassword"),
	}

	_, err := as.Update(account)

	if err == nil {
		return false, "testCannotUpdateUnregisteredAccount", "Failed to reject unregistered account for update"
	}

	return true, "testCannotUpdateUnregisteredAccount", "Successfully rejected unregistered account for update"
}

func TestCanGetAccountById() (bool, string, string) {
	as := SetUpAccount()
	defer as.Close()

	accountId := 1

	_, err := as.GetById(accountId)

	if err != nil {
		return false, "testCanGetAccountById", "Failed to get account with id: " + strconv.Itoa(accountId)
	}

	return true, "testCanGetAccountById", "Successfully got account for id: " + strconv.Itoa(accountId)
}

func TestCannotGetAccountWithFakeId() (bool, string, string) {
	as := SetUpAccount()
	defer as.Close()

	accountId := 1337

	_, err := as.GetById(accountId)

	if err == nil {
		return false, "testCannotGetAccountWithFakeId", "Failed to reject getting account with false account id"
	}

	return true, "testCannotGetAccountWithFakeId", "Successfully rejected getting account with false account id"
}
