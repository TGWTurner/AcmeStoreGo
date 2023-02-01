package blTests

import (
	"bjssStoreGo/backend/layers/businessLogic"
	"bjssStoreGo/backend/layers/dataAccess"
	"fmt"
)

func SetUpAccount() businessLogic.AccountService {
	db := dataAccess.InitiateConnection()
	accountService := businessLogic.NewAccountService(db.Account)

	return *accountService
}

func RunTests() {
	// fmt.Println("OrderTests:")
	// successes, outOf := runTestSet(orderBlTests)
	// fmt.Println(successes, "/", outOf)

	fmt.Println("AccountTests:")
	successes, outOf := runTestSet(accountBlTests)
	fmt.Println(successes, "/", outOf)

	// fmt.Println("ProductTests:")
	// successes, outOf = runTestSet(productBlTests)
	// fmt.Println(successes, "/", outOf)
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

var accountBlTests = []func() (bool, string, string){
	TestCanSignInWithExistingAccount,
	TestCannotSignInWithIncorrectPassword,
	TestCannotSignInWithNonRegisteredEmail,
	TestCanSignUpWithNewAccount,
	TestCannotSignUpWithAlreadyRegisteredAccount,
	TestCanUpdateExistingAccount,
	TestCannotUpdateUnregisteredAccount,
	TestCanGetAccountById,
	TestCannotGetAccountWithFakeId,
}
var orderBlTests = []func() (bool, string, string){}

var productBlTests = []func() (bool, string, string){}
