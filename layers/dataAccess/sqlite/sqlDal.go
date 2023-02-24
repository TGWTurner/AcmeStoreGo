package sqlite

import (
	"backend/utils"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func NewDatabase(connection string) utils.Database {
	db, err := gorm.Open(sqlite.Open(connection), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
	}

	autoMigration(db)

	return utils.Database{
		Account: NewAccountDatabase(db),
		Product: NewProductDatabase(db),
		Order:   NewOrderDatabase(db),
	}
}

func autoMigration(db *gorm.DB) {
	db.AutoMigrate(
		&Order{},
		&Product{},
		&OrderItem{},
		&Account{},
		&Deal{},
		&Category{},
	)
}
