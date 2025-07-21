package main

import (
	"fmt"
	"os"

	"github.com/spka2008/myLedger/db"
	"github.com/spka2008/myLedger/item"
	"github.com/spka2008/myLedger/receipt"
)

func main() {
	db := db.DataBaseConnect()
	receipt := receipt.NewCheck(os.Args[1])
	answer := "n"
	var tr *item.Transaction
	for answer != "y" {
		tr = item.CheckToTransaction(receipt, db.Path)
		tr.Print()
		fmt.Print("Записать?y/n")
		fmt.Scanln(&answer)
	}
	tr.SaveTransaction()
	os.Remove(os.Args[1])

}
