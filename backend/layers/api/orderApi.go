package api

import (
	bl "bjssStoreGo/backend/layers/businessLogic"
	"bjssStoreGo/backend/utils"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func NewOrderApi(orderService *bl.OrderService, s *sessions.CookieStore) *OrderApi {
	return &OrderApi{
		os: *orderService,
		s:  s,
	}
}

func (o *OrderApi) Close() {
	o.os.Close()
}

func (o *OrderApi) getBasket(r *http.Request) utils.Basket {
	session, _ := o.s.Get(r, "session-name")
	basket, ok := session.Values["basket"]

	if !ok {
		return utils.Basket{}
	}

	return basket.(utils.Basket)
}

func (o *OrderApi) setBasket(w http.ResponseWriter, r *http.Request, basket utils.Basket) {
	session, _ := o.s.Get(r, "session-name")
	session.Values["basket"] = basket

	session.Save(r, w)
}

func (o *OrderApi) GetBasket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	basket := o.getBasket(r)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(basket)
}

func (o *OrderApi) PostBasket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	currentBasket := o.getBasket(r)

	var basket utils.Basket

	err := json.NewDecoder(r.Body).Decode(&basket)

	if err != nil {
		Error(w, r, http.StatusInternalServerError, "error", err.Error())
	}

	newBasket, err := o.os.UpdateBasket(basket.Items, currentBasket)

	if err != nil {
		Error(w, r, http.StatusUnauthorized, "error", err.Error())
		return
	}

	o.setBasket(w, r, newBasket)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newBasket)
}

func (o *OrderApi) GetHistory(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement get history
}

func (o *OrderApi) GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	orderToken, ok := params["id"]

	if !ok {
		Error(w, r, http.StatusBadRequest, "error", "Missing id in request")
		return
	}

	order, err := o.os.GetOrderByToken(orderToken)

	if err != nil {
		Error(w, r, http.StatusNotFound, "error", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

type OrderRequest struct {
	PaymentToken    string
	ShippingDetails utils.ShippingDetails
	Items           []utils.OrderItem
}

func (o *OrderApi) PostCheckout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var order OrderRequest

	err := json.NewDecoder(r.Body).Decode(&order)

	if err != nil {
		Error(w, r, http.StatusInternalServerError, "error", err.Error())
	}

	session, _ := o.s.Get(r, "session-name")
	customerId := session.Values["customerId"]

	newOrder, err := o.os.CreateOrder(customerId.(int), order.ShippingDetails, order.Items)

	if err != nil {
		Error(w, r, http.StatusInternalServerError, "error", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newOrder)
}

type OrderApi struct {
	os bl.OrderService
	s  *sessions.CookieStore
}
