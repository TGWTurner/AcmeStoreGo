package dbTests

import (
	"bjssStoreGo/backend/layers/dataAccess"
	"bjssStoreGo/backend/utils"
	"fmt"
)

func SetUp() utils.Database {
	db := dataAccess.InitiateConnection()

	return db
}

func CloseDb(db utils.Database) {
	db.Order.Close()
	db.Account.Close()
	db.Product.Close()
}

func printTestResult(pass bool, testName string, message string) {
	fmt.Print("TEST ")
	if pass {
		fmt.Print("PASSED")
	} else {
		fmt.Print("FAILED")
	}
	fmt.Println(" -- " + testName)
	fmt.Printf("\t" + message + "\n")
}

func RunTests() {
	fmt.Println("OrderTests:")
	successes, outOf := runTestSet(orderDbTests)
	fmt.Println(successes, "/", outOf)

	fmt.Println("AccountTests:")
	successes, outOf = runTestSet(accountDbTests)
	fmt.Println(successes, "/", outOf)

	fmt.Println("ProductTests:")
	successes, outOf = runTestSet(productDbTests)
	fmt.Println(successes, "/", outOf)
}

func runTestSet(testSet []func() (bool, string, string)) (int, int) {
	successes := 0

	for _, test := range testSet {
		pass, name, message := test()

		if pass {
			successes++
			continue
		}

		printTestResult(false, name, message)
	}

	return successes, len(testSet)
}

var orderDbTests = []func() (bool, string, string){
	TestCreateOrder,
	TestSucceedsCreatingTwoIdenticalOrders,
	TestGetOrderByToken,
	TestFailsToGetOrderWithFakeToken,
	TestGetOrdersByCustomerId,
}

var productDbTests = []func() (bool, string, string){
	TestGetProductGivenId,
	TestGetProductsGivenIds,
	TestGetCategoriesReturnsCorrectCategories,
	TestGetProductsByCategoryProvidesCorrectProducts,
	TestGetProductsByText,
	TestGetProductsWithCurrentDeals,
	TestDecreaseStockReducesStockByCorrectQuantity,
	TestDecreaseStockFailsForFakeProduct,
}

var accountDbTests = []func() (bool, string, string){
	TestCreateAccountWithNewAccountPasses,
	TestCreateAccountWithExistingAccountFails,
	TestGetAccountByEmail,
	TestGetAccountByNonExistingEmailFails,
	TestGetByIdSucceedsForExistingAccount,
	TestGetByIdFailsForNonExistingId,
	TestUpdateAccountSucceedsForExistingAccount,
	TestUpdateAccountFailsForNonExistingAccount,
}
