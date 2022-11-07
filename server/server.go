package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type USDBRL struct {
	USDBRL Currency
}

type Currency struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type Cotacao struct {
	Value string `json:"value"`
}

func main() {
	http.HandleFunc("/cotacao", GetExchangeHandler)
	http.ListenAndServe(":8080", nil)
}

func GetExchangeHandler(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Request error: %v\n", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Response body read error: %v\n", err)
	}

	var data USDBRL
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unmarshal error: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	cotacao := Cotacao{Value: data.USDBRL.Bid}

	json.NewEncoder(w).Encode(cotacao)
}
