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
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/alexsasharegan/dotenv"
	"github.com/prometheus/client_golang/prometheus"
)

// Global vars for metrics
var (
	streamViewers = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "lmhd",
			Subsystem: "twitch",
			Name:      "stream_viewers",
			Help:      "Number of viewers of a stream",
		},
		[]string{
			// Which twitch channel?
			"channel",
		},
	)
	channelFollowers = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "lmhd",
			Subsystem: "twitch",
			Name:      "channel_followers",
			Help:      "Number of followers of a channel",
		},
		[]string{
			// Which twitch channel?
			"channel",
		},
	)
	channelViews = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "lmhd",
			Subsystem: "twitch",
			Name:      "channel_views",
			Help:      "Number of views of a channel",
		},
		[]string{
			// Which twitch channel?
			"channel",
		},
	)
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
	// SimpleGaugeVec Example
	//

	prometheus.MustRegister(streamViewers)
	prometheus.MustRegister(channelFollowers)
	prometheus.MustRegister(channelViews)

	//
	// Loops
	//

	go metricsHandler()

	metricsUpdate()
}

func metricsHandler() {

	flag.Parse()
	http.Handle("/metrics", prometheus.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func metricsUpdate() {

	// Loop to update gauges
	// Ideally this would just be a GaugeVecFunc, but that doesn't exist
	for {
		log.Debugf("Querying...")

		// TODO: start of loop, take a set of channel names, get ids

		// Channel names to IDs
		channels := map[string]string{
			"seraphimkimiko":  "182561942",
			"slytq":           "89854598",
			"nintendo":        "37319",
			"twitch":          "12826",
			"yogscast":        "20786541",
			"el_funko":        "28160719",
			"nerdcubed":       "29660771",
			"loadingreadyrun": "27132299",
			"kate":            "73625408",
			"bengineering":    "113481237",
		}

		for name, id := range channels {
			streamData, err := KrakenStreamsRequest(id)

			if err != nil {
				log.Errorf("Error: %s", err)

				// TODO: if not live, do something else
			} else {
				streamViewers.With(prometheus.Labels{"channel": name}).Set(float64(streamData.Stream.Viewers))
				channelFollowers.With(prometheus.Labels{"channel": name}).Set(float64(streamData.Stream.Channel.Followers))
				channelViews.With(prometheus.Labels{"channel": name}).Set(float64(streamData.Stream.Channel.Views))
				// TODO: more metrics!
			}
		}

		log.Debugf("Sleeping")

		// TODO: Calculate rate limit (30 reqs per second)
		time.Sleep(30 * time.Second)
	}
}
