package main

import (
	"bjssStoreGo/backend/layers/dataAccess"
	"bjssStoreGo/backend/utils"
	"fmt"
	"os"
)

func main() {
	file := "./test.db"

	unlinkFile(file)

	db := dataAccess.DataAccess{}.InitiateConnection()

	db.Account.Add(utils.Account{
		Id:           0,
		Email:        "test@test.com",
		Name:         "test",
		Address:      "test Address",
		Postcode:     "TE57 7ST",
		PasswordHash: "This Is A Password Hash",
	})

	fmt.Println(db.Account.GetByEmail("test@test.com"))

}

func unlinkFile(file string) {
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		if err := os.Remove(file); err != nil {
			panic("Failed to remove file: " + file)
		}
	}
}
