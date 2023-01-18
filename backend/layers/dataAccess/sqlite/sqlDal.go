package sqlite

import (
	"bjssStoreGo/backend/utils"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func NewDatabase(connection string) Database {
	db, err := gorm.Open(sqlite.Open(connection), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
	}

	autoMigration(db)

	return Database{
		Account: NewAccountDatabase(db),
		Product: NewProductDatabase(db),
		Order:   NewOrderDatabase(db),
	}
}

func autoMigration(db *gorm.DB) {
	db.AutoMigrate(
		&utils.Order{},
		&utils.Product{},
		&utils.OrderItem{},
		&utils.Account{},
		&utils.Deal{},
		&utils.Category{},
	)
}

type Database struct {
	Account AccountDatabase
	Product ProductDatabase
	Order   OrderDatabase
}
