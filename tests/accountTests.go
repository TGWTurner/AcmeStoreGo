package tests

import (
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
	"reflect"
	"strconv"

	"golang.org/x/exp/slices"
)

func TestCreateAccountWithNewAccountPasses() (bool, string, string) {
	db := SetUp()
	defer CloseDb(db)

	expected := utils.Account{
		PasswordHash: "This is a password",
		ShippingDetails: utils.ShippingDetails{
			Email:    "testEmail@email.com",
			Name:     "testName",
			Address:  "testAddress",
			Postcode: "testPostcode",
		},
	}

	actual, err := db.Account.Add(expected)

	if err != nil {
		return false,
			"testCreateAccountWithNewAccountPasses",
			"Account addition error: " + err.Error()
	}

	if reflect.DeepEqual(expected, actual) {
		return false,
			"testCreateAccountWithNewAccountPasses",
			"Returned account differs from expected response"
	}

	return true,
		"testCreateAccountWithNewAccountPasses",
		"Created account successfully"
}

func TestCreateAccountWithExistingAccountFails() (bool, string, string) {
	db := SetUp()
	defer CloseDb(db)

	accountId := 1

	account, err := db.Account.GetById(accountId)

	if err != nil {
		return false,
			"testCreateAccountWithExistingAccountFails",
			"Failed to get account with Id: " + strconv.Itoa(accountId)
	}

	_, err = db.Account.Add(account)

	if err == nil {
		return false,
			"testCreateAccountWithExistingAccountFails",
			"Succeeded creating already existing account"
	}

	return true,
		"testCreateAccountWithExistingAccountFails",
		"Failed to create already existing account"
}

func TestGetAccountByEmail() (bool, string, string) {
	db := SetUp()
	defer CloseDb(db)

	expected, err := createAccount(&db)

	if err != nil {
		return false,
			"testGetAccountByEmail",
			"Failed to create account, err: " + err.Error()
	}

	email := "testEmail@email.com"

	actual, err := db.Account.GetByEmail(email)

	if err != nil {
		return false,
			"testGetAccountByEmail",
			"Failed to get account for email:" + email
	}

	if !reflect.DeepEqual(expected, actual) {
		return false,
			"testGetAccountByEmail",
			"Returned account differs from expected response for email: " + email
	}

	return true,
		"testGetAccountByEmail",
		"Account with email " + email + " successfully retrieved"
}

func TestGetAccountByNonExistingEmailFails() (bool, string, string) {
	db := SetUp()
	defer CloseDb(db)

	email := "emailThatDoesntExist@email.com"

	_, err := db.Account.GetByEmail(email)

	if err == nil {
		return false,
			"testGetAccountByNonExistingEmailFails",
			"Retrieved account when should not have"
	}

	return true,
		"testGetAccountByNonExistingEmailFails",
		"Successfully returned error when retrieved for non existant email"
}

func TestGetByIdSucceedsForExistingAccount() (bool, string, string) {
	db := SetUp()
	defer CloseDb(db)

	index := 1

	actual, err := db.Account.GetById(index)

	if err != nil {
		return false,
			"testGetByIdSucceedsForExistingAccount",
			"Failed to get Account with id: " + strconv.Itoa(index)
	}

	expected := getTestAccountById(index)

	if !reflect.DeepEqual(expected, actual) {
		return false,
			"testGetByIdSucceedsForExistingAccount",
			"Failed to get Account with id: " + strconv.Itoa(index)
	}

	return true,
		"testGetById",
		"Successfully got Account with id: " + strconv.Itoa(index)
}

func TestGetByIdFailsForNonExistingId() (bool, string, string) {
	db := SetUp()
	defer CloseDb(db)

	index := 1337

	_, err := db.Account.GetById(index)

	if err == nil {
		return false,
			"testGetByIdFailsForNonExistingId",
			"Recieved account when should not have"
	}

	return true,
		"testGetById",
		"Successfully failed to get account with non existant id"
}

func TestUpdateAccountSucceedsForExistingAccount() (bool, string, string) {
	db := SetUp()
	defer CloseDb(db)

	accountId := 1
	newEmail := "differentTestEmail@test.com"

	account, err := db.Account.GetById(accountId)

	if err != nil {
		return false,
			"testUpdateAccountSucceedsForExistingAccount",
			"Failed to get account with Id: " + strconv.Itoa(accountId)
	}

	account.Email = newEmail

	actual, err := db.Account.Update(account)
	expected := getTestAccountById(accountId)

	expected.Email = newEmail

	if !reflect.DeepEqual(expected, actual) || err != nil {
		return false,
			"testUpdateAccountSucceedsForExistingAccount",
			"Failed to update email of account with Id: " + strconv.Itoa(actual.Id)
	}

	return true,
		"testUpdateAccountSucceedsForExistingAccount",
		"Succeeded updating email for account with id: " + strconv.Itoa(actual.Id)
}

func TestUpdateAccountFailsForNonExistingAccount() (bool, string, string) {
	db := SetUp()
	defer CloseDb(db)

	account := utils.Account{
		Id:           1337,
		PasswordHash: "This is a non existing password",
		ShippingDetails: utils.ShippingDetails{
			Email:    "nonExisting@email.com",
			Name:     "nonExistingName",
			Address:  "nonExistingAddress",
			Postcode: "nonExistingPostcode",
		},
	}

	_, err := db.Account.Update(account)

	if err == nil {
		return false,
			"testUpdateAccountFailsForNonExistingAccount",
			"Succeeded updating account that does not exist"
	}

	return true,
		"testUpdateAccountFailsForNonExistingAccount",
		"Successfully failed to update non existant account"
}

func createAccount(db *utils.Database) (utils.Account, error) {
	expected := utils.Account{
		PasswordHash: "This is a password",
		ShippingDetails: utils.ShippingDetails{
			Email:    "testEmail@email.com",
			Name:     "testName",
			Address:  "testAddress",
			Postcode: "testPostcode",
		},
	}

	return db.Account.Add(expected)
}

func getTestAccountById(id int) utils.Account {
	expectedAccounts := testData.GetAccountTestData()
	idx := slices.IndexFunc(
		expectedAccounts,
		func(a utils.Account) bool { return a.Id == id },
	)

	return expectedAccounts[idx]
}
