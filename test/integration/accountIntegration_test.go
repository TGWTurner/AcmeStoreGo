package integration

import (
	td "backend/layers/dataAccess/testData"
	"backend/test"
	"backend/utils"
	"encoding/json"
	"reflect"
	"testing"
)

func AssertSignedIn(t *testing.T, ar *test.ApiRequester) {
	response := ar.Get("/api/account")

	test.AssertResponseCode(t, 200, response.Code)
}

func AssertNotSignedIn(t *testing.T, ar *test.ApiRequester) {
	response := ar.Get("/api/account")

	test.AssertResponseCode(t, 401, response.Code)
}

func TestNotSignedInByDefault(t *testing.T) {
	w := test.SetUpApi(t)
	apiRequester := test.NewApiRequester(w)
	defer w.Close()

	AssertNotSignedIn(t, apiRequester)
}

func TestSignsInUsingPrePopulatedTestAccount(t *testing.T) {
	w := test.SetUpApi(t)
	apiRequester := test.NewApiRequester(w)
	defer w.Close()

	body, err := json.Marshal(td.GetTestAccountCredentials())

	test.AssertNil(t, err)

	response := apiRequester.Post("/api/account/sign-in", body)

	test.AssertResponseCode(t, 200, response.Code)

	AssertSignedIn(t, apiRequester)
}

func TestSignInReturnsTheRightAccount(t *testing.T) {
	w := test.SetUpApi(t)
	apiRequester := test.NewApiRequester(w)
	defer w.Close()

	body, err := json.Marshal(td.GetTestAccountCredentials())

	test.AssertNil(t, err)

	response := apiRequester.Post("/api/account/sign-in", body)

	test.AssertResponseCode(t, 200, response.Code)

	var actual utils.AccountApiResponse
	err = json.NewDecoder(response.Body).Decode(&actual)

	test.AssertNil(t, err)

	expected := td.GetAccountTestData()[0].OmitPasswordHash()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expeced account: %v, got account: %v", expected, actual)
	}
}

func TestSignUp(t *testing.T) {
	w := test.SetUpApi(t)
	apiRequester := test.NewApiRequester(w)
	defer w.Close()

	newUser := makeTestUser("signup-basic@example.com")

	body, err := json.Marshal(newUser)

	test.AssertNil(t, err)

	response := apiRequester.Post("/api/account/sign-up", body)

	test.AssertResponseCode(t, 200, response.Code)
}

func TestSignupReturnsTheNewlyCreatedAccount(t *testing.T) {
	w := test.SetUpApi(t)
	apiRequester := test.NewApiRequester(w)
	defer w.Close()

	newUser := makeTestUser("signup-return@example.com")

	account := signUp(t, apiRequester, newUser)

	if account.Id == 0 {
		t.Errorf("Expected non 0 account id")
	}

	account.Id = 0

	expected := utils.AccountApiResponse{
		Id: 0,
		ShippingDetails: utils.ShippingDetails{
			Email:    newUser.Email,
			Name:     newUser.Name,
			Address:  newUser.Address,
			Postcode: newUser.Postcode,
		},
	}

	if !reflect.DeepEqual(expected, account) {
		t.Errorf("Expected account %v, got account %v", expected, account)
	}
}

func TestIsSignedInAfterSignUp(t *testing.T) {
	w := test.SetUpApi(t)
	apiRequester := test.NewApiRequester(w)
	defer w.Close()

	newUser := makeTestUser("signup-return@example.com")

	signUp(t, apiRequester, newUser)

	AssertSignedIn(t, apiRequester)
}

func TestCanSignInWithANewlyCreatedUser(t *testing.T) {
	w := test.SetUpApi(t)
	apiRequester := test.NewApiRequester(w)
	defer w.Close()

	newUser := makeTestUser("signup-issignedin@example.com")

	signUp(t, apiRequester, newUser)

	apiRequester2 := test.NewApiRequester(w)

	credentials := struct {
		Email    string
		Password string
	}{
		Email:    "signup-issignedin@example.com",
		Password: "secret",
	}

	body, err := json.Marshal(credentials)

	test.AssertNil(t, err)

	response := apiRequester2.Post("/api/account/sign-in", body)

	test.AssertResponseCode(t, 200, response.Code)
	AssertSignedIn(t, apiRequester2)
}

func TestInvalidCredentialsRejected(t *testing.T) {
	w := test.SetUpApi(t)
	apiRequester := test.NewApiRequester(w)
	defer w.Close()

	invalidCredentials := struct {
		Email    string
		Password string
	}{
		Email:    "nosuch@example.com",
		Password: "password",
	}

	body, err := json.Marshal(invalidCredentials)

	test.AssertNil(t, err)

	response := apiRequester.Post("/api/account/sign-in", body)

	test.AssertResponseCode(t, 401, response.Code)

	AssertNotSignedIn(t, apiRequester)
}

func TestRetrievesAnAccount(t *testing.T) {
	w := test.SetUpApi(t)
	apiRequester := test.NewApiRequester(w)
	defer w.Close()

	original := test.SignInPrePopulatedUser(t, apiRequester)

	response := apiRequester.Get("/api/account")

	test.AssertResponseCode(t, 200, response.Code)

	var actual utils.AccountApiResponse
	err := json.NewDecoder(response.Body).Decode(&actual)

	test.AssertNil(t, err)

	if !reflect.DeepEqual(original, actual) {
		t.Errorf("Expected account: %v, got account: %v", original, actual)
	}
}

func TestUpdatesAccount(t *testing.T) {
	w := test.SetUpApi(t)
	apiRequester := test.NewApiRequester(w)
	defer w.Close()

	initial := signUp(t, apiRequester, makeTestUser("updateaccount@example.com"))

	toUpdate := utils.UpdateAccount{
		Id:              initial.Id,
		ShippingDetails: initial.ShippingDetails,
	}
	toUpdate.Name = "changed"

	body, err := json.Marshal(toUpdate)

	test.AssertNil(t, err)

	response := apiRequester.Post("/api/account", body)

	test.AssertResponseCode(t, 200, response.Code)

	var expected utils.AccountApiResponse
	err = json.NewDecoder(response.Body).Decode(&expected)

	test.AssertNil(t, err)

	initial.Name = "changed"

	if !reflect.DeepEqual(initial, expected) {
		t.Errorf("Expected returned account: %v, got account: %v", initial, response)
	}

	fetched := apiRequester.Get("/api/account")

	test.AssertResponseCode(t, 200, fetched.Code)

	var actual utils.AccountApiResponse
	err = json.NewDecoder(fetched.Body).Decode(&actual)

	test.AssertNil(t, err)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected response account: %v, got: %v", expected, actual)
	}
}

func TestHandlesPasswordChange(t *testing.T) {
	w := test.SetUpApi(t)
	apiRequester := test.NewApiRequester(w)
	defer w.Close()

	newUser := makeTestUser("passwordchange@example.com")
	returnedUser := signUp(t, apiRequester, newUser)

	toUpdate := utils.UpdateAccount{
		Id:              returnedUser.Id,
		ShippingDetails: returnedUser.ShippingDetails,
		Password:        "changed-password",
	}

	body, err := json.Marshal(toUpdate)

	test.AssertNil(t, err)

	modified := apiRequester.Post("/api/account", body)

	test.AssertResponseCode(t, 200, modified.Code)

	// check old password fails
	oldCredentials := struct {
		Email    string
		Password string
	}{Email: newUser.Email, Password: newUser.Password}

	apiRequester2 := test.NewApiRequester(w)

	body, err = json.Marshal(oldCredentials)

	test.AssertNil(t, err)

	invalid := apiRequester2.Post("/api/account/sign-in", body)

	test.AssertResponseCode(t, 401, invalid.Code)

	// check new password works
	updatedCredentials := struct {
		Email    string
		Password string
	}{Email: newUser.Email, Password: "changed-password"}

	body, err = json.Marshal(updatedCredentials)

	test.AssertNil(t, err)

	ok := apiRequester2.Post("/api/account/sign-in", body)

	test.AssertResponseCode(t, 200, ok.Code)
}

type account struct {
	Email    string
	Name     string
	Address  string
	Postcode string
	Password string
}

func makeTestUser(email string) account {
	return account{
		Email:    email,
		Name:     "a",
		Address:  "b",
		Postcode: "abc123",
		Password: "secret",
	}
}

func signUp(t *testing.T, ar *test.ApiRequester, newUser account) utils.AccountApiResponse {
	body, err := json.Marshal(newUser)

	//wants password not password hash

	test.AssertNil(t, err)

	response := ar.Post("/api/account/sign-up", body)

	test.AssertResponseCode(t, 200, response.Code)

	var newAccount utils.AccountApiResponse
	err = json.NewDecoder(response.Body).Decode(&newAccount)

	test.AssertNil(t, err)

	return newAccount
}
