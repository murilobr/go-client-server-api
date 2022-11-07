package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Cotacao struct {
	Value string `json:"value"`
}

func main() {
	resp, err := http.Get("http://localhost:8080/cotacao")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	var cotacao Cotacao
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		println(err)
	}

	fmt.Println(cotacao.Value)
}
