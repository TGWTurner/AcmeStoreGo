package main

import (
	"bjssStoreGo/tests"
	"fmt"
	"os"
)

func main() {
	setUp("sql")
	fmt.Println("Sql")
	tests.RunDbTests()

	fmt.Println()

	setUp("memory")
	fmt.Println("Memory")
	tests.RunDbTests()
}

func setUp(mode string) {
	//"sql" or "sql-mem" or ""
	os.Setenv("DB_CONNECTION", mode)
}
