package sqlite

import (
	"backend/layers/dataAccess/testData"
	"backend/utils"
	"errors"
	"reflect"
	"strconv"

	"gorm.io/gorm"
)

func NewAccountDatabase(db *gorm.DB) *AccountDatabaseImpl {
	ad := AccountDatabaseImpl{
		db: db,
	}

	testAccounts := ConvertToDbAccounts(testData.GetAccountTestData())

	if res := db.Create(&testAccounts); res.Error != nil {
		panic("Failed to create test accounts")
	}

	return &ad
}

func (ad AccountDatabaseImpl) Close() {
	db, err := ad.db.DB()

	if err != nil {
		panic("Failed to get account db instance")
	}

	db.Close()
}

func (ad AccountDatabaseImpl) Add(account utils.Account) (utils.Account, error) {
	dbAccount := ConvertToDbAccount(account)

	_, err := ad.GetByEmail(account.Email)

	if err == nil {
		return utils.Account{}, errors.New("Account with email: " + account.Email + " already exists")
	}

	result := ad.db.Create(&dbAccount)

	if result.Error != nil {
		return utils.Account{}, result.Error
	}

	return ad.GetById(int(dbAccount.Id))
}

func (ad AccountDatabaseImpl) GetByEmail(email string) (utils.Account, error) {
	var account Account

	result := ad.db.Where("email = ?", email).Limit(1).Find(&account)

	if reflect.DeepEqual(account, Account{}) || result.Error != nil {
		return utils.Account{}, errors.New("Account with email: " + email + " not found")
	}

	return ConvertFromDbAccount(account), nil
}

func (ad AccountDatabaseImpl) GetById(accountId int) (utils.Account, error) {
	var account Account

	result := ad.db.Limit(1).Find(&account, accountId)

	if reflect.DeepEqual(account, Account{}) || result.Error != nil {
		return utils.Account{}, errors.New("Account with id: " + strconv.Itoa(accountId) + " not found")
	}

	return ConvertFromDbAccount(account), nil
}

func (ad AccountDatabaseImpl) Update(account utils.Account) (utils.Account, error) {
	dbAccount := ConvertToDbAccount(account)

	_, err := ad.GetById(account.Id)

	if err != nil {
		return utils.Account{}, err
	}

	result := ad.db.Save(&dbAccount)

	if result.Error != nil {
		return utils.Account{}, result.Error
	}

	return ad.GetById(int(dbAccount.Id))
}

type AccountDatabaseImpl struct {
	db *gorm.DB
}
