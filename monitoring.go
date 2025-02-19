package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "stock_requests_total",
		Help: "Total number of stock requests",
	},
	[]string{"ndays", "symbol"},
)

func init() {
	prometheus.MustRegister(totalRequests)
}

func RecordRequest(symbol string, ndays string) {
	totalRequests.WithLabelValues(ndays, symbol).Inc()
}
