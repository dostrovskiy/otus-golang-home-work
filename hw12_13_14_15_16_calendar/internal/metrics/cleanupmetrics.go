package appmetrics //nolint

import "github.com/prometheus/client_golang/prometheus"

var (
	EventsCleanupStatus = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "event_cleanup_process_status",
            Help: "Status of the event cleanup process (1 = running, 0 = stopped).",
        },
    )

	EventsDeletedCounter = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "event_cleanup_events_deleted_total",
            Help: "Total number of events deleted by the cleanup process.",
        },
    )

	EventsCleanupDuration = prometheus.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "event_cleanup_duration_seconds",
            Help:    "Duration of the event cleanup process.",
            Buckets: prometheus.DefBuckets,
        },
    )
)

func init() {
	prometheus.MustRegister(EventsCleanupStatus)
	prometheus.MustRegister(EventsDeletedCounter)
	prometheus.MustRegister(EventsCleanupDuration)
	// initial values
	EventsDeletedCounter.Add(0)
	EventsCleanupDuration.Observe(0)
}
//nolint