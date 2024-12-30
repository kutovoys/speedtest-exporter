package web

import (
	"net/http"
	"speedtest-exporter/config"
	"time"
)

func IndexHandler(version, commit string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		if config.CLIConfig.ProtectedMetrics {
			user, pass, ok := r.BasicAuth()
			if !ok || user != config.CLIConfig.MetricsUsername || pass != config.CLIConfig.MetricsPassword {
				w.Header().Set("WWW-Authenticate", `Basic realm="metrics"`)
				http.Error(w, "Unauthorized.", http.StatusUnauthorized)
				return
			}
		}

		data := PageData{
			Version:        version,
			Commit:         commit,
			Port:           config.CLIConfig.Port,
			UpdateInterval: config.CLIConfig.UpdateInterval,
			ServerIDs:      config.CLIConfig.ServerIDs,
		}

		if err := RenderIndex(w, data); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := http.Client{
			Timeout: 3 * time.Second,
		}
		_, err := client.Get("https://clients3.google.com/generate_204")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("No Internet Connection"))
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		}
	}
}
