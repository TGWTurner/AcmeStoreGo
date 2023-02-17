package integration

import (
	td "bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/test"
	"bjssStoreGo/backend/utils"
	"encoding/json"
	"strconv"
	"testing"
)

func TestListsProducts(t *testing.T) {
	w := test.SetUpApi()
	apiRequester := test.NewApiRequester(w)
	defer w.Close()

	response := apiRequester.Get("/api/product/catalogue")

	test.AssertResponseCode(t, 200, response.Code)

	expected := td.GetProductTestData().Products
	var actual []utils.Product

	err := json.NewDecoder(response.Body).Decode(&actual)

	test.AssertNil(t, err)

	test.AssertProductSetsMatch(t, expected, actual)
}

func TestListsDeals(t *testing.T) {
	w := test.SetUpApi()
	apiRequester := test.NewApiRequester(w)
	defer w.Close()

	response := apiRequester.Get("/api/product/deals")

	test.AssertResponseCode(t, 200, response.Code)

	currentDate := utils.GetFormattedDate()
	expected := test.GetTestProductsWithCurrentDeals(currentDate)

	var actual []utils.Product
	err := json.NewDecoder(response.Body).Decode(&actual)

	test.AssertNil(t, err)

	test.AssertProductSetsMatch(t, expected, actual)
}

func TestGetsCategories(t *testing.T) {
	w := test.SetUpApi()
	apiRequester := test.NewApiRequester(w)
	defer w.Close()

	response := apiRequester.Get("/api/product/categories")

	test.AssertResponseCode(t, 200, response.Code)

	var actual []utils.ProductCategory
	err := json.NewDecoder(response.Body).Decode(&actual)

	test.AssertNil(t, err)

	expected := td.GetProductTestData().Categories

	test.AssertCategorySetsMatch(t, expected, actual)
}

func TestListsProductsInASingleCategory(t *testing.T) {
	w := test.SetUpApi()
	apiRequester := test.NewApiRequester(w)
	defer w.Close()

	categoryId := 1

	response := apiRequester.Get("/api/product/catalogue?category=" + strconv.Itoa(categoryId))

	test.AssertResponseCode(t, 200, response.Code)

	var actual []utils.Product
	err := json.NewDecoder(response.Body).Decode(&actual)

	test.AssertNil(t, err)

	expected := test.GetTestProductsByCategory(categoryId)

	test.AssertProductSetsMatch(t, expected, actual)
}

func TestListsProductsWithSearchTerm(t *testing.T) {
	w := test.SetUpApi()
	apiRequester := test.NewApiRequester(w)
	defer w.Close()

	searchTerm := "Apricot"

	response := apiRequester.Get("/api/product/catalogue?search=" + searchTerm)

	test.AssertResponseCode(t, 200, response.Code)

	var actual []utils.Product
	err := json.NewDecoder(response.Body).Decode(&actual)

	test.AssertNil(t, err)

	expected := test.GetTestProductsByText(searchTerm)

	test.AssertProductSetsMatch(t, expected, actual)
}
