package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "Dólar: %s", response.Bid)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Cotação salva com sucesso.")
}
