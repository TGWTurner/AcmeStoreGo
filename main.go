package main

import (
	"bjssStoreGo/blTests"
	"bjssStoreGo/dbTests"
	"fmt"
	"os"
)

func main() {
	runDbTests()
	runBlTests()
}

func runBlTests() {
	fmt.Println()

	setUp("sql")
	fmt.Println("Sql")
	blTests.RunTests()

	fmt.Println()

	setUp("memory")
	fmt.Println("Memory")
	blTests.RunTests()
}

func runDbTests() {
	fmt.Println()

	setUp("sql")
	fmt.Println("Sql")
	dbTests.RunTests()

	fmt.Println()

	setUp("memory")
	fmt.Println("Memory")
	dbTests.RunTests()
}

func setUp(mode string) {
	//"sql" or "sql-mem" or ""
	os.Setenv("DB_CONNECTION", mode)
}
