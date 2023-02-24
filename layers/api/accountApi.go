package api

import (
	bl "backend/layers/businessLogic"
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
// @Failure 401 {object} ApiErrorResponse "Invalid credentials"
// @Router /api/account/sign-in [post]
func (a *AccountApi) PostSignIn(w http.ResponseWriter, r *http.Request) {
	// To implement
}

// PostSignUp godoc
// @Summary Registers a new user
// @Description Signs up, deletes any existing session, creates a new one for this user. Will give an error if the user already exists.
// @ID PostSignUp
// @Accept json
// @Produce json
// @Param user body AccountDetails true "account information"
// @Success 200 {object} utils.AccountApiResponse "The user's account"
// @Failure 400 {object} ApiErrorResponse "Malformed request or account already exists"
// @Router /api/account/sign-up [post]
func (a *AccountApi) PostSignUp(w http.ResponseWriter, r *http.Request) {
	// To implement
}

// GetAccount godoc
// @Summary Gets the user's Account
// @Description
// @ID GetAccount
// @Accept json
// @Produce json
// @Success 200 {object} utils.AccountApiResponse "The user's account"
// @Failure 401 {object} ApiErrorResponse "User is not signed in"
// @Router /api/account [get]
func (a *AccountApi) GetAccount(w http.ResponseWriter, r *http.Request) {
	// To implement
}

// PostAccount godoc
// @Summary Updates the users account
// @Description
// @ID PostAccount
// @Accept json
// @Produce json
// @Param user body AccountDetails true "account information"
// @Success 200 {object} utils.AccountApiResponse "The user's account"
// @Failure 401 {object} ApiErrorResponse "User is not signed in"
// @Router /api/account [post]
func (a *AccountApi) PostAccount(w http.ResponseWriter, r *http.Request) {
	// To implement
}

type AccountApi struct {
	as bl.AccountService
	s  *sessions.CookieStore
}
