package main

import "bjssStoreGo/tests"

func main() {
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

/*
TESTING:
	- Accounts -
	1. Add Account
	2. Get Account by email
	3. Get Account by id
	4. Update an Account
*/
