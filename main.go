package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"speedtest-exporter/config"
	"speedtest-exporter/metrics"
	"speedtest-exporter/speed"
	"speedtest-exporter/web"

	"github.com/go-co-op/gocron"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	version = "unknown"
	commit  = "unknown"
)

func BasicAuthMiddleware(username, password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.CLIConfig.ProtectedMetrics {
				user, pass, ok := r.BasicAuth()
				if !ok || user != username || pass != password {
					w.Header().Set("WWW-Authenticate", `Basic realm="metrics"`)
					http.Error(w, "Unauthorized.", http.StatusUnauthorized)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	config.Parse(version, commit)
	fmt.Println("Speedtest Exporter", version)

	metrics.UpdateLabelNames(config.CLIConfig.Instance)

	reg := prometheus.NewRegistry()

	for _, metric := range metrics.GetMetrics() {
		reg.MustRegister(metric)
	}

	s := gocron.NewScheduler(time.UTC)

	var serverIDs []string
	if config.CLIConfig.ServerIDs != "" {
		serverIDs = strings.Split(config.CLIConfig.ServerIDs, ",")
	}

	s.Every(config.CLIConfig.UpdateInterval).Minutes().Do(func() {
		log.Printf("Starting speed test iteration...")
		if err := speed.RunTests(serverIDs); err != nil {
			log.Printf("Tests failed: %v", err)
		} else {
			log.Printf("All tests completed successfully")
		}
	})

	s.StartAsync()

	http.Handle("/", web.IndexHandler(version, commit))
	http.Handle("/health", web.HealthHandler())
	http.Handle("/metrics", BasicAuthMiddleware(
		config.CLIConfig.MetricsUsername,
		config.CLIConfig.MetricsPassword,
	)(promhttp.HandlerFor(reg, promhttp.HandlerOpts{})))

	log.Printf("Starting server on :%s", config.CLIConfig.Port)
	if err := http.ListenAndServe(":"+config.CLIConfig.Port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
