package dbTests

import (
	"backend/layers/dataAccess/testData"
	"backend/utils"
	"reflect"
	"testing"

	"golang.org/x/exp/slices"
)

func TestCreateAccountWithNewAccountPasses(t *testing.T) {
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

	AssertNil(t, err)

	if reflect.DeepEqual(expected, actual) {
		t.Errorf("Returned account differs from expected response")
	}
}

func TestCreateAccountWithExistingAccountFails(t *testing.T) {
	db := SetUp()
	defer CloseDb(db)

	accountId := 1

	account, err := db.Account.GetById(accountId)

	AssertNil(t, err)

	_, err = db.Account.Add(account)

	if err == nil {
		t.Errorf("Succeeded createing already existing account")
	}
}

func TestGetAccountByEmail(t *testing.T) {
	db := SetUp()
	defer CloseDb(db)

	expected, err := createAccount(&db)

	AssertNil(t, err)

	email := "testEmail@email.com"

	actual, err := db.Account.GetByEmail(email)

	AssertNil(t, err)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected account: %v, got: %v", expected, actual)
	}
}

func TestGetAccountByNonExistingEmailFails(t *testing.T) {
	db := SetUp()
	defer CloseDb(db)

	email := "emailThatDoesntExist@email.com"

	_, err := db.Account.GetByEmail(email)

	if err == nil {
		t.Errorf("Expected error: nil, got error: %s", err)
	}
}

func TestGetByIdSucceedsForExistingAccount(t *testing.T) {
	db := SetUp()
	defer CloseDb(db)

	index := 1

	actual, err := db.Account.GetById(index)

	AssertNil(t, err)

	expected := getTestAccountById(index)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected account: %v, got account: %v", expected, actual)
	}
}

func TestGetByIdFailsForNonExistingId(t *testing.T) {
	db := SetUp()
	defer CloseDb(db)

	index := 1337

	_, err := db.Account.GetById(index)

	if err == nil {
		t.Errorf("Expected error: nil, got error: %s", err)
	}
}

func TestUpdateAccountSucceedsForExistingAccount(t *testing.T) {
	db := SetUp()
	defer CloseDb(db)

	accountId := 1
	newEmail := "differentTestEmail@test.com"

	account, err := db.Account.GetById(accountId)

	AssertNil(t, err)

	account.Email = newEmail

	actual, err := db.Account.Update(account)

	AssertNil(t, err)

	expected := getTestAccountById(accountId)

	expected.Email = newEmail

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected account: %v, got account: %v", expected, actual)
	}
}

func TestUpdateAccountFailsForNonExistingAccount(t *testing.T) {
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
		t.Errorf("Expected error: nil, got error: %s", err)
	}
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
