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

type UserDetails struct {
	Email    string
	Password string
}

// PostSignIn godoc
// @Summary Signs in
// @Description Signs in, deletes any existing session, creates a new one for this user.
// @ID PostSignIn
// @Accept json
// @Produce json
// @Param user body UserDetails true "account information"
// @Success 200 {object} utils.AccountApiResponse "The user's account"
// @Failure 401 {object} utils.ApiErrorResponse "Invalid credentials"
// @Router /api/account/sign-in [post]
func (a *AccountApi) PostSignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var userDetails UserDetails

	err := json.NewDecoder(r.Body).Decode(&userDetails)

	if err != nil {
		Error(w, r, http.StatusInternalServerError, "error", err.Error())
		return
	}

	account, err := a.as.SignIn(userDetails.Email, userDetails.Password)

	if err != nil {
		Error(w, r, http.StatusInternalServerError, "forbidden", "Invalid credentials")
		return
	}

	a.setSignedInUserId(w, r, account.Id)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

// PostSignUp godoc
// @Summary Registers a new user
// @Description Signs up, deletes any existing session, creates a new one for this user. Will give an error if the user already exists.
// @ID PostSignUp
// @Accept json
// @Produce json
// @Param user body utils.AccountDetails true "account information"
// @Success 200 {object} utils.AccountApiResponse "The user's account"
// @Failure 400 {object} utils.ApiErrorResponse "Malformed request or account already exists"
// @Router /api/account/sign-up [post]
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
		Error(w, r, http.StatusBadRequest, "error", "Malformed request or account already exists")
		return
	}

	account, err := a.as.SignUp(acc.Email, acc.Password, acc.Name, acc.Address, acc.Postcode)

	if err != nil {
		Error(w, r, http.StatusBadRequest, "error", "Malformed request or account already exists")
		return
	}

	a.setSignedInUserId(w, r, account.Id)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

// GetAccount godoc
// @Summary Gets the user's Account
// @Description
// @ID GetAccount
// @Accept json
// @Produce json
// @Param Cookie header string false "token"
// @Success 200 {object} utils.AccountApiResponse "The user's account"
// @Failure 401 {object} utils.ApiErrorResponse "User is not signed in"
// @Router /api/account [get]
func (a *AccountApi) GetAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, err := a.getSignedInUserId(r)

	if err != nil {
		Error(w, r, http.StatusUnauthorized, "error", "User is not signed in")
		return
	}

	account, err := a.as.GetById(userId)

	if err != nil {
		Error(w, r, http.StatusInternalServerError, "error", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

// PostAccount godoc
// @Summary Updates the users account
// @Description
// @ID PostAccount
// @Accept json
// @Produce json
// @Param user body utils.AccountDetails true "account information"
// @Success 200 {object} utils.AccountApiResponse "The user's account"
// @Failure 401 {object} utils.ApiErrorResponse "User is not signed in"
// @Router /api/account [post]
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
		Error(w, r, http.StatusUnauthorized, "error", "user is not signed in")
		return
	}

	account, err := a.as.Update(acc)

	if err != nil {
		Error(w, r, http.StatusInternalServerError, "error", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

type AccountApi struct {
	as bl.AccountService
	s  *sessions.CookieStore
}
