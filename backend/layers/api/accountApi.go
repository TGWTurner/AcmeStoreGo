package api

import (
	bl "bjssStoreGo/backend/layers/businessLogic"
	"bjssStoreGo/backend/utils"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

func NewAccountApi(accountService *bl.AccountService, s *sessions.CookieStore) *AccountApi {
	return &AccountApi{
		as: *accountService,
		s:  s,
	}
}

func (a *AccountApi) Close() {
	a.as.Close()
}

func (a *AccountApi) getSignedInUserId(r *http.Request) (int, error) {
	session, _ := a.s.Get(r, "session-name")

	customerId, ok := session.Values["customerId"]

	if !ok {
		return 0, errors.New("Failed to get customerId from session")
	}

	return customerId.(int), nil
}

func (a *AccountApi) setSignedInUserId(w http.ResponseWriter, r *http.Request, customerId int) {
	session, _ := a.s.Get(r, "session-name")
	session.Values["customerId"] = customerId
	session.Save(r, w)
}

func (a *AccountApi) PostSignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var eap struct {
		Email    string
		Password string
	}

	err := json.NewDecoder(r.Body).Decode(&eap)

	if err != nil {
		//log the error?
	}

	account, err := a.as.SignIn(eap.Email, eap.Password)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(utils.ApiErrorResponse{
			Error: "forbidden",
			Msg:   "Invalid credentials",
		})
		return
	}

	a.setSignedInUserId(w, r, account.Id)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

func (a *AccountApi) PostSignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var acc struct {
		Email    string
		Name     string
		Address  string
		Postcode string
		Password string
	}

	err := json.NewDecoder(r.Body).Decode(&acc)

	if err != nil {
		json.NewEncoder(w).Encode(utils.ApiErrorResponse{
			Error: "",
			Msg:   "Malformed request or account already exists",
		})
		w.WriteHeader(http.StatusBadRequest)
	}

	account, err := a.as.SignUp(acc.Email, acc.Password, acc.Name, acc.Address, acc.Postcode)

	if err != nil {
		json.NewEncoder(w).Encode(utils.ApiErrorResponse{
			Error: "",
			Msg:   "Malformed request or account already exists",
		})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

func (a *AccountApi) GetAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, err := a.getSignedInUserId(r)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(utils.ApiErrorResponse{
			Error: "forbidden",
			Msg:   "user is not signed in",
		})
		return
	}

	account, err := a.as.GetById(userId)

	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

func (a *AccountApi) PostAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, err := a.getSignedInUserId(r)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	var acc utils.UpdateAccount

	acc.Id = userId

	err = json.NewDecoder(r.Body).Decode(&acc)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(utils.ApiErrorResponse{
			Error: "forbidden",
			Msg:   "user is not signed in",
		})
		return
	}

	account, err := a.as.Update(acc)

	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

type AccountApi struct {
	as bl.AccountService
	s  *sessions.CookieStore
}
