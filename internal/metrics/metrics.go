package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var RequestsByEndpoint = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "requests_processed_by_endpoint",
		Help:    "The total number of processed requests",
		Buckets: []float64{10, 50, 100, 200, 500, 1000},
	},
	[]string{"endpoint"},
)

var StatusCodes = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_status_codes",
		Help: "HTTP response status codes",
	},
	[]string{"status_code"},
)
