package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

type StockData struct {
	TimeSeries map[string]DailyData `json:"Time Series (Daily)"`
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
	Symbol          string  `json:"2. Symbol"`
	NDays           int     `json:"3. NDays"`
	AvgClosingPrice float64 `json:"4. Avg Closing Price"`
}

func main() {
	apiKey := os.Getenv("API_KEY") // Retrieve API key from environment variable

	mux := http.NewServeMux()
	mux.HandleFunc("/stockticker/{symbol}/lastndays/{ndays}", func(w http.ResponseWriter, r *http.Request) {
		handleStockRequest(w, r, apiKey)
	})
	log.Print("starting server on :8080")

	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
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
	currentTime := time.Now()

	for i := 0; i < ndays; i++ {
		date := currentTime.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")

		dailyDataInterface, ok := timeSeries[dateStr].(map[string]interface{})
		fmt.Println(dailyDataInterface)
		if !ok {
			continue // Skip to the next day
		}

		dailyData := DailyData{}
		dailyDataBytes, _ := json.Marshal(dailyDataInterface)
		json.Unmarshal(dailyDataBytes, &dailyData)
		stockResponse.TimeSeriesDaily[dateStr] = dailyData

		closePrice, err := strconv.ParseFloat(dailyData.Close, 64)
		if err != nil {
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

	avgPrice := 0.0
	if len(closingPrices) > 0 {
		avgPrice = sum / float64(len(closingPrices))
	}

	stockResponse.MetaData.AvgClosingPrice = avgPrice

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(stockResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
