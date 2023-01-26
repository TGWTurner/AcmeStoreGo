package tests

import (
	"bjssStoreGo/backend/layers/dataAccess"
	"bjssStoreGo/backend/layers/dataAccess/sqlite"
	"fmt"
	"os"
)

func SetUp() sqlite.Database {
	file := "./sqlite.db"

	if _, err := os.Stat(file); !os.IsNotExist(err) {
		if err := os.Remove(file); err != nil {
			fmt.Println(err)
			panic("Failed to remove file: " + file)
		}
	}
	db := dataAccess.InitiateConnection()

	return db
}

func CloseDb(db sqlite.Database) {
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
