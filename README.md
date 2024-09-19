# Client-Server API for USD-BRL Exchange Rate

This project consists of two Go programs, `client.go` and `server.go`, which implement a client-server architecture to fetch and store the USD to BRL exchange rate. The exchange rate is fetched from [AwesomeAPI](https://economia.awesomeapi.com.br/json/last/USD-BRL), and the server stores it in an SQLite database. The client receives the exchange rate from the server and saves it in a local file.

## Desafio

Neste desafio vamos aplicar o que aprendemos sobre webserver http, contextos,
banco de dados e manipulação de arquivos com Go.
 
Você precisará nos entregar dois sistemas em Go:
- client.go
- server.go
 
Os requisitos para cumprir este desafio são:
 
O client.go deverá realizar uma requisição HTTP no server.go solicitando a cotação do dólar.
 
O server.go deverá consumir a API contendo o câmbio de Dólar e Real no endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL e em seguida deverá retornar no formato JSON o resultado para o cliente.
 
Usando o package "context", o server.go deverá registrar no banco de dados SQLite cada cotação recebida, sendo que o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms e o timeout máximo para conseguir persistir os dados no banco deverá ser de 10ms.
 
O client.go precisará receber do server.go apenas o valor atual do câmbio (campo "bid" do JSON). Utilizando o package "context", o client.go terá um timeout máximo de 300ms para receber o resultado do server.go.
 
Os 3 contextos deverão retornar erro nos logs caso o tempo de execução seja insuficiente.
 
O client.go terá que salvar a cotação atual em um arquivo "cotacao.txt" no formato: Dólar: {valor}
 
O endpoint necessário gerado pelo server.go para este desafio será: /cotacao e a porta a ser utilizada pelo servidor HTTP será a 8080.

## Project Structure

```bash
client-server-api/
├── client/
│   └── client.go     # Client to request the exchange rate from the server
├── server/
│   └── server.go     # Server to fetch the exchange rate and serve it to the client
├── go.mod            # Go module definition
└── README.md         # Project documentation
```

## Features

- Client: Requests the current USD-BRL exchange rate from the server and saves it in a file (cotacao.txt).
- Server: Fetches the exchange rate from the AwesomeAPI, stores it in an SQLite database, and serves the exchange rate to the client via an HTTP API.
- Timeouts: Both the client and the server implement timeouts using Go's context package.
    The client has a 300ms timeout to receive the exchange rate from the server.
    The server has a 200ms timeout to fetch the exchange rate from the external API and a 10ms timeout to persist data in the database.
- SQLite Database: The server logs each exchange rate fetched in an exchange_rates table.

## Endpoints

- /cotacao: The server exposes this endpoint on port 8080 to return the current USD-BRL exchange rate in JSON format.

## Setup Instructions
1. Prerequisites

    Go (version 1.16 or higher)
    SQLite3 installed

2. Clone the repository

```bash
git clone git@github.com:JRRGomes/client-server-api.git
cd client-server-api
```
3. Initialize the Go Module

```bash
go mod tidy
```
4. Running the Server

Navigate to the server directory and run the server:

```bash
cd server
go run server.go
```
The server will start on http://localhost:8080.
5. Running the Client

Open a new terminal window, navigate to the client directory, and run the client:

```bash
cd client
go run client.go
```
The client will request the exchange rate from the server and save it in a file called cotacao.txt.
6. Verifying Data in the SQLite Database

To check the exchange rates saved by the server, you can use the SQLite3 command-line tool:

```bash
sqlite3 exchange.db
sqlite> SELECT * FROM exchange_rates;
```
Example of Usage

    Start the server:

```bash
go run server/server.go
```

Request the exchange rate using the client:

```bash
go run client/client.go
```

Check the saved exchange rate in cotacao.txt:

```bash
cat cotacao.txt
```
Query the SQLite database for historical exchange rates:

```bash
sqlite3 server/exchange.db
```

## Error Handling

The client and server both log timeout errors if the API call or database operation takes longer than the defined time limits. You can view the logs in the terminal where each service is running.
License

## This project is licensed under the MIT License. See the LICENSE file for details.
Author

JRRGomes - [GitHub](https://github.com/JRRGomes/)
