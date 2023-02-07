package main

import (
	"bjssStoreGo/backend/layers/api"
	"bjssStoreGo/blTests"
	"bjssStoreGo/dbTests"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func page(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "hello page")
}

func landingPage() http.Handler {
	h := http.HandlerFunc(page)
	return h
}

func main() {
	//Run tests:
	//tests()
	wiring := api.NewWiring()

	r := mux.NewRouter()

	wiring.SetUpRoutes(r)
}

func testFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("Got to the endpoint response"))
}

func tests() {
	runDbTests()
	runBlTests()
}

func runBlTests() {
	fmt.Println()

	setUp("sql")
	fmt.Println("Sql")
	blTests.RunTests()

	fmt.Println()

	setUp("memory")
	fmt.Println("Memory")
	blTests.RunTests()
}

func runDbTests() {
	fmt.Println()

	setUp("sql")
	fmt.Println("Sql")
	dbTests.RunTests()

	fmt.Println()

	setUp("memory")
	fmt.Println("Memory")
	dbTests.RunTests()
}

func setUp(mode string) {
	//"sql" or "sql-mem" or ""
	os.Setenv("DB_CONNECTION", mode)
}
