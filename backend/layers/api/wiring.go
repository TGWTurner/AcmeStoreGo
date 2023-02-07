package api

import (
	bl "bjssStoreGo/backend/layers/businessLogic"
	"bjssStoreGo/backend/layers/dataAccess"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func NewWiring() *Wiring {
	db := dataAccess.InitiateConnection()

	return &Wiring{
		accountApi: NewAccountApi(bl.NewAccountService(db.Account)),
		productApi: NewProductApi(bl.NewProductService(db.Product)),
		orderApi:   NewOrderApi(bl.NewOrderService(db.Order)),
	}
}

func (w *Wiring) SetUpRoutes(r *mux.Router) {
	db := dataAccess.InitiateConnection()
	accountApi := NewAccountApi(bl.NewAccountService(db.Account))
	productApi := NewProductApi(bl.NewProductService(db.Product))
	orderApi := NewOrderApi(bl.NewOrderService(db.Order))

	app := r.PathPrefix("/api").Subrouter()

	account := app.PathPrefix("/account").Subrouter()
	product := app.PathPrefix("/product").Subrouter()
	order := app.PathPrefix("/order").Subrouter()

	accountSignIn := account.PathPrefix("/sign-in").Subrouter()
	accountSignIn.Use(invalidateSession)
	accountSignIn.HandleFunc("", accountApi.PostSignIn).Methods("POST")

	accountSignUp := account.PathPrefix("/sign-up").Subrouter()
	accountSignUp.Use(invalidateSession)
	accountSignUp.HandleFunc("", accountApi.PostSignUp).Methods("POST")

	accountPostGet := account.PathPrefix("").Subrouter()
	accountPostGet.Use(mustBeSignedIn)
	accountPostGet.HandleFunc("", accountApi.PostAccount).Methods("POST")
	accountPostGet.HandleFunc("", accountApi.GetAccount).Methods("GET")

	product.HandleFunc("/catalogue", productApi.Search).Methods("GET")
	product.HandleFunc("/categories", productApi.Categories).Methods("GET")
	product.HandleFunc("/deals", productApi.Deals).Methods("GET")

	order.HandleFunc("/basket", orderApi.PostBasket).Methods("POST")
	order.HandleFunc("/basket", orderApi.GetBasket).Methods("GET")
	order.HandleFunc("/checkout", orderApi.PostCheckout).Methods("POST")

	history := order.PathPrefix("/history").Subrouter()
	history.Use(mustBeSignedIn)
	history.HandleFunc("", orderApi.GetHistory).Methods("POST")

	// We don't use mustBeSignedIn here.
	// This allows guest customers with an order ID from a 'track my order'
	// email to fetch the order. What problems could there be with this?
	// What properties does order ID need to be secure. OWASP Top 10 #3 and #5
	order.HandleFunc("/{id}", orderApi.GetOrder).Methods("GET")

	app.PathPrefix("/").HandlerFunc(error404Handler)

	printEndpoints(r)

	http.ListenAndServe(":3000", r)
}

func printEndpoints(r *mux.Router) {
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}
		methods, err := route.GetMethods()
		if err != nil {
			return nil
		}
		fmt.Printf("%v %s\n", methods, path)
		return nil
	})
}

func mustBeSignedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: FIX THE SESSION?
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

func error404Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Uh Oh</h1>")
}

type Wiring struct {
	accountApi *AccountApi
	productApi *ProductApi
	orderApi   *OrderApi
}
