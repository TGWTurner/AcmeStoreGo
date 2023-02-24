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

// GetBasket godoc
// @Summary Gets the user's Basket
// @Description The same session cookie that created the basket is needed
// @ID GetBasket
// @Produce json
// @Success 200 {object} utils.Basket "A Basket"
// @Router /api/order/basket [get]
func (o *OrderApi) GetBasket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	basket := o.getBasket(r)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(basket)
}

//TODO remove Total from post

// PostBasket godoc
// @Summary Creates or updates the user's Basket
// @Description Sets a session cookie which is needed to later get the basket
// @ID PostBasket
// @Accept json
// @Produce json
// @Param Basket body utils.Basket true "A Basket"
// @Success 200 {object} utils.Basket "A Basket"
// @Router /api/order/basket [post]
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

// GetHistory godoc
// @Summary Gets the user's Order history
// @ID GetHistory
// @Produce json
// @Success 200 {object} utils.Basket "A Basket"
// @Router /api/order/history [get]
func (o *OrderApi) GetHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	customerId := GetSignedInUserId(r, o.s)

	if customerId == 0 {
		Error(w, r, http.StatusUnauthorized, "error", "User is not signed in")
		return
	}

	orderHistory, err := o.os.GetOrdersByCustomerId(customerId)

	if err != nil {
		Error(w, r, http.StatusInternalServerError, "error", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orderHistory)
}

// GetOrder godoc
// @Summary Fetches an order given a token
// @Description Does not require a signed in user so that we can implement getting an order via a link in an email, etc.
// @ID GetOrder
// @Accept json
// @Produce json
// @Param token path string true "Order token. Currently same as order.id"
// @Success 200 {object} OrderApiResponse "The newly created order"
// @Failure 404 {object} ApiErrorResponse "No such order"
// @Router /api/order/{token} [get]
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

	response := convertOrderToOrderApiResponse(order)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// PostCheckout godoc
// @Summary Creates a new order
// @Description Checks the stock levels and paymentToken. If Ok creates a new order. If not gives an error and the products there is not enough stock for. Sets a session cookie which can be used later to tie this order to a signed in user. Does not require a signed in user so guests can check out
// @ID PostCheckout
// @Accept json
// @Produce json
// @Param order body OrderRequest true "An order"
// @Success 200 {object} OrderApiResponse "The newly created order"
// @Failure 400 {object} ApiErrorResponse "An error. If the request was well formed this will be payment or stock level error. If stock level error, the quantityRemaining is returned for products with not enough stock."
// @Router /api/order/checkout [post]
func (o *OrderApi) PostCheckout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var order OrderRequest

	err := json.NewDecoder(r.Body).Decode(&order)

	if err != nil {
		Error(w, r, http.StatusInternalServerError, "error", err.Error())
	}

	customerId := GetSignedInUserId(r, o.s)

	newOrder, err := o.os.CreateOrder(customerId, order.ShippingDetails, order.Items)

	if err != nil {
		Error(w, r, http.StatusInternalServerError, "error", err.Error())
		return
	}

	response := convertOrderToOrderApiResponse(newOrder)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func convertOrderToOrderApiResponse(order utils.Order) OrderApiResponse {
	return OrderApiResponse{
		Id:              order.Id,
		Total:           order.Total,
		UpdatedDate:     order.UpdatedDate,
		ShippingDetails: order.ShippingDetails,
		Items:           order.Items,
	}
}

type OrderApi struct {
	os bl.OrderService
	s  *sessions.CookieStore
}
