// Структура данных myLedger
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Record Еллемент транзакции
type Record struct {
	account string
	comment string
	sum     float64
}

// Transaction Элемент базы данных
type Transaction struct {
	Date        string
	status      bool
	Destination string
	Records     []Record
}

func (p Product) ProductToRecord() *Record {
	var r Record
	r.sum = p.Sum
	r.comment = fmt.Sprintf("|%.2f * %.3f|%v", p.Price, p.Quantity, p.Name)
	return &r
}

func (r *Record) format() string {
	l := 41 - len(fmt.Sprintf("%.2f", r.sum))
	str := "    %-" + fmt.Sprint(l) + "s$%.2f"
	if len(r.comment) != 0 {
		str += "  ;  " + r.comment
	}
	return fmt.Sprintf(str, r.account, r.sum)
}

func (t *Transaction) Sum() float64 {
	var sum float64
	for _, el := range t.Records {
		sum += el.sum
	}
	return sum
}

func (t *Transaction) toStrings() []string {
	var res []string
	var status string
	if t.status {
		status = "* "
	} else {
		status = ""
	}
	res = append(res, t.Date+" "+status+t.Destination)
	for i := 0; i < len(t.Records); i++ {
		res = append(res, t.Records[i].format())
	}
	return res
}

func (ch Check) CheckToTransaction(db dataBase) *Transaction {

	if !db.Connect {
		panic("соединение с базой не установлено")
	}
	var t Transaction
	t.database = db
	var flag string
	var prompt string
	fmt.Print("Destination " + t.Destination + "? y/n ")
	fmt.Scan(&flag)
	if flag == "n" {
		t.Destination = findFild("payee", "Получатель", t.database)
	}
	if !isExist(t.Destination, "payee", t.database) {
		appendFild(t.Destination, "payee", t.database)
	}
	for _, el := range ch.Products {
		for b := true; b; {
			prompt = "Счет для " + el.Name
			acc := findFild("account", prompt, t.database)
			b = false
			if !isExist(acc, "account", t.database) {
				fmt.Println("Добавить " + acc + "?y/n")
				fmt.Scan(&flag)
				if b = (flag == "n"); !b {
					appendFild(acc, "account", t.database)
				}
			}
		}
		t.Records = append(t.Records, *el.ProductToRecord())
	}
	prompt = "Счет для оплаты "
	str := findFild("account", prompt, t.database)
	fmt.Println(str)
	if strings.Contains(str, "Наличные") {
		fmt.Printf("Сумма %.2f  - ", ch.TotalSum)
		fmt.Scan(&ch.TotalSum)
		if chS := ch.TotalSum - t.Sum(); chS != 0 {
			b := Record{account: "Баланс:Корректировка", sum: chS}
			t.Records = append(t.Records, b)
		}
	}
	b := Record{account: str, sum: ch.TotalSum * -1}
	t.Records = append(t.Records, b)
	strSpl := strings.Split(t.Date, " ")[0]
	strDate := strings.Split(strSpl, ".")
	if len(strDate[2]) < 4 {
		t.Date = "20" + strDate[2]
	} else {
		t.Date = strDate[2]
	}
	t.Date += "/" + strDate[1] + "/" + strDate[0]
	t.status = true
	return &t
}

func dataBaseConnect() dataBase {
	var db dataBase
	path, exist := os.LookupEnv("LEDGERPATH")
	if !exist {
		panic("нет пути к базе")
	}
	err := exec.Command("bash", "-c", "git -C /home/serg/money push").Run()
	if err != nil {
		panic("Проблеммы синхронизации git")
	}
	db.Path = path + "/"
	db.Connect = true
	return db
}

func appendFild(fild string, pat string, db dataBase) {
	file, _ := os.OpenFile(db.Path+pat+".dat", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	defer file.Close()
	file.WriteString(pat + " " + fild + "\n")
}

func findFild(pat string, prompt string, db dataBase) string {
	exec.Command("setxkbmap", "-layout", "ru").Run()
	defer exec.Command("setxkbmap", "-layout", "us,ru").Run()
	cmd := "cat " + db.Path + pat + ".dat | dmenu -i -p '" + prompt + "'"
	out, _ := exec.Command("bash", "-c", cmd).Output()
	return strings.Replace(strings.TrimSpace(string(out)), pat+" ", "", 1)
}

func isExist(fild string, pat string, db dataBase) bool {
	file, err := os.Open(db.Path + pat + ".dat")
	if err != nil {
		panic("Ошибка чтения")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == pat+" "+fild {
			return true
		}
	}
	return false
}

func (t *Transaction) saveTransaction() {
	path, exist := os.LookupEnv("LEDGER")
	if !exist {
		panic("нет пути к журналу")
	}
	file, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	defer file.Close()
	for _, el := range t.toStrings() {
		file.WriteString(el + "\n")
	}
}
