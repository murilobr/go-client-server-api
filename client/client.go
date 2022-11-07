package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Cotacao struct {
	Value string `json:"value"`
}

func main() {
	resp, err := http.Get("http://localhost:8080/cotacao")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot request endpoint: %v\n", err)
	}

	if resp.StatusCode == http.StatusRequestTimeout {
		fmt.Fprintf(os.Stderr, "Request timeout\n")
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read response body: %v\n", err)
	}

	var cotacao Cotacao
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot unmarshal response to struct: %v\n", err)
		return
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar o arquivo cotacao.txt: %v\n", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", cotacao.Value))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao escrever no arquivo cotacao.txt: %v\n", err)
	}
}
