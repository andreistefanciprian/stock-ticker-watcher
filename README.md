# Stock Ticker Watcher

This is a simple Go based stock ticker watcher application that retrieves stock data from the AlphaVantage API.

## Run

```bash
export API_KEY=<YOUR_API_KEY>
# Replace `<YOUR_API_KEY>` with your actual [AlphaVantage](https://www.alphavantage.co/support/#api-key) API key.

# Run as binary
go build -o stock-ticker-watcher
./stock-ticker-watcher --apikey $API_KEY

# Run as container
docker build -t stock-ticker-watcher . -f infra/Dockerfile
docker run -p 8080:8080 stock-ticker-watcher --apikey $API_KEY
```

## Usage

To check the stock price for a particular symbol for the last N days, access the following URL in your browser:

```
http://localhost:8080/stockticker/<SYMBOL>/lastndays/<NDAYS>
```

*   Replace `<SYMBOL>` with the stock symbol (e.g., UBER, AAPL, MSFT).
*   Replace `<NDAYS>` with the number of days for which you want the data (e.g., 20, 50, 100).

**Example:**

To check the stock price for UBER for the last 20 days:

```
http://localhost:8080/stockticker/UBER/lastndays/20
```