package api

import (
	bl "bjssStoreGo/backend/layers/businessLogic"
	"bjssStoreGo/backend/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func NewWiring(db utils.Database, r *mux.Router) *Wiring {

	return &Wiring{
		accountApi: NewAccountApi(bl.NewAccountService(db.Account)),
		productApi: NewProductApi(bl.NewProductService(db.Product)),
		orderApi:   NewOrderApi(bl.NewOrderService(db.Order)),
	}
}

func (w *Wiring) Run() {
	http.ListenAndServe(":3000", &w.router)
}

func (w *Wiring) SetUpRoutes() {
	app := w.router.PathPrefix("/api").Subrouter()

	account := app.PathPrefix("/account").Subrouter()
	product := app.PathPrefix("/product").Subrouter()
	order := app.PathPrefix("/order").Subrouter()

	accountSignIn := account.PathPrefix("/sign-in").Subrouter()
	accountSignIn.Use(invalidateSession)
	accountSignIn.HandleFunc("", w.accountApi.PostSignIn).Methods("POST")

	accountSignUp := account.PathPrefix("/sign-up").Subrouter()
	accountSignUp.Use(invalidateSession)
	accountSignUp.HandleFunc("", w.accountApi.PostSignUp).Methods("POST")

	accountPostGet := account.PathPrefix("").Subrouter()
	accountPostGet.Use(mustBeSignedIn)
	accountPostGet.HandleFunc("", w.accountApi.PostAccount).Methods("POST")
	accountPostGet.HandleFunc("", w.accountApi.GetAccount).Methods("GET")

	product.HandleFunc("/catalogue", w.productApi.Search).Methods("GET")
	product.HandleFunc("/categories", w.productApi.Categories).Methods("GET")
	product.HandleFunc("/deals", w.productApi.Deals).Methods("GET")

	order.HandleFunc("/basket", w.orderApi.PostBasket).Methods("POST")
	order.HandleFunc("/basket", w.orderApi.GetBasket).Methods("GET")
	order.HandleFunc("/checkout", w.orderApi.PostCheckout).Methods("POST")

	history := order.PathPrefix("/history").Subrouter()
	history.Use(mustBeSignedIn)
	history.HandleFunc("", w.orderApi.GetHistory).Methods("POST")

	// We don't use mustBeSignedIn here.
	// This allows guest customers with an order ID from a 'track my order'
	// email to fetch the order. What problems could there be with this?
	// What properties does order ID need to be secure. OWASP Top 10 #3 and #5
	order.HandleFunc("/{id}", w.orderApi.GetOrder).Methods("GET")

	w.router.PathPrefix("/").HandlerFunc(w.error404Handler)
}

func mustBeSignedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: GET THE SESSION
		if "session" != "" {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "forbidden", http.StatusUnauthorized)
		}
	})
}

func invalidateSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//INVALIDATE SESSION
	})
}

func (w Wiring) error404Handler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html")

	fmt.Fprint(writer, "<h1>Uh Oh, route not found</h1>")
	fmt.Fprint(writer, "<h2>Available routes:</h2>")
	fmt.Fprint(writer, "<ul>")
	w.printEndpoints(&w.router, writer)
	fmt.Fprint(writer, "</ul>")
}

func (w Wiring) printEndpoints(r *mux.Router, writer http.ResponseWriter) {
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}
		methods, err := route.GetMethods()
		if err != nil {
			return nil
		}
		str := "<li>[" + strings.Join(methods, "") + "] - " + path + "</li>"
		fmt.Fprint(writer, str)
		return nil
	})
}

type Wiring struct {
	router     mux.Router
	accountApi *AccountApi
	productApi *ProductApi
	orderApi   *OrderApi
}
