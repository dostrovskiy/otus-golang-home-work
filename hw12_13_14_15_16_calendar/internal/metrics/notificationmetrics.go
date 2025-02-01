package appmetrics //nolint

import "github.com/prometheus/client_golang/prometheus"

var (
	EventsNotificationStatus = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "event_notification_process_status",
            Help: "Status of the event notification process (1 = running, 0 = stopped).",
        },
    )

	EventsNotificationCounter = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "event_notification_sent_total",
            Help: "Total number of notifications sent.",
        },
    )
)

func init() {
	prometheus.MustRegister(EventsNotificationStatus)
	prometheus.MustRegister(EventsNotificationCounter)
    // initial values
    EventsNotificationCounter.Add(0)
}
//nolint