package memory

import (
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
	"errors"
	"strconv"
)

func NewAccountDatabase() *AccountDatabaseImpl {
	testAccounts := testData.GetAccountTestData()

	return &AccountDatabaseImpl{
		accounts: testAccounts,
	}
}

func (ad *AccountDatabaseImpl) Close() {}

func (ad *AccountDatabaseImpl) Add(account utils.Account) (utils.Account, error) {
	index := ad.getIndexFromEmail(account.Email)

	if index != len(ad.accounts) {
		return utils.Account{}, errors.New(account.Email + " already registered")
	}

	account.Id = len(ad.accounts) + 1

	ad.accounts = append(ad.accounts, account)

	return ad.GetById(account.Id)
}

func (ad *AccountDatabaseImpl) GetByEmail(email string) (utils.Account, error) {
	index := ad.getIndexFromEmail(email)

	if index == len(ad.accounts) {
		return utils.Account{},
			errors.New("Record with email: " + email + " not found")
	}

	return ad.accounts[index], nil
}

func (ad *AccountDatabaseImpl) GetById(accountId int) (utils.Account, error) {
	index := ad.getIndexFromId(accountId)

	if index == len(ad.accounts) {
		return utils.Account{}, errors.New("Record with Id: " + strconv.Itoa(accountId) + " not found")
	}

	return ad.accounts[index], nil
}

func (ad *AccountDatabaseImpl) Update(updateAccount utils.Account) (utils.Account, error) {
	index := ad.getIndexFromId(updateAccount.Id)

	if index == len(ad.accounts) {
		return utils.Account{}, errors.New("Could not find record for account for id: " + strconv.Itoa(updateAccount.Id))
	}

	ad.accounts[index] = updateAccount

	return ad.accounts[index], nil
}

func (ad *AccountDatabaseImpl) getIndexFromEmail(email string) int {
	index := len(ad.accounts)

	for i, a := range ad.accounts {
		if a.Email == email {
			index = i
			break
		}
	}

	return index
}

func (ad *AccountDatabaseImpl) getIndexFromId(id int) int {
	index := len(ad.accounts)

	for i, a := range ad.accounts {
		if a.Id == id {
			index = i
			break
		}
	}

	return index
}

type AccountDatabaseImpl struct {
	accounts []utils.Account
}
