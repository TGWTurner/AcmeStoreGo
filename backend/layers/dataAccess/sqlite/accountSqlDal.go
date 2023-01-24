package sqlite

import (
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
	"strconv"

	"gorm.io/gorm"
)

func NewAccountDatabase(db *gorm.DB) AccountDatabase {
	ad := AccountDatabase{
		db: db,
	}

	db.AutoMigrate(&utils.Account{})

	testAccounts := testData.AccountTestData{}.GetTestData()
	db.Create(&testAccounts)

	return ad
}

func (ad AccountDatabase) Add(account utils.Account) {
	result := ad.db.Create(&account)

	if result.Error != nil {
		panic(account.Email + " already registered")
	}
}

func (ad AccountDatabase) GetByEmail(email string) utils.Account {
	var account utils.Account

	result := ad.db.Where("email = ?", email).First(&account)

	if result.Error != nil {
		panic("Record with email: " + email + " not found")
	}

	return account
}

func (ad AccountDatabase) GetById(accountId int) utils.Account {
	var account utils.Account

	result := ad.db.First(&account, accountId)

	if result.Error != nil {
		panic("Record with Id: " + strconv.Itoa(accountId) + " not found")
	}

	return account
}

func (ad AccountDatabase) Update(account utils.Account) {
	result := ad.db.Save(&account)

	if result.Error != nil {
		panic("Could not save record" + result.Error.Error())
	}
}

type AccountDatabase struct {
	db *gorm.DB
}
