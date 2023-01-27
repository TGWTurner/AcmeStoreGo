package main

import (
	"bjssStoreGo/tests"
	"os"
)

func main() {
	setUp()

	//Order Tests
	tests.TestCreateOrder()
	tests.TestGetOrderByToken()
	tests.TestGetOrdersByCustomerId()
	//Product Tests
	tests.TestGetProductGivenId()
	tests.TestGetProductsGivenIds()
	tests.TestGetCategoriesReturnsCorrectCategories()
	tests.TestGetProductsByCategoryProvidesCorrectProducts()
	tests.TestGetProductsByText()
	tests.TestGetProductsWithCurrentDeals()
	tests.TestDecreaseStockReducesStockByCorrectQuantity()
	//Account Tests
	tests.TestCreateAccount()
	tests.TestGetAccountByEmail()
	tests.TestGetById()
	tests.TestUpdateAccount()
}

func setUp() {
	os.Setenv("DB_CONNECTION", "sql")
}
