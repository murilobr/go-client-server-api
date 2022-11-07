package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
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
	ctxReq := context.Background()
	ctxReq, cancelCtxReq := context.WithTimeout(ctxReq, time.Millisecond*200)
	defer cancelCtxReq()

	req, err := http.NewRequestWithContext(ctxReq, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Create new request error: %v\n", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Request error: %v\n", err)
		w.WriteHeader(http.StatusRequestTimeout)
		w.Write([]byte(""))
		return
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

	// ----------------------------------------------------------------
	// SAVE TO DATABASE

	ctxDb := context.Background()
	ctxDb, cancelCtxDb := context.WithTimeout(ctxDb, time.Millisecond*10)
	defer cancelCtxDb()

	db, err := sql.Open("sqlite3", "cotacao.db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot connect to database error: %v\n", err)
	}
	defer db.Close()

	createCotacaoTable := `CREATE TABLE cotacao (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, value TEXT);`

	stmt, err := db.Prepare(createCotacaoTable)
	if err == nil {
		stmt.Exec()
	}

	insertCotacaoSQL := `INSERT INTO cotacao(value) VALUES (?)`
	stmt, err = db.Prepare(insertCotacaoSQL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create statement to insert into table cotacao error: %v\n", err)
	}

	_, err = stmt.ExecContext(ctxDb, cotacao.Value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot insert into table cotacao error: %v\n", err)
	}
}
