package monitor

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

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

func NewOtelTracer(ctx context.Context, host string) (*sdktrace.TracerProvider, error) {
	// resource.WithContainer() adds container.id which the agent will leverage to fetch container tags via the tagger.
	res, err := resource.New(ctx, resource.WithContainer(),
		resource.WithAttributes(semconv.ServiceNameKey.String(host)),
		resource.WithFromEnv(),
	)
	if err != nil {
		return nil, err
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}
