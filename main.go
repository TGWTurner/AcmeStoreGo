package main

import (
	"bjssStoreGo/backend/layers/api"
	da "bjssStoreGo/backend/layers/dataAccess"
	"fmt"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func main() {
	//SET UP <<<<<<<<<<<<<<<
	setUp("sql")

	db := da.InitiateConnection()
	r := mux.NewRouter()
	store := sessions.NewCookieStore([]byte("my session encryption secret"))

	wiring := api.NewWiring(db, r, store)

	wiring.SetUpRoutes()
	//SET UP <<<<<<<<<<<<<<<

	fmt.Println("Starting...")
	wiring.AsyncListen(":4001")
}

func setUp(mode string) {
	//"sql" or "sql-mem" or ""
	os.Setenv("DB_CONNECTION", mode)
}
