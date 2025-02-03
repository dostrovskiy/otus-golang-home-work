package appmetrics //nolint

import "github.com/prometheus/client_golang/prometheus"

var (
	EventsRequestsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "event_service_requests_total",
			Help: "Total number of requests.",
		},
		[]string{"method", "endpoint", "status_group"},
	)

	EventsRequestsDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "event_service_request_duration_seconds",
			Help:    "Duration of requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	EventsConcurrentRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "event_service_concurrent_requests",
			Help: "Number of concurrent requests.",
		},
		[]string{"method", "endpoint"},
	)
)

func init() {
	prometheus.MustRegister(EventsRequestsCounter)
	prometheus.MustRegister(EventsRequestsDuration)
	prometheus.MustRegister(EventsConcurrentRequests)
	// initial values
	EventsRequestsCounter.WithLabelValues("GET", "/event:id", "2xx").Add(0)
	EventsRequestsDuration.WithLabelValues("GET", "/event:id").Observe(0)
	EventsConcurrentRequests.WithLabelValues("GET", "/event:id").Set(0)
}
//nolint