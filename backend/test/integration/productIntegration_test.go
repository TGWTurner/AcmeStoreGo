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
	defer w.Close()

	method := "GET"
	path := "/api/product/catalogue"

	response := test.ApiRequest(
		t,
		w,
		method,
		path,
		nil,
	)

	test.AssertResponseCode(t, response.Code, 200)

	expected := td.GetProductTestData().Products
	var actual []utils.Product

	err := json.NewDecoder(response.Body).Decode(&actual)

	test.AssertNil(t, err)

	test.AssertProductSetsMatch(t, actual, expected)
}

func TestListsDeals(t *testing.T) {
	w := test.SetUpApi()
	defer w.Close()

	method := "GET"
	path := "/api/product/deals"

	response := test.ApiRequest(
		t,
		w,
		method,
		path,
		nil,
	)

	test.AssertResponseCode(t, response.Code, 200)

	currentDate := utils.GetFormattedDate()
	expected := test.GetTestProductsWithCurrentDeals(currentDate)

	var actual []utils.Product
	err := json.NewDecoder(response.Body).Decode(&actual)

	test.AssertNil(t, err)

	test.AssertProductSetsMatch(t, actual, expected)
}

func TestGetsCategories(t *testing.T) {
	w := test.SetUpApi()
	defer w.Close()

	method := "GET"
	path := "/api/product/categories"

	response := test.ApiRequest(
		t,
		w,
		method,
		path,
		nil,
	)

	test.AssertResponseCode(t, response.Code, 200)

	var actual []utils.ProductCategory
	err := json.NewDecoder(response.Body).Decode(&actual)

	test.AssertNil(t, err)

	expected := td.GetProductTestData().Categories

	test.AssertCategorySetsMatch(t, actual, expected)
}

func TestListsProductsInASingleCategory(t *testing.T) {
	w := test.SetUpApi()
	defer w.Close()

	categoryId := 1

	method := "GET"
	path := "/api/product/catalogue?category=" + strconv.Itoa(categoryId)

	response := test.ApiRequest(
		t,
		w,
		method,
		path,
		nil,
	)

	test.AssertResponseCode(t, response.Code, 200)

	var actual []utils.Product
	err := json.NewDecoder(response.Body).Decode(&actual)

	test.AssertNil(t, err)

	expected := test.GetTestProductsByCategory(categoryId)

	test.AssertProductSetsMatch(t, actual, expected)
}
