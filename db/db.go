package db

import (
	"os"
)

type dataBase struct {
	Path    string
	Connect bool
}

func DataBaseConnect() dataBase {
	var db dataBase
	path, exist := os.LookupEnv("LEDGERPATH")
	if !exist {
		panic("нет пути к базе")
	}
	/*err := exec.Command("bash", "-c", "git -C /home/serg/money push").Run()
	if err != nil {
		panic("Проблеммы синхронизации git")
	}*/
	db.Path = path + "/"
	db.Connect = true
	return db
}
