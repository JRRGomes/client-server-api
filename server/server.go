package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ExchangeRate struct {
	USD_BRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func main() {
	db, err := sql.Open("sqlite3", "./exchange.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS exchange_rates (id INTEGER PRIMARY KEY, rate TEXT, timestamp DATETIME)`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/cotacao", handleCotacao)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fetchExchangeRate(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return "", err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var rate ExchangeRate
	if err := json.Unmarshal(body, &rate); err != nil {
		return "", err
	}

	return rate.USD_BRL.Bid, nil
}

func saveToDatabase(ctx context.Context, db *sql.DB, rate string) error {
	query := `INSERT INTO exchange_rates (rate, timestamp) VALUES (?, ?)`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, rate, time.Now())
	return err
}

func handleCotacao(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Contexto para a requisição da API com timeout de 200ms
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	rate, err := fetchExchangeRate(ctx)
	if err != nil {
		http.Error(w, "Error fetching exchange rate", http.StatusInternalServerError)
		log.Println("Error:", err)
		return
	}

	// Contexto para salvar no banco com timeout de 10ms
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer dbCancel()

	db, err := sql.Open("sqlite3", "./exchange.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := saveToDatabase(dbCtx, db, rate); err != nil {
		log.Println("Database error:", err)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"bid": rate})
}
