package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
)

var (
	apiKey string
	port   int
)

func initFlags() {
	// define and parse cli params
	flag.IntVar(&port, "port", 8080, "Port to listen on")
	flag.StringVar(&apiKey, "apikey", "", "AlphaVantage API key")
}

type DailyData struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type StockResponse struct {
	MetaData        StockMetaData        `json:"Meta Data"`
	TimeSeriesDaily map[string]DailyData `json:"Time Series (Daily)"`
}

type StockMetaData struct {
	Symbol          string    `json:"1. Symbol"`
	NDays           int       `json:"2. NDays"`
	ClosingPrices   []float64 `json:"3. Closing Prices"`
	AvgClosingPrice float64   `json:"4. Average Closing Price for last NDays"`
}

func main() {
	// parse CLI params
	initFlags()
	flag.Parse()
	httpPort := ":" + strconv.Itoa(port)
	mux := http.NewServeMux()
	mux.HandleFunc("/stockticker/{symbol}/lastndays/{ndays}", func(w http.ResponseWriter, r *http.Request) {
		handleStockRequest(w, r, apiKey)
	})
	// Add health check handler
	mux.HandleFunc("/healthz", healthCheckHandler) // New health check endpoint

	log.Printf("starting server on %s", httpPort)

	err := http.ListenAndServe(httpPort, mux)
	log.Fatal(err)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK") // Simple "OK" response
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintln(w, "Method Not Allowed")
}

func handleStockRequest(w http.ResponseWriter, r *http.Request, apiKey string) {
	// get the symbol and ndays from the URL
	symbol := r.PathValue("symbol")
	ndaysStr := r.PathValue("ndays")

	ndays, err := strconv.Atoi(ndaysStr)
	if err != nil {
		http.Error(w, "ndays must be an integer", http.StatusBadRequest)
		return
	}

	if ndays <= 0 {
		http.Error(w, "ndays must be greater than 0", http.StatusBadRequest)
		return
	}

	// call the stock API
	url := fmt.Sprintf("https://www.alphavantage.co/query?apikey=%s&function=TIME_SERIES_DAILY&symbol=%s", apiKey, symbol)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data map[string]interface{} // Use a generic map
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	timeSeries, ok := data["Time Series (Daily)"].(map[string]interface{})
	if !ok {
		http.Error(w, "Error: 'Time Series (Daily)' not found in response", http.StatusInternalServerError)
		return
	}

	stockResponse := StockResponse{
		MetaData: StockMetaData{
			Symbol: symbol,
			NDays:  ndays,
		},
		TimeSeriesDaily: make(map[string]DailyData),
	}

	closingPrices := make([]float64, 0, ndays)
	count := 0

	// 1. Collect and sort the dates
	var dates []string
	for dateStr := range timeSeries {
		dates = append(dates, dateStr)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(dates))) // Sort in descending order (most recent first)

	// 2. Iterate over the sorted dates
	for _, dateStr := range dates {
		if count >= ndays {
			break // Stop when ndays data points are collected
		}

		fmt.Println("Processing date from API:", dateStr) // Debugging

		dailyDataInterface, ok := timeSeries[dateStr].(map[string]interface{})
		if !ok {
			fmt.Println("Data not found for:", dateStr) // Debugging
			continue                                    // Skip if data is not found for date
		}

		dailyData := DailyData{}
		dailyDataBytes, _ := json.Marshal(dailyDataInterface)
		json.Unmarshal(dailyDataBytes, &dailyData)
		stockResponse.TimeSeriesDaily[dateStr] = dailyData

		closePrice, err := strconv.ParseFloat(dailyData.Close, 64)
		if err != nil {
			fmt.Println("Parse Float Error:", err)
			http.Error(w, "Error parsing closing price", http.StatusInternalServerError)
			return
		}
		closingPrices = append(closingPrices, closePrice)
		count++
	}

	sum := 0.0
	for _, price := range closingPrices {
		sum += price
	}

	average := sum / float64(len(closingPrices))
	fmt.Printf("no of closing prices %d", len(closingPrices))
	fmt.Printf("The average closing price is: %.2f\n", average)

	stockResponse.MetaData.AvgClosingPrice = average
	stockResponse.MetaData.ClosingPrices = closingPrices

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(stockResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
