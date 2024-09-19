package main

import (
	"fmt"
	"bufio"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
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
	r.comment =fmt.Sprintf("|%.2f * %.3f|%v", r.Price, r.Quantity, r.Name)
}

func (r *Record) Format() string {
	l := 41 - len(fmt.Sprintf("%.2f", r.Sum))
	str := "    %-" + fmt.Sprint(l) + "s$%.2f"
	if (len(r.comment) != 0) {
		str += "  ;  " + r.comment
	}
	return fmt.Sprintf(str, r.account, r.Sum)
}

func (t *Transaction) CheckSum() float64 {
	var sum float64 = 0.0
	for _, el := range t.Records {
		sum += el.Sum
	}
	return  t.TotalSum - sum
}

func (t *Transaction) ToStrings() []string {
	var res []string
	var status string
	if t.status {
		status = "* "
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
	for i, _  := range t.Records {
		t.Records[i].CollectComment()
	}
	t.getKeyboard()
	strSpl := strings.Split(t.Date, " ")[0]
	strDate := strings.Split(strSpl, ".")
	if len(strDate[2]) < 4 {
		t.Date = "20" + strDate[2]
	} else {
		t.Date = strDate[2]
	}
	t.Date += "/" + strDate[1] + "/" + strDate[0]
	t.status = true
	os.Remove(path)
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

func (tr *Transaction) getKeyboard() {
	var flag string
	var prompt string
	fmt.Print("Destination " + tr.Destination + "? y/n ")
	fmt.Scan(&flag)
	if flag == "n" {
		tr.Destination = findFild("payee", "Получатель", tr.database)
	}
	if !isExist(tr.Destination, "payee", tr.database) {
		appendFild(tr.Destination, "payee", tr.database)
	}
	for i, el := range  tr.Records {
		prompt = "Счет для " + el.Name
		tr.Records[i].account = findFild("account", prompt, tr.database)
		fmt.Printf("%s\t%.2f\n", tr.Records[i].account,tr.Records[i].Sum)
	}
	prompt ="Счет для оплаты "
	str := findFild("account", prompt, tr.database) 
	fmt.Println(str)
	if strings.Contains(str, "Наличные") {
		fmt.Printf("Сумма %.2f  - ", tr.TotalSum)
		fmt.Scan(&tr.TotalSum)
		if chS := tr.CheckSum(); chS != 0 {
			b := Record{account: "Баланс:Корректировка", Sum: chS}
			tr.Records = append(tr.Records, b)
		}
	}
	r := Record{account: str, Sum: tr.TotalSum * -1}
	tr.Records = append(tr.Records, r)
	fmt.Printf("%s\t%.2f\n",r.account, r.Sum)

}

func appendFild(fild string, pat string, db DB) {
	file, _ := os.OpenFile(db.Path + pat + ".dat", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	defer file.Close()
	file.WriteString(pat + " " + fild + "\n")
}

func findFild(pat string, prompt string, db DB) string {
	exec.Command("setxkbmap", "-layout", "ru").Run()
	defer exec.Command("setxkbmap","-layout","us,ru").Run()
	cmd := "cat " + db.Path + pat +".dat | dmenu -i -p '" + prompt + "'"
	out, _ := exec.Command("bash","-c",cmd).Output()
	return strings.Replace(strings.TrimSpace(string(out)), pat + " ", "", 1)
}

func isExist(fild string, pat string, db DB) bool {
	file, err := os.Open(db.Path + pat + ".dat")
	if err != nil {
		panic("Ошибка чтения")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == pat + " " + fild {
			return true
		}
	}
	return false
}

func (tr *Transaction) SaveTransaction() {
	path, exist := os.LookupEnv("LEDGER")
	if !exist {
		panic("нет пути к журналу")
	}
	file, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	defer file.Close()
	for _, el := range tr.ToStrings() {
		file.WriteString(el+"\n")
	}
}
