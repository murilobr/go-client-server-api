package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Value string `json:"value"`
}

func main() {
	ctxReq := context.Background()
	ctxReq, cancelCtxReq := context.WithTimeout(ctxReq, time.Millisecond*300)
	defer cancelCtxReq()

	req, err := http.NewRequestWithContext(ctxReq, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Create new request error: %v\n", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Request error: %v\n", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusRequestTimeout {
		fmt.Fprintf(os.Stderr, "Request timeout\n")
		return
	}

	body, err := io.ReadAll(res.Body)
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
