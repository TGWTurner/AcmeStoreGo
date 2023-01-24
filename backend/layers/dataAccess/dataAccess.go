package dataAccess

import (
	"bjssStoreGo/backend/layers/dataAccess/sqlite"
	"os"
)

// needs to be an interface that is returned
func (da DataAccess) InitiateConnection() sqlite.Database /*sqlite.Database*/ {
	dbConnection := "sql" //"sql" //TODO: Change to env variable check

	if dbConnection == "sql" {
		return connectToSqlite()
	} else if dbConnection == "sql-mem" {
		//In-memory SQL database. Makes unit tests a whole lot faster
		// return connectToInMemory()
	}

	// return memory.NewDatabase()
	return connectToSqlite()
}

func connectToSqlite() sqlite.Database {
	file := "./sqlite.db"

	if _, err := os.Stat(file); !os.IsNotExist(err) {
		if err := os.Remove(file); err != nil {
			panic("Failed to remove file: " + file)
		}
	}

	return sqlite.NewDatabase(file)
}

func connectToInMemory() sqlite.Database {
	return sqlite.NewDatabase(":memory:")
}

type DataAccess struct{}
