package api

import (
	"bjssStoreGo/backend/layers/businessLogic"
	"bjssStoreGo/backend/utils"
	"encoding/json"
	"net/http"
)

// TODO: Implement Get signed in user id
func getSignedInUserId(req string) (int, error) {
	return 1, nil
}

// TODO: Implement Set signed in user id
func setSignedInUserId(req string, customerId string) {
	customerId = "userId"
}

func NewAccountApi(accountService *businessLogic.AccountService) *AccountApi {
	return &AccountApi{
		as: *accountService,
	}
}

func (a AccountApi) PostSignIn(w http.ResponseWriter, r *http.Request) {
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

func (a AccountApi) PostSignUp(w http.ResponseWriter, r *http.Request) {
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
		//log the error?
	}

	account, err := a.as.SignUp(acc.Email, acc.Password, acc.Name, acc.Address, acc.Postcode)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

func (a AccountApi) GetAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, err := getSignedInUserId("Session?")

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)

		response := struct {
			Error string
			Msg   string
		}{
			Error: "forbidden",
			Msg:   "user is not signed in",
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	account, err := a.as.GetById(userId)

	if err != nil {
		//log the error?
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

func (a AccountApi) PostAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_, err := getSignedInUserId("session?")

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	var acc utils.Account

	err = json.NewDecoder(r.Body).Decode(&acc)

	if err != nil {
		//log the error?
	}

	account, err := a.as.Update(acc)

	if err != nil {
		//log the error?
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

type AccountApi struct {
	as businessLogic.AccountService
}
