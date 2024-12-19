package main

import (
	"flag"
	"log"
	"os"
	"time"

	monitor "github.com/deividaspetraitis/price-monitor"
	"github.com/deividaspetraitis/price-monitor/provider"
)

var (
	threshold float64
	interval  int
)

func init() {
	flag.Float64Var(&threshold, "threshold", 0.1, "Price difference threshold for logging")
	flag.IntVar(&interval, "interval", 60, "Interval between price checks in seconds")
	flag.Parse()
}

// pairs is a list of cryptocurrency pairs to monitor.
var pairs = monitor.Pairs{
	{Base: monitor.OSMO, Quote: monitor.USD},
}

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	providers := []monitor.Provider{
		provider.NewCoinGeckoClient(),
		provider.NewSQSClient(),
	}

	monitor.Compare(monitor.Fetch(providers, pairs, logger), threshold, logger)

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			monitor.Compare(monitor.Fetch(providers, pairs, logger), threshold, logger)
		}
	}
}
