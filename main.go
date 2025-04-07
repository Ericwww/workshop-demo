package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

func main() {
	loggerInit()
	podIp := os.Getenv("POD_IP")
	ctx := context.Background()
	// The exporter embeds a default OpenTelemetry Reader and
	// implements prometheus.Collector, allowing it to be used as
	// both a Reader and Collector.
	exporter, err := prometheus.New()
	if err != nil {
		log.Error().Err(err).Msg("failed to new prometheus exporter")
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter("my-app")

	helloCounter, err := meter.Int64Counter("hello", api.WithUnit("time(s)"), api.WithDescription("a counter of hello api"))

	// Start the prometheus HTTP server and pass the exporter Collector to it
	go serveMetrics(podIp)
	log.Debug().Msg("prom http server is working")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		helloCounter.Add(ctx, 1)
		log.Info().Msg("Hello World!")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	})

	addr := fmt.Sprintf("%s:3333", podIp)
	log.Debug().Msgf("tring to start server at %s", addr)
	http.ListenAndServe(addr, r)
}

func serveMetrics(podIp string) {
	fmt.Println("serving metrics at localhost:2223/metrics")
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(fmt.Sprintf("%s:2223", podIp), nil) //nolint:gosec // Ignoring G114: Use of net/http serve function that has no support for setting timeouts.
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}

func loggerInit() {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	multi := zerolog.MultiLevelWriter(consoleWriter)
	logPath := os.Getenv("LOG_PATH")
	if logPath != "" {
		f, err := os.Create(logPath)
		if err != nil {
			log.Fatal().Msg("create log path failed")
		}
		multi = zerolog.MultiLevelWriter(consoleWriter, f)
	}
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()
}
