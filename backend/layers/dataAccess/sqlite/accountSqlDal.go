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

	testAccounts := ConvertToDbAccounts(testData.GetAccountTestData())

	if res := db.Create(&testAccounts); res.Error != nil {
		panic("Failed to create test accounts")
	}

	return ad
}

func (ad AccountDatabase) Close() {
	db, err := ad.db.DB()

	if err != nil {
		panic("Failed to get account db instance")
	}

	db.Close()
}

func (ad AccountDatabase) Add(account utils.Account) utils.Account {
	dbAccount := ConvertToDbAccount(account)

	result := ad.db.Create(&dbAccount)

	if result.Error != nil {
		panic(account.Email + " already registered")
	}

	return ad.GetById(int(dbAccount.ID))
}

func (ad AccountDatabase) GetByEmail(email string) utils.Account {
	var account Account

	result := ad.db.Where("email = ?", email).First(&account)

	if result.Error != nil {
		panic("Record with email: " + email + " not found")
	}

	return ConvertFromDbAccount(account)
}

func (ad AccountDatabase) GetById(accountId int) utils.Account {
	var account Account

	result := ad.db.First(&account, accountId)

	if result.Error != nil {
		panic("Record with Id: " + strconv.Itoa(accountId) + " not found")
	}

	return ConvertFromDbAccount(account)
}

func (ad AccountDatabase) Update(account utils.Account) utils.Account {
	dbAccount := ConvertToDbAccount(account)

	result := ad.db.Save(&dbAccount)

	if result.Error != nil {
		panic("Could not save record" + result.Error.Error())
	}

	return ad.GetById(int(dbAccount.ID))
}

type AccountDatabase struct {
	db *gorm.DB
}
