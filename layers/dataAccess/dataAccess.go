package dataAccess

import (
	"backend/layers/dataAccess/memory"
	"backend/layers/dataAccess/sqlite"
	"backend/utils"
	"fmt"
	"os"
)

func InitiateConnection() utils.Database {
	dbConnection := os.Getenv("DB_CONNECTION")

	fmt.Println(dbConnection)

	if dbConnection == "sql" {
		return connectToSqlite()
	} else if dbConnection == "sql-mem" {
		//In-memory SQL database. Makes unit tests a whole lot faster
		return connectToInMemory()
	}

	return memory.NewDatabase()
}

func connectToSqlite() utils.Database {
	file := "./sqlite.db"

	if _, err := os.Stat(file); !os.IsNotExist(err) {
		if err := os.Remove(file); err != nil {
			panic("Failed to remove file: " + file)
		}
	}

	return sqlite.NewDatabase(file)
}

func connectToInMemory() utils.Database {
	return sqlite.NewDatabase(":memory:")
}
