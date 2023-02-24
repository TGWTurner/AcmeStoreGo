package sqlite

import (
	"backend/utils"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Id           int    `gorm:"primaryKey"`
	Email        string `gorm:"not null; unique"`
	Name         string `gorm:"not null"`
	Address      string `gorm:"not null"`
	Postcode     string `gorm:"not null"`
	PasswordHash string `gorm:"not null"`
}

func ConvertToDbAccount(account utils.Account) Account {
	return Account{
		Id:           account.Id,
		Email:        account.Email,
		Name:         account.Name,
		Address:      account.Address,
		Postcode:     account.Postcode,
		PasswordHash: account.PasswordHash,
	}
}

func ConvertToDbAccounts(accounts []utils.Account) []Account {
	dbAccounts := make([]Account, len(accounts))

	for i, account := range accounts {
		dbAccounts[i] = ConvertToDbAccount(account)
	}

	return dbAccounts
}

func ConvertFromDbAccount(account Account) utils.Account {
	return utils.Account{
		Id:           account.Id,
		PasswordHash: account.PasswordHash,
		ShippingDetails: utils.ShippingDetails{
			Email:    account.Email,
			Name:     account.Name,
			Address:  account.Address,
			Postcode: account.Postcode,
		},
	}
}

func ConvertFromDbAccounts(dbAccounts []Account) []utils.Account {
	accounts := make([]utils.Account, len(dbAccounts))

	for i, account := range dbAccounts {
		accounts[i] = ConvertFromDbAccount(account)
	}

	return accounts
}
