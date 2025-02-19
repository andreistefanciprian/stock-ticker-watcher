# Stock Ticker Watcher

This is a Go-based stock ticker application that retrieves historical stock data from the AlphaVantage API and calculates the average closing price over a specified number of days.  It provides a simplified view of the data, focusing on the last N days and their average closing price, rather than displaying the full API response.

The application queries the AlphaVantage API using a URL similar to this: ```https://www.alphavantage.co/query?apikey=export%20API_KEY=<YOUR_API_KEY>&function=TIME_SERIES_DAILY&symbol=<TICKER_SYMBOL>```

However, instead of returning all the data from the API, the application filters the results to include only the data for the last N days and calculates the average closing price for that period.

**Key Features:**

*   Retrieves historical stock data from AlphaVantage.
*   Filters data to show only the last N days.
*   Calculates and returns the average closing price for the specified period.

**Example Usage:**

To check the stock price for UBER for the last 20 days, use the following URL: ```http://<STOCK_TICKER_WATCHER_ADDRESS>/stockticker/UBER/lastndays/20```

## Run

```bash
export API_KEY=<YOUR_API_KEY>

# Run as container
docker build -t stock-ticker-watcher . -f infra/Dockerfile
docker run -p 8080:8080 stock-ticker-watcher --apikey $API_KEY
# Or
docker-compose --build up

# Run as k8s deployment
kubectl create secret generic api-key-secret --from-literal=apiKey=$API_KEY
helm install stock-ticker-watcher infra/k8s/stock-ticker-watcher
k port-forward svc/stock-ticker-watcher 8080:8080
# Check logs
k logs -l app.kubernetes.io/name=stock-ticker-watcher -f
# Uninstall k8s deployment
helm uninstall stock-ticker-watcher
```

**Note:** Replace `<YOUR_API_KEY>` with your actual [AlphaVantage](https://www.alphavantage.co/support/#api-key) API key.

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

## Prometheus Metrics

This application exposes Prometheus metrics at the `/metrics` endpoint.  These metrics can be scraped by a Prometheus server to monitor the application's performance and behavior:

```
# HELP stock_requests_total Total number of stock requests
# TYPE stock_requests_total counter
stock_requests_total{ndays="10",symbol="ABNB"} 11
stock_requests_total{ndays="1000",symbol="MSFT"} 1
stock_requests_total{ndays="100",symbol="MSFT"} 1
stock_requests_total{ndays="20",symbol="ABNB"} 2
stock_requests_total{ndays="20",symbol="MSFT"} 1
```