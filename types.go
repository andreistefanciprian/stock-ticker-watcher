package main

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
