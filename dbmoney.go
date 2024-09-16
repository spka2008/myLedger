package main

import (
	"fmt"
//	"bufio"
	"encoding/json"
	"os"
	"os/exec"
	//"github.com/jcmuller/dmenu"
)

type DB struct {
	Path string
	Connect bool 
}

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
	database DB
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

func newTransactionGetJSON(path string, db DB) *Transaction {

	if !db.Connect {
		panic("соединение с базой не установлено")
	}
	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("file not open %v", err))
	}
	defer file.Close()
	var t Transaction
	t.database = db
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&t); err !=nil {
		panic(fmt.Sprintf("ERROR Serializ: %v", err))
	}
	for _, el := range t.Records {
		el.CollectComment()
	}
	t.getDiscription()
	return &t
}

func  DBInit() DB {
	var db DB
	path, exist := os.LookupEnv("LEDGERPATH")
	if !exist {
		panic("нет пути к базе")
	}
	db.Path = path
	db.Connect = true
	return db
}

func (tr *Transaction) getDiscription() {
/*	file, err := os.Open(tr.database.Path + "payee.dat")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var payers []string
	s := bufio.NewScanner(file) 
	for s.Scan() {
		payers = append(payers, s.Text())
	}
	fmt.Println(payers)
//	var payer string */
	var flag string
	fmt.Print("Destination " + tr.Destination + "? y/n ")
	fmt.Scan(&flag)
	if flag == "n" {
		cmd := "cat " + tr.database.Path + "payee.dat | dmenu"
		out, _ := exec.Command("bash","-c",cmd).Output()
		tr.Destination = string(out)
	}
}

	
