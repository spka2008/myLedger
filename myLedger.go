package main

import (
	"fmt"
	//"io"
	"os"
//	"encoding/json"
)

func main() {
	db := DBInit()
	tr := newTransactionGetJSON(os.Args[1],db)
	for _, el := range tr.ToStrings() {
		fmt.Println(el)
	}
	tr.SaveTransaction()
}
