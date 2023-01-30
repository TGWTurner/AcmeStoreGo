package main

import (
	"bjssStoreGo/tests"
	"fmt"
	"os"
)

var orderTests = []func() (bool, string, string){
	tests.TestCreateOrder,
	tests.TestGetOrderByToken,
	tests.TestGetOrdersByCustomerId,
}
var productTests = []func() (bool, string, string){
	tests.TestGetProductGivenId,
	tests.TestGetProductsGivenIds,
	tests.TestGetCategoriesReturnsCorrectCategories,
	tests.TestGetProductsByCategoryProvidesCorrectProducts,
	tests.TestGetProductsByText,
	tests.TestGetProductsWithCurrentDeals,
	tests.TestDecreaseStockReducesStockByCorrectQuantity,
}

var accountTests = []func() (bool, string, string){
	tests.TestCreateAccountWithNewAccountPasses,
	tests.TestCreateAccountWithExistingAccountFails,
	tests.TestGetAccountByEmail,
	tests.TestGetAccountByNonExistingEmailFails,
	tests.TestGetByIdSucceedsForExistingAccount,
	tests.TestGetByIdFailsForNonExistingId,
	tests.TestUpdateAccountSucceedsForExistingAccount,
	tests.TestUpdateAccountFailsForNonExistingAccount,
}

func main() {
	setUp("sql")
	fmt.Println("Sql")
	runTests()

	fmt.Println()

	setUp("memory")
	fmt.Println("Memory")
	runTests()
}

func runTests() {
	fmt.Println("OrderTests:")
	successes, outOf := runTestSet(orderTests)
	fmt.Println(successes, "/", outOf)

	fmt.Println("AccountTests:")
	successes, outOf = runTestSet(accountTests)
	fmt.Println(successes, "/", outOf)

	fmt.Println("ProductTests:")
	successes, outOf = runTestSet(productTests)
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

		PrintTestResult(false, name, message)
	}

	return successes, len(testSet)
}

func setUp(mode string) {
	//"sql" or "sql-mem" or ""
	os.Setenv("DB_CONNECTION", mode)
}

func PrintTestResult(pass bool, testName string, message string) {
	fmt.Printf("\tTEST ")
	if pass {
		fmt.Print("PASSED")
	} else {
		fmt.Print("FAILED")
	}
	fmt.Println(" -- " + testName)
	fmt.Printf("\t\t" + message + "\n")
}
