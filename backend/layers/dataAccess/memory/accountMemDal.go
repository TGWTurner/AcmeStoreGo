package memory

import (
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
	"sort"
	"strconv"
)

func NewAccountDatabase() AccountDatabase {
	testAccounts := testData.GetAccountTestData()

	return AccountDatabase{
		accounts: testAccounts,
	}
}

func (ad *AccountDatabase) Close() {}

func (ad *AccountDatabase) GetAccounts() []utils.Account {
	return ad.accounts
}

func (ad *AccountDatabase) Add(account utils.Account) utils.Account {
	index := sort.Search(len(ad.accounts), func(i int) bool { return ad.accounts[i].Email == account.Email })

	if index != len(ad.accounts) {
		panic(account.Email + " already registered")
	}

	account.Id = len(ad.accounts) + 1

	ad.accounts = append(ad.accounts, account)

	return ad.GetById(account.Id)
}

func (ad *AccountDatabase) GetByEmail(email string) utils.Account {

	index := sort.Search(len(ad.accounts), func(i int) bool { return ad.accounts[i].Email == email })

	if index == len(ad.accounts) {
		panic("Record with email: " + email + " not found")
	}

	return ad.accounts[index]
}

func (ad *AccountDatabase) GetById(accountId int) utils.Account {
	index := len(ad.accounts)

	for i, a := range ad.accounts {
		if a.Id == accountId {
			index = i
			break
		}
	}

	if index == len(ad.accounts) {
		panic("Record with Id: " + strconv.Itoa(accountId) + " not found")
	}

	return ad.accounts[index]
}

func (ad *AccountDatabase) Update(account utils.Account) utils.Account {
	index := sort.Search(len(ad.accounts), func(i int) bool { return ad.accounts[i].Email == account.Email })

	if index == len(ad.accounts) {
		panic("Could not find record for account with email: " + account.Email)
	}

	ad.accounts[index] = account

	return ad.accounts[index]
}

type AccountDatabase struct {
	accounts []utils.Account
}
