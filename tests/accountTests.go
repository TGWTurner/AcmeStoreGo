package tests

import (
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
	"reflect"
	"strconv"

	"golang.org/x/exp/slices"
)

func TestCreateAccount() {
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

	actual := db.Account.Add(expected)

	if reflect.DeepEqual(expected, actual) {
		PrintTestResult(
			false,
			"testCreateAccount",
			"Returned account differs from expected response",
		)
		return
	}

	PrintTestResult(
		true,
		"testCreateAccount",
		"Created account successfully",
	)
}

func TestGetAccountByEmail() {
	db := SetUp()
	defer CloseDb(db)

	expected := createAccount(&db)
	email := "testEmail@email.com"

	actual := db.Account.GetByEmail(email)

	if !reflect.DeepEqual(expected, actual) {
		PrintTestResult(
			false,
			"testGetAccountByEmail",
			"Returned account differs from expected response for email: "+email,
		)
		return
	}

	PrintTestResult(
		true,
		"testGetAccountByEmail",
		"Account with email "+email+" successfully retrieved",
	)
}

func TestGetById() {
	db := SetUp()
	defer CloseDb(db)

	index := 1

	actual := db.Account.GetById(index)
	expected := getTestAccountById(index)

	if !reflect.DeepEqual(expected, actual) {
		PrintTestResult(
			false,
			"testGetById",
			"Failed to get Account with id: "+strconv.Itoa(index),
		)
		return
	}

	PrintTestResult(
		true,
		"testGetById",
		"Successfully got Account with id: "+strconv.Itoa(index),
	)
}

func TestUpdateAccount() {
	db := SetUp()
	defer CloseDb(db)

	accountId := 1
	newEmail := "differentTestEmail@test.com"

	account := db.Account.GetById(accountId)
	account.Email = newEmail

	actual := db.Account.Update(account)
	expected := getTestAccountById(accountId)

	expected.Email = newEmail

	if !reflect.DeepEqual(expected, actual) {
		PrintTestResult(
			false,
			"testUpdateAccount",
			"Failed to update email of account with Id: "+strconv.Itoa(actual.Id),
		)
	}

	PrintTestResult(
		true,
		"testUpdateAccount",
		"Succeeded updating email for account with id: "+strconv.Itoa(actual.Id),
	)
}

func createAccount(db *utils.Database) utils.Account {
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
