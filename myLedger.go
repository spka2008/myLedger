package main

import (
	"fmt"
	"os"
)

func main() {
	db := DBInit()
	tr := newTransactionGetJSON(os.Args[1],db)
	for _, el := range tr.ToStrings() {
		fmt.Println(el)
	}
	var answer string
	fmt.Print("Записать?y/n")
	fmt.Scanln(&answer)
	if answer == "y" {
		tr.SaveTransaction()
	}
}
