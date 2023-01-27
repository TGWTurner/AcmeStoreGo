package dataAccess

import (
	"bjssStoreGo/backend/layers/dataAccess/memory"
	"bjssStoreGo/backend/layers/dataAccess/sqlite"
	"bjssStoreGo/backend/utils"
	"os"
)

func InitiateConnection() utils.Database {
	dbConnection := os.Getenv("DB_CONNECTION")

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
