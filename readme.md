# Bitcoin LTP API

This is a simple Go API application that retrieves the Last Traded Price (LTP) of Bitcoin for the currency pairs BTC/USD, BTC/CHF, and BTC/EUR using the Kraken public API.

## Requirements

- Go 1.20 or higher
- Docker (optional, for containerized application)

## Setup

### Local Setup

1. Clone the repository:
   ```sh
   git clone <repository-url>
   cd <repository-name>
   ```

1. Build the application: ```go build -o main ./main.go```
2. Run the application ```PORT=8080 ./main```


The application will be available at http://localhost:8080/api/v1/ltp.

## Docker Setup
Build the Docker image:

1. `docker build -t bitcoin-ltp-api`

2. Run the Docker container

3. `docker run -p 8080:8080 -e PORT=8080 bitcoin-ltp-api`

The application will be available at http://localhost:8080/api/v1/ltp.

## Tests
To run tests, execute:

`go test ./...`
``` API
GET /api/v1/ltp
Response:

json
{
  "ltp": [
    {
      "pair": "BTCUSD",
      "amount": "52000.12"
    },
    {
      "pair": "BTCCHF",
      "amount": "49000.12"
    },
    {
      "pair": "BTCEUR",
      "amount": "50000.12"
    }
  ]
}
```

### `go.mod`
```go
module bitcoin-ltp-api

go 1.20

require (
	github.com/stretchr/testify v1.7.0
)