package sqlite

import (
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
		&Order{},
		&Product{},
		&OrderItem{},
		&Account{},
		&Deal{},
		&Category{},
	)
}

type Database struct {
	Account AccountDatabase
	Product ProductDatabase
	Order   OrderDatabase
}
