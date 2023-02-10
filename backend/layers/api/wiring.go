package api

import (
	bl "bjssStoreGo/backend/layers/businessLogic"
	"bjssStoreGo/backend/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func NewWiring(db utils.Database, r *mux.Router, s *sessions.CookieStore) *Wiring {
	as := bl.NewAccountService(db.Account)
	ps := bl.NewProductService(db.Product)
	os := bl.NewOrderService(db.Order, *ps)

	return &Wiring{
		store:      s,
		router:     r,
		accountApi: NewAccountApi(as, s),
		productApi: NewProductApi(ps, s),
		orderApi:   NewOrderApi(os, s),
	}
}

func (w *Wiring) AsyncListen(port string) {
	http.ListenAndServe(port, w.router)
}

func (w *Wiring) SetUpRoutes() {
	app := w.router.PathPrefix("/api").Subrouter()

	app.HandleFunc("/values", w.Values)

	account := app.PathPrefix("/account").Subrouter()
	product := app.PathPrefix("/product").Subrouter()
	order := app.PathPrefix("/order").Subrouter()

	postSignIn := http.HandlerFunc(w.accountApi.PostSignIn)
	account.Handle("/sign-in", w.invalidateSession(postSignIn)).Methods("POST")

	postSignUp := http.HandlerFunc(w.accountApi.PostSignUp)
	account.Handle("/sign-up", w.invalidateSession(postSignUp)).Methods("POST")

	postAccount := http.HandlerFunc(w.accountApi.PostAccount)
	getAccount := http.HandlerFunc(w.accountApi.GetAccount)
	account.Handle("", w.mustBeSignedIn(postAccount)).Methods("POST")
	account.Handle("", w.mustBeSignedIn(getAccount)).Methods("GET")

	product.HandleFunc("/catalogue", w.productApi.Search).Methods("GET")
	product.HandleFunc("/categories", w.productApi.Categories).Methods("GET")
	product.HandleFunc("/deals", w.productApi.Deals).Methods("GET")

	order.HandleFunc("/basket", w.orderApi.PostBasket).Methods("POST")
	order.HandleFunc("/basket", w.orderApi.GetBasket).Methods("GET")
	order.HandleFunc("/checkout", w.orderApi.PostCheckout).Methods("POST")

	getHistory := http.HandlerFunc(w.orderApi.GetHistory)
	order.Handle("/history", w.mustBeSignedIn(getHistory)).Methods("POST")

	// We don't use mustBeSignedIn here.
	// This allows guest customers with an order ID from a 'track my order'
	// email to fetch the order. What problems could there be with this?
	// What properties does order ID need to be secure. OWASP Top 10 #3 and #5
	order.HandleFunc("/{id}", w.orderApi.GetOrder).Methods("GET")

	w.router.PathPrefix("/").HandlerFunc(w.error404Handler)

	fmt.Println("Paths created")
}

// TEST \/\/\/\/
func (w *Wiring) Values(writer http.ResponseWriter, request *http.Request) {
	session, _ := w.store.Get(request, "session-name")
	fmt.Println()
	fmt.Println("Current session: ", session.Values)
	fmt.Println()
}

//TEST /\/\/\/\

func (w *Wiring) mustBeSignedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		session, _ := w.store.Get(request, "session-name")

		if _, ok := session.Values["customerId"]; ok {
			next.ServeHTTP(writer, request)
		} else {
			response := struct {
				Error string
				Msg   string
			}{
				Error: "forbidden",
				Msg:   "user is not signed in",
			}

			json.NewEncoder(writer).Encode(response)
			writer.WriteHeader(http.StatusUnauthorized)
		}
	})
}

func (w *Wiring) invalidateSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		session, _ := w.store.Get(request, "session-name")

		for i := range session.Values {
			delete(session.Values, i)
		}

		session.Save(request, writer)

		next.ServeHTTP(writer, request)
	})
}

func (w *Wiring) error404Handler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html")

	fmt.Fprint(writer, "<h1>Uh Oh, route not found</h1>")
	fmt.Fprint(writer, "<h2>Available routes:</h2>")
	fmt.Fprint(writer, "<ul>")
	w.printEndpoints(w.router, writer)
	fmt.Fprint(writer, "</ul>")
}

func (w *Wiring) printEndpoints(r *mux.Router, writer http.ResponseWriter) {
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
	store      *sessions.CookieStore
	router     *mux.Router
	accountApi *AccountApi
	productApi *ProductApi
	orderApi   *OrderApi
}
