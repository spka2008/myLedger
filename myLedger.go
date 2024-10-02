package main

import (
	"fmt"
	"os"
)

func main() {
	db := dataBaseConnect()
	tr := newTransactionGetJSON(os.Args[1], db)
	for _, el := range tr.toStrings() {
		fmt.Println(el)
	}
	var answer string
	fmt.Print("Записать?y/n")
	fmt.Scanln(&answer)
	if answer == "y" {
		tr.saveTransaction()
	}
}
