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
	"strings"

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

	streamFps = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "lmhd",
			Subsystem: "twitch",
			Name:      "stream_average_fps",
			Help:      "Average FPS of a stream",
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
	LogLevel       string
	KrakenClientID string
	Channels       []string
}

var cfg Config

var kraken *KrakenClient

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

	if len(os.Getenv("KRAKEN_CLIENT_ID")) == 0 {
		log.Fatalf("KRAKEN_CLIENT_ID not set!")
	}
	cfg.KrakenClientID = os.Getenv("KRAKEN_CLIENT_ID")

	if len(os.Getenv("TWITCH_CHANNELS")) == 0 {
		log.Fatalf("TWITCH_CHANNELS not set!")
	}
	cfg.Channels = strings.Split(os.Getenv("TWITCH_CHANNELS"), ",")

	//
	// Twitch API Client
	//

	kraken = NewKrakenClient(cfg.KrakenClientID)

	//
	// Register Prometheus Metrics
	//

	prometheus.MustRegister(streamViewers)
	prometheus.MustRegister(streamFps)

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

	krakenUsersResponse, err := kraken.Users(cfg.Channels)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// map of IDs to Names
	var channelIDs = make(map[string]string)
	// slice of IDs
	var channels []string

	for _, user := range krakenUsersResponse.Users {
		channelIDs[user.Name] = user.ID
		channels = append(channels, user.ID)
	}

	log.Debugf("Channel IDs: %s", channelIDs)

	// Loop to update gauges
	for {
		// Keep track of which channels we've had updates for
		// i.e. those which are live
		channelsSeen := make(map[string]string)
		// and those which are not
		channelsUnseen := make(map[string]string)
		for k, v := range channelIDs {
			channelsUnseen[k] = v
		}

		log.Debugf("Querying %d channels", len(channelIDs))
		krakenStreamsResponse, err := kraken.Streams(channels)
		if err != nil {
			log.Fatalf("%s", err)
		}

		log.Debugf("Updating metrics for %d live channels", len(krakenStreamsResponse.Streams))

		for _, stream := range krakenStreamsResponse.Streams {
			name := stream.Channel.Name

			streamViewers.With(prometheus.Labels{"channel": name}).Set(float64(stream.Viewers))
			streamFps.With(prometheus.Labels{"channel": name}).Set(float64(stream.AverageFps))
			channelFollowers.With(prometheus.Labels{"channel": name}).Set(float64(stream.Channel.Followers))
			channelViews.With(prometheus.Labels{"channel": name}).Set(float64(stream.Channel.Views))

			channelsSeen[name] = channelsUnseen[name]
			delete(channelsUnseen, name)
		}

		// For all channels we haven't seen, delete their live viewers
		for name := range channelsUnseen {
			log.Debugf("Deleting metrics for unseen channel %s", name)
			streamViewers.Delete(prometheus.Labels{"channel": name})
			streamFps.Delete(prometheus.Labels{"channel": name})

			// TODO: query API to get these for non-live channels
			channelFollowers.Delete(prometheus.Labels{"channel": name})
			channelViews.Delete(prometheus.Labels{"channel": name})
		}
	}
}
