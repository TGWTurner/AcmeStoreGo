package test

import (
	"bjssStoreGo/backend/layers/api"
	da "bjssStoreGo/backend/layers/dataAccess"
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/exp/slices"
)

func SetUpApi() *api.Wiring {
	db := da.InitiateConnection()
	r := mux.NewRouter()
	store := sessions.NewCookieStore([]byte("my session encryption secret"))

	wiring := api.NewWiring(db, r, store)

	wiring.SetUpRoutes()

	return wiring
}

func AssertErrorString(t *testing.T, err error, msg string) {
	if err.Error() != msg {
		t.Errorf("got error %s, expected error %s", err.Error(), msg)
	}
}

func AssertNil(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Error was not nil, got: %s", err.Error())
	}
}

func ApiRequest(t *testing.T, wiring *api.Wiring, method string, path string, body []byte) *httptest.ResponseRecorder {
	requestBody := bytes.NewBuffer(body)

	req, err := http.NewRequest(method, path, requestBody)

	if err != nil {
		t.Errorf("Expected error: nil, got error: %s", err)
	}

	return executeRequest(req, wiring)
}

func executeRequest(req *http.Request, wiring *api.Wiring) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	wiring.Router.ServeHTTP(rr, req)

	return rr
}

func AssertResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code: %d, got response code: %d", expected, actual)
	}
}

func AssertProductSetsMatch(t *testing.T, actual, expected []utils.Product) {
	if len(actual) != len(expected) {
		t.Errorf("Expected %d products, got %d products", len(expected), len(actual))
	}

	for _, actualProduct := range actual {
		found := false
		for _, expectedProduct := range expected {
			if reflect.DeepEqual(actualProduct, expectedProduct) {
				found = true
				break
			}
		}
		if found == false {
			t.Errorf("Returned product not found in expected products")
		}
	}
}

func AssertCategorySetsMatch(t *testing.T, actual, expected []utils.ProductCategory) {
	if len(actual) != len(expected) {
		t.Errorf("Expected %d categories, got %d categories", len(expected), len(actual))
	}

	for _, category := range actual {
		found := false
		for _, expectedCategory := range expected {
			if reflect.DeepEqual(expectedCategory, category) {
				found = true
				break
			}
		}

		if found != true {
			t.Errorf("Returned category not found in expected categories")
		}
	}
}

func AssertSignedIn(t *testing.T, w *api.Wiring) {
	method := "GET"
	path := "/api/account"

	response := ApiRequest(
		t,
		w,
		method,
		path,
		nil,
	)

	fmt.Println("================")
	fmt.Println(response)
	fmt.Println("================")

	AssertResponseCode(t, 200, response.Code)
}

func AssertNotSignedIn(t *testing.T, w *api.Wiring) {
	method := "GET"
	path := "/api/account"

	response := ApiRequest(
		t,
		w,
		method,
		path,
		nil,
	)

	AssertResponseCode(t, 401, response.Code)
}

func GetTestProductById(id int) utils.Product {
	expectedProducts := testData.GetProductTestData().Products
	idx := slices.IndexFunc(
		expectedProducts,
		func(p utils.Product) bool { return p.Id == id },
	)

	return expectedProducts[idx]
}

func GetTestProductsByText(text string) []utils.Product {
	expectedProducts := testData.GetProductTestData().Products
	resProducts := []utils.Product{}

	for _, product := range expectedProducts {
		if strings.Contains(product.ShortDescription, text) || strings.Contains(product.LongDescription, text) {
			resProducts = append(resProducts, product)
		}
	}

	return resProducts
}

func GetTestProductsByCategory(categoryId int) []utils.Product {
	expectedProducts := testData.GetProductTestData().Products
	resProducts := []utils.Product{}

	for _, product := range expectedProducts {
		if categoryId == product.CategoryId {
			resProducts = append(resProducts, product)
		}
	}

	return resProducts
}

func GetTestProductsWithCurrentDeals(currentDate string) []utils.Product {
	deals := testData.GetProductTestData().Deals
	products := []utils.Product{}

	for _, deal := range deals {
		if deal.StartDate < currentDate && deal.EndDate > currentDate {
			products = append(products, GetTestProductById(deal.ProductId))
		}
	}

	return products
}
