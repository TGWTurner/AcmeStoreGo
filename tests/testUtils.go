package tests

import (
	"bjssStoreGo/backend/layers/dataAccess"
	"bjssStoreGo/backend/layers/dataAccess/sqlite"
	"fmt"
)

func SetUp() /*memory.Database*/ sqlite.Database {
	db := dataAccess.InitiateConnection()

	return db
}

func CloseDb(db /*memory.Database*/ sqlite.Database) {
	db.Order.Close()
	db.Account.Close()
	db.Product.Close()
}

func PrintTestResult(pass bool, testName string, message string) {
	fmt.Print("TEST ")
	if pass {
		fmt.Print("PASSED")
	} else {
		fmt.Print("FAILED")
	}
	fmt.Println(" -- " + testName)
	fmt.Printf("\t" + message + "\n")
}
