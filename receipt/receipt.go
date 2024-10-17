package receipt

import (
	"encoding/json"
	"fmt"
	"os"
)

// Product товар
type Product struct {
	Sum      float64 `json:"sum"`
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
}

// Receipt Чек
type Receipt struct {
	Date     string    `json:"date"`
	ShopName string    `json:"shopName"`
	Products []Product `json:"products"`
	TotalSum float64   `json:"totalSum"`
}

func NewCheck(path string) Receipt {
	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("file not open %v", err))
	}
	defer file.Close()
	var ch Receipt
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&ch); err != nil {
		panic(fmt.Sprintf("ERROR Serializ: %v", err))
	}
	return ch
}
