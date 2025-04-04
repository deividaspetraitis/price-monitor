package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/deividaspetraitis/price-monitor"
	"github.com/deividaspetraitis/price-monitor/errors"
	ihttp "github.com/deividaspetraitis/price-monitor/http"
	"github.com/deividaspetraitis/price-monitor/log"
	"github.com/deividaspetraitis/price-monitor/provider"
)

// shutdowntimeout is the duration the service will wait for outstanding requests to complete before shutting down.
var shutdowntimeout = time.Duration(5) * time.Second

// Program flags
var (
	host        string
	httpAddress string
	sqsBaseURL  string
	threshold   float64
	interval    int
	otel        bool
)

func init() {
	flag.StringVar(&host, "host", "price-monitor", "the name of the host")
	flag.StringVar(&httpAddress, "http", ":8080", "HTTP service address")
	flag.StringVar(&sqsBaseURL, "sqs-base-url", "http://localhost:9092", "SQS provider base URL")
	flag.Float64Var(&threshold, "threshold", 0.1, "Price difference threshold for logging")
	flag.IntVar(&interval, "interval", 60, "Interval between price checks in seconds")
	flag.BoolVar(&otel, "otel", false, "Enable OpenTelemetry")
	flag.Parse()
}

// pairs is a list of cryptocurrency pairs to monitor.
var pairs = monitor.Pairs{
	{Base: monitor.OSMO, Quote: monitor.USD},
}

// main program entry point.
func main() {
	ctx := context.Background()
	logger := log.Default()

	if otel {
		tp, err := monitor.NewOtelTracer(ctx, host)
		if err != nil {
			logger.Fatalf("Error creating tracer provider: %v", err)
		}

		defer func() {
			if err := tp.Shutdown(ctx); err != nil {
				log.Fatalf("Error shutting down tracer provider: %v", err)
			}
		}()
	}

	if err := run(ctx, logger); err != nil {
		logger.WithError(err).Error("unable to start service")
		os.Exit(1)
	}
}

func run(ctx context.Context, logger log.Logger) error {
	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// =========================================================================
	// Start HTTP server

	api := http.Server{
		Addr:    httpAddress,
		Handler: ihttp.API(shutdown),
	}

	go func() {
		logger.Printf("http server listening on %s", httpAddress)
		serverErrors <- api.ListenAndServe()
	}()

	// =========================================================================
	// Start Service

	// TODO: Handle shutdown gracefully
	providers := []monitor.Provider{
		provider.NewCoinGeckoClient(),
		provider.NewSQSClient(sqsBaseURL),
	}

	monitor.Compare(monitor.Fetch(providers, pairs, logger), threshold, logger)

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				monitor.Compare(monitor.Fetch(providers, pairs, logger), threshold, logger)
			case <-ctx.Done():
				return
			}
		}
	}()

	// ========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		logger.Printf("http server start shutdown caused by %v", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), shutdowntimeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			logger.WithError(err).Error("graceful shutdown did not complete")
			api.Close()
		}

		// Log the status of this shutdown.
		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}
