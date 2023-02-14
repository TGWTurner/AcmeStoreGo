package integration

import (
	"bjssStoreGo/backend/layers/api"
	td "bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/test"
	"bjssStoreGo/backend/utils"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestNotSignedInByDefault(t *testing.T) {
	w := test.SetUpApi()
	defer w.Close()

	test.AssertNotSignedIn(t, w)
}

func TestSignsInUsingPrePopulatedTestAccount(t *testing.T) {
	w := test.SetUpApi()
	defer w.Close()

	body, err := json.Marshal(td.GetTestAccountCredentials())

	test.AssertNil(t, err)

	response := test.ApiRequest(
		t,
		w,
		"POST",
		"/api/account/sign-in",
		body,
	)

	fmt.Println("================")
	fmt.Println(response)
	fmt.Println("================")

	test.AssertResponseCode(t, 200, response.Code)

	test.AssertSignedIn(t, w)
}

func TestSignInReturnsTheRightAccount(t *testing.T) {
	w := test.SetUpApi()
	defer w.Close()

	body, err := json.Marshal(td.GetTestAccountCredentials())

	test.AssertNil(t, err)

	response := test.ApiRequest(
		t,
		w,
		"POST",
		"/api/account/sign-in",
		body,
	)

	test.AssertResponseCode(t, 200, response.Code)

	var actual []utils.AccountApiResponse
	err = json.NewDecoder(response.Body).Decode(&actual)

	test.AssertNil(t, err)

	expected := td.GetAccountTestData()[0].OmitPasswordHash()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expeced account: %v, got account: %v", expected, actual)
	}
}

func TestSignUp(t *testing.T) {
	w := test.SetUpApi()
	defer w.Close()

	newUser := makeTestUser("signup-basic@example.com")

	body, err := json.Marshal(newUser)

	test.AssertNil(t, err)

	response := test.ApiRequest(
		t,
		w,
		"POST",
		"/api/account/sign-up",
		body,
	)

	test.AssertResponseCode(t, 200, response.Code)
}

func TestSignupReturnsTheNewlyCreatedAccount(t *testing.T) {
	w := test.SetUpApi()
	defer w.Close()

	newUser := makeTestUser("signup-return@example.com")

	account := signUp(t, w, newUser)

	if account.Id == 0 {
		t.Errorf("Expected non 0 account id")
	}

	account.Id = 0

	expected := newUser.OmitPasswordHash()

	if !reflect.DeepEqual(expected, account) {
		t.Errorf("Expected account %v, got account %v", expected, account)
	}
}

func TestIsSignedInAfterSignUp(t *testing.T) {
	w := test.SetUpApi()
	defer w.Close()

	newUser := makeTestUser("signup-return@example.com")

	signUp(t, w, newUser)

	test.AssertSignedIn(t, w)
}

func TestCanSignInWithANewlyCreatedUser(t *testing.T) {
	w := test.SetUpApi()
	defer w.Close()

	newUser := makeTestUser("signup-issignedin@example.com")

	signUp(t, w, newUser)

	w2 := test.SetUpApi()
	defer w2.Close()

	credentials := struct {
		Email    string
		Password string
	}{
		Email:    newUser.Email,
		Password: newUser.PasswordHash,
	}

	body, err := json.Marshal(credentials)

	test.AssertNil(t, err)

	response := test.ApiRequest(
		t,
		w2,
		"GET",
		"/api/account/sign-in",
		body,
	)

	test.AssertResponseCode(t, 200, response.Code)
	test.AssertSignedIn(t, w2)
}

func TestInvalidCredentialsRejected(t *testing.T) {
	w := test.SetUpApi()
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

	response := test.ApiRequest(
		t,
		w,
		"GET",
		"/api/account/sign-in",
		body,
	)

	test.AssertResponseCode(t, 401, response.Code)

	test.AssertSignedIn(t, w)
}

func TestRetrievesAnAccount(t *testing.T) {
	w := test.SetUpApi()
	defer w.Close()

	original := signInPrePopulatedUser(t, w)

	response := test.ApiRequest(
		t,
		w,
		"GET",
		"/api/account",
		nil,
	)

	test.AssertResponseCode(t, 200, response.Code)

	var actual []utils.AccountApiResponse
	err := json.NewDecoder(response.Body).Decode(&actual)

	test.AssertNil(t, err)

	if !reflect.DeepEqual(original, actual) {
		t.Errorf("Expected account: %v, got account: %v", original, actual)
	}
}

func TestUpdatesAccount(t *testing.T) {
	w := test.SetUpApi()
	defer w.Close()

	initial := signUp(t, w, makeTestUser("updateaccount@example.com"))

	toUpdate := utils.UpdateAccount{
		Id:              initial.Id,
		ShippingDetails: initial.ShippingDetails,
	}
	toUpdate.Name = "changed"

	body, err := json.Marshal(toUpdate)

	test.AssertNil(t, err)

	modified := test.ApiRequest(
		t,
		w,
		"POST",
		"/api/account",
		body,
	)

	test.AssertResponseCode(t, 200, modified.Code)

	initial.Name = "changed"

	if !reflect.DeepEqual(initial, modified) {
		t.Errorf("Expected returned account: %v, got account: %v", initial, modified)
	}

	fetched := test.ApiRequest(
		t,
		w,
		"GET",
		"/api/account",
		nil,
	)

	test.AssertResponseCode(t, 200, fetched.Code)

	if !reflect.DeepEqual(fetched.Body, modified.Body) {
		t.Errorf("Expected response body: %v, got: %v", modified.Body, fetched.Body)
	}
}

func TestHandlesPasswordChange(t *testing.T) {
	w := test.SetUpApi()
	defer w.Close()

	newUser := makeTestUser("passwordchange@example.com")
	signUp(t, w, newUser)

	toUpdate := utils.UpdateAccount{
		Id:              newUser.Id,
		ShippingDetails: newUser.ShippingDetails,
		Password:        "changed-password",
	}

	body, err := json.Marshal(toUpdate)

	test.AssertNil(t, err)

	modified := test.ApiRequest(
		t,
		w,
		"POST",
		"/api/account",
		body,
	)

	test.AssertResponseCode(t, 200, modified.Code)

	// check old password fails
	oldCredentials := struct {
		Email    string
		Password string
	}{Email: newUser.Email, Password: newUser.PasswordHash}

	w2 := test.SetUpApi()
	defer w2.Close()

	body, err = json.Marshal(oldCredentials)

	test.AssertNil(t, err)

	invalid := test.ApiRequest(
		t,
		w2,
		"POST",
		"/api/account/sign-in",
		body,
	)

	test.AssertResponseCode(t, 401, invalid.Code)

	// check new password works
	updatedCredentials := struct {
		Email    string
		Password string
	}{Email: newUser.Email, Password: "changed-password"}

	w3 := test.SetUpApi()
	defer w3.Close()

	body, err = json.Marshal(updatedCredentials)

	test.AssertNil(t, err)

	ok := test.ApiRequest(
		t,
		w2,
		"POST",
		"/api/account/sign-in",
		body,
	)

	test.AssertResponseCode(t, 200, ok.Code)
}

func signInPrePopulatedUser(t *testing.T, w *api.Wiring, index ...int) utils.AccountApiResponse {
	var credentials struct {
		Email    string
		Password string
	}

	if len(index) == 0 {
		credentials = td.GetTestAccountCredentials(0)
	} else {
		credentials = td.GetTestAccountCredentials(index[0])
	}

	method := "POST"
	path := "/api/accont/sign-in"
	body, err := json.Marshal(credentials)

	test.AssertNil(t, err)

	response := test.ApiRequest(
		t,
		w,
		method,
		path,
		body,
	)

	var newAccount utils.AccountApiResponse
	err = json.NewDecoder(response.Body).Decode(&newAccount)

	test.AssertNil(t, err)

	return newAccount
}

func makeTestUser(email string) utils.Account {
	return utils.Account{
		ShippingDetails: utils.ShippingDetails{
			Email:    email,
			Name:     "a",
			Address:  "b",
			Postcode: "abc123",
		},
		PasswordHash: "secret",
	}
}

func signUp(t *testing.T, w *api.Wiring, newUser utils.Account) utils.AccountApiResponse {
	method := "GET"
	path := "/api/account/sign-up"
	body, err := json.Marshal(newUser)

	test.AssertNil(t, err)

	response := test.ApiRequest(
		t,
		w,
		method,
		path,
		body,
	)

	test.AssertResponseCode(t, 200, response.Code)

	var newAccount utils.AccountApiResponse
	err = json.NewDecoder(response.Body).Decode(&newAccount)

	test.AssertNil(t, err)

	return newAccount
}
