// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// A minimal example of how to include Prometheus instrumentation.
package main

import (
	"flag"
	"net/http"
	"os"
	"runtime"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/alexsasharegan/dotenv"
	"github.com/prometheus/client_golang/prometheus"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

type Config struct {
	LogLevel string
}

var cfg Config

func main() {
	// Load env vars from .env file, if present
	// Ignore errors caused by file not existing
	_ = dotenv.Load()

	// Setup logging before anything else
	if len(os.Getenv("LOG_LEVEL")) == 0 {
		cfg.LogLevel = "info"
	} else {
		cfg.LogLevel = os.Getenv("LOG_LEVEL")
	}
	switch cfg.LogLevel {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	}

	log.Infof("Hello World")

	//
	// Simple Gauge Example
	//

	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "lmhd",
		Subsystem: "twitch",
		Name:      "gauge",
		Help:      "Testing a simple gauge",
	})
	prometheus.MustRegister(gauge)

	//
	// Simple GaugeFunc Example
	//

	if err := prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: "lmhd",
			Subsystem: "twitch",
			Name:      "gauge_func",
			Help:      "Number of goroutines that currently exist.",
		},
		func() float64 { return float64(runtime.NumGoroutine()) },
	)); err == nil {
		log.Debugf("GaugeFunc 'goroutines_count' registered.")
	}

	//
	// SimpleGaugeVec Example
	//

	opsQueued := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "lmhd",
			Subsystem: "twitch",
			Name:      "gauge_vec",
			Help:      "Number of blob storage operations waiting to be processed, partitioned by user and type.",
		},
		[]string{
			// Which user has requested the operation?
			"user",
			// Of what type is the operation?
			"type",
		},
	)
	prometheus.MustRegister(opsQueued)

	//
	// Loops
	//

	go handler()

	// Loop to update gauges
	// Ideally this would just be a GaugeVecFunc, but that doesn't exist
	for {
		log.Debugf("Querying...")
		gauge.Set(3)
		opsQueued.With(prometheus.Labels{"type": "delete", "user": "alice"}).Inc()
		time.Sleep(10 * time.Second)
	}
}

func handler() {

	flag.Parse()
	http.Handle("/metrics", prometheus.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
