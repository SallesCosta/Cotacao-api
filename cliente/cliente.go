package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
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

func main() {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*3000)
	defer cancel()
	getCotacao(ctx)

}

func getCotacao(ctx context.Context) {
	select {
	case <-ctx.Done():
		fmt.Println("deu Timeout.")
		return
	case <-time.After(1 * time.Second):
		fmt.Println("ok")
		Request(ctx)
		return
	}
}

func Request(ctx context.Context) {
	e := http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}

	response, err := e.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	name := "cotacao.txt"
	file, err := os.Create(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar arquivo: %v\n", err)
	}

	defer file.Close()
	file.WriteString(fmt.Sprintf("Dólar: %s", string(body)))
	fmt.Printf("Arquivo criado com sucesso")
}
