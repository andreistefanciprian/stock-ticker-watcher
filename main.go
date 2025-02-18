package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func main() {
	// parse CLI params
	initFlags()
	flag.Parse()
	httpPort := ":" + strconv.Itoa(port)

	baseUrl := "https://www.alphavantage.co"

	mux := http.NewServeMux()
	mux.HandleFunc("GET /stockticker/{symbol}/lastndays/{ndays}", func(w http.ResponseWriter, r *http.Request) {
		handleStockRequest(w, r, baseUrl, apiKey)
	})
	mux.HandleFunc("GET /healthz", healthCheckHandler) // Health check endpoint
	mux.Handle("GET /metrics", promhttp.Handler())     // Metrics endpoint

	log.Printf("starting server on %s", httpPort)

	err := http.ListenAndServe(httpPort, mux)
	log.Fatal(err)
}
