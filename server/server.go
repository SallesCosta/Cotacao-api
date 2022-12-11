package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Usdbrl struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"var_bid"`
		PctChange  string `json:"pct_change"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"usdbrl"`
}

type ToDb struct {
	ID    int `gorm:"primaryKey"`
	Valor string
}

func main() {
	http.HandleFunc("/cotacao", CotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

func CotacaoHandler(w http.ResponseWriter, r *http.Request) {
	dns := "root:root@tcp(localhost:3306)/escapi?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&ToDb{})

	ctx200 := r.Context()
	ctx200, ctx200Cancel := context.WithTimeout(ctx200, time.Millisecond*200)
	defer ctx200Cancel()

	ctx10 := context.Background()
	ctx10, ctx10Cancel := context.WithTimeout(ctx10, time.Millisecond*10)
	defer ctx10Cancel()

	req, err := http.NewRequestWithContext(ctx200, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var c Usdbrl
	err = json.Unmarshal(body, &c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer o aprse da resposta: %v\n", err)
	}

	result := db.WithContext(ctx10).Create(&ToDb{
		Valor: c.Usdbrl.Bid,
	})

	if result.Error != nil {
		panic(result.Error)
	}

	w.Write([]byte(c.Usdbrl.Bid))
}
