package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"Health":"OK"}`))
}

func handleStockRequest(w http.ResponseWriter, r *http.Request, baseUrl, apiKey string) {
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
	url := fmt.Sprintf("%s/query?apikey=%s&function=TIME_SERIES_DAILY&symbol=%s", baseUrl, apiKey, symbol)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	RecordRequest(symbol, ndaysStr) // Record the request

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
