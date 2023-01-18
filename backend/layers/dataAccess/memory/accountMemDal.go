package memory

import (
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
	"fmt"
	"sort"
	"strconv"
)

func NewAccountDatabase() AccountDatabase {
	testAccounts := testData.AccountTestData{}.GetTestData()

	return AccountDatabase{
		accounts: testAccounts,
	}
}

func (ad *AccountDatabase) GetAccount() []utils.Account {
	return ad.accounts
}

func (ad *AccountDatabase) Add(account utils.Account) {
	index := sort.Search(len(ad.accounts), func(i int) bool { return ad.accounts[i].Email == account.Email })

	fmt.Println("adding account: " + account.Email)
	fmt.Println("Index of acc: " + strconv.Itoa(index))

	if index != len(ad.accounts) {
		panic(account.Email + " already registered")
	}

	ad.accounts = append(ad.accounts, account)

	for _, acc := range ad.accounts {
		fmt.Println("=======")
		fmt.Println(acc.Address)
		fmt.Println(acc.Name)
		fmt.Println(acc.Email)
		fmt.Println(acc.Id)
		fmt.Println(acc.Postcode)
	}

	fmt.Println("finish add")
}

func (ad *AccountDatabase) GetByEmail(email string) utils.Account {
	index := sort.Search(len(ad.accounts), func(i int) bool { return ad.accounts[i].Email == email })

	if index == len(ad.accounts) {
		panic("Record with email: " + email + " not found")
	}

	return ad.accounts[index]
}

func (ad *AccountDatabase) GetById(accountId int) utils.Account {
	index := sort.Search(len(ad.accounts), func(i int) bool { return ad.accounts[i].Id == accountId })

	if index == len(ad.accounts) {
		panic("Record with Id: " + strconv.Itoa(accountId) + " not found")
	}

	return ad.accounts[index]
}

func (ad *AccountDatabase) Update(account utils.Account) {
	index := sort.Search(len(ad.accounts), func(i int) bool { return ad.accounts[i].Email == account.Email })

	if index == len(ad.accounts) {
		panic("Could not find record for account with email: " + account.Email)
	}

	ad.accounts[index] = account
}

type AccountDatabase struct {
	accounts []utils.Account
}
