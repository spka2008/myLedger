package main

import (
	"fmt"
	//"io"
	"encoding/json"
	"os"
)

type Record struct {
	account string
	Sum float64 `json:"sum"`
	Name string `json:"name"`
	Quantity float64 `json:"quantity"`
	Price float64 `json:"price"`
	comment string
}

type Transaction struct {
	Date string `json:"date"`
	status bool
	Destination string `json:"shopName"`
	Records []Record `json:"products"`
	TotalSum float64 `json:"totalSum"`
}

func (r *Record) CollectComment() {
	r.comment =fmt.Sprintf("|%.2f * %f|%v", r.Price, r.Quantity, r.Name)
}
func (r *Record) Format() string {
	l := 41 - len(fmt.Sprintf("%.2f", r.Sum))
	str := "    %-" + fmt.Sprint(l) + "s$%.2f"
	if (len(r.comment) != 0) {
		str += "  ;  " + r.comment
	}
	return fmt.Sprintf(str, r.account, r.Sum)
}

func (t *Transaction) CheckSum() bool {
	var sum float64 = 0.0
	for _, el := range t.Records {
		sum += el.Sum
	}
	return sum == 0
}

func (t *Transaction) ToStrings() []string {
	var res []string
	var status string
	if t.status {
		status = "*"
	} else {
		status = ""
	}
	res = append(res, t.Date + " " + status + t.Destination)
	for i := 0; i < len(t.Records); i++ {
		res = append(res, t.Records[i].Format())
	}
	return res
}

func newTransactionGetJSON(path string) *Transaction {

	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("file not open %v", err))
	}
	defer file.Close()
	var t Transaction
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&t); err !=nil {
		panic(fmt.Sprintf("ERROR Serializ: %v", err))
	}
	for _, el := range t.Records {
		el.CollectComment()
	}
	return &t
}


