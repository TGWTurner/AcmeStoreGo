package main

import (
	"bjssStoreGo/backend/layers/api"
	da "bjssStoreGo/backend/layers/dataAccess"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
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
	setUp("sql")

	db := da.InitiateConnection()
	r := mux.NewRouter()
	store := sessions.NewCookieStore([]byte("my session encryption secret"))

	wiring := api.NewWiring(db, r, store)

	wiring.SetUpRoutes()
	// wiring.AsyncListen(":4001")
	testApiRequest(wiring)
}

func testApiRequest(wiring *api.Wiring) {
	req, err := http.NewRequest("GET", "http://localhost:4001/api/product/catalogue", nil)

	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	response := executeRequest(req, wiring)

	fmt.Println(response.Code)
}

func executeRequest(req *http.Request, wiring *api.Wiring) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	wiring.Router.ServeHTTP(rr, req)

	return rr
}

func testFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("Got to the endpoint response"))
}

func setUp(mode string) {
	//"sql" or "sql-mem" or ""
	os.Setenv("DB_CONNECTION", mode)
}
