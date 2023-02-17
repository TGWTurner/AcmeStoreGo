package test

import (
	"bjssStoreGo/backend/layers/api"
	da "bjssStoreGo/backend/layers/dataAccess"
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
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

func AssertResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code: %d, got response code: %d", expected, actual)
	}
}

func AssertProductSetsMatch(t *testing.T, expected, actual []utils.Product) {
	if len(actual) != len(expected) {
		t.Errorf("Expected %d products, got %d products", len(expected), len(actual))
	}

	for _, actualProduct := range actual {
		found := false
		for _, expectedProduct := range expected {
			if reflect.DeepEqual(expectedProduct, actualProduct) {
				found = true
				break
			}
		}
		if found == false {
			t.Errorf("Returned product not found in expected products")
		}
	}
}

func AssertCategorySetsMatch(t *testing.T, expected, actual []utils.ProductCategory) {
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

func AssertSignedIn(t *testing.T, ar *ApiRequester) {
	response := ar.Get("/api/account")

	AssertResponseCode(t, 200, response.Code)
}

func AssertNotSignedIn(t *testing.T, ar *ApiRequester) {
	response := ar.Get("/api/account")

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
