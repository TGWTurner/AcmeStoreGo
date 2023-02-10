package dbTests

import (
	"bjssStoreGo/backend/layers/dataAccess"
	"bjssStoreGo/backend/utils"
	"testing"
)

func SetUp() utils.Database {
	db := dataAccess.InitiateConnection()

	return db
}

func CloseDb(db utils.Database) {
	db.Order.Close()
	db.Account.Close()
	db.Product.Close()
}

func AssertNil(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Expected error: nil, got: %s", err.Error())
	}
}
