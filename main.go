package main

import (
	"fmt"
	//"io"
	"os"
//	"encoding/json"
)

func main() {
	tr := newTransactionGetJSON(os.Args[1])
	for _, el := range tr.ToStrings() {
		fmt.Println(el)
	}
}
