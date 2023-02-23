package main

import (
	"bjssStoreGo/backend/layers/api"
	da "bjssStoreGo/backend/layers/dataAccess"
	"bjssStoreGo/backend/utils"
	"fmt"
	"os"

	"encoding/gob"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// @title BJSS Store
// @version 1.0
// @description Simple store for teaching and learning
// @host  localhost:4001
// @BasePath /
func main() {
	//SET UP <<<<<<<<<<<<<<<
	setUp("sql")

	db := da.InitiateConnection()
	r := mux.NewRouter()
	store := sessions.NewCookieStore([]byte("my session encryption secret"))

	wiring := api.NewWiring(db, r, store)

	wiring.SetUpRoutes()
	//SET UP <<<<<<<<<<<<<<<

	fmt.Println("server listening on port 4001, try http://localhost:4001/api-docs/")
	wiring.AsyncListen(":4001")
}

func setUp(mode string) {
	//"sql" or "sql-mem" or ""
	os.Setenv("DB_CONNECTION", mode)

	//register basket to be used with session
	gob.Register(utils.Basket{})
}
