package monitor

import "github.com/prometheus/client_golang/prometheus"

var (
	// PriceMonitorErrorCounterMetricName is the name of the Prometheus metric for measuring the number of pricing errors.
	PriceMonitorErrorCounterMetricName = "price_monitor_price_errors_total"

	// PriceMonitorErrorCounterMetricName is the name of the Prometheus metric for sending a heartbeat signal of the price monitor.
	PriceMonitorHeartbeatMetricName = "price_monitor_heartbeat"

	// PricingErrorCounter is a Prometheus counter that measures the number of pricing errors.
	// This metric can be used to monitor errors in the price monitor.
	PricingErrorCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: PriceMonitorErrorCounterMetricName,
			Help: "Total number of pricing errors",
		},
	)

	// PricingHeartbeatCounter is a Prometheus counter that sends a heartbeat signal of the price monitor.
	// This metric can be used to monitor the health of the price monitor.
	PricingHeartbeatCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: PriceMonitorHeartbeatMetricName,
			Help: "Total number of pricing measurements",
		},
	)
)

// init registers metrics with Prometheus
func init() {
	prometheus.MustRegister(PricingErrorCounter)
	prometheus.MustRegister(PricingHeartbeatCounter)
}
