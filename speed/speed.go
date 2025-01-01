package speed

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"speedtest-exporter/config"
	"speedtest-exporter/metrics"
	"speedtest-exporter/models"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/showwin/speedtest-go/speedtest"
)

func RunTests(serverIDs []string) error {
	defer func() {
		runtime.GC()
		debug.FreeOSMemory()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	for _, id := range serverIDs {
		id = strings.TrimSpace(id)
		serverCtx, serverCancel := context.WithTimeout(ctx, 2*time.Minute)

		if err := runAllTests(serverCtx, id); err != nil {
			log.Printf("Tests failed for server ID %s: %v", id, err)
		}

		serverCancel()
		runtime.GC()
		debug.FreeOSMemory()
		time.Sleep(5 * time.Second)
	}

	return nil
}

func runAllTests(ctx context.Context, serverID string) error {
	client := speedtest.New()
	defer func() {
		if manager := client.Manager; manager != nil {
			manager.Reset()
			if snapshots := manager.Snapshots(); snapshots != nil {
				snapshots.Clean()
			}
			manager.Wait()
		}
	}()

	var server *speedtest.Server
	var err error

	if serverID == "" || serverID == "0" {
		serverList, err := client.FetchServerListContext(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch server list: %v", err)
		}
		if len(serverList) == 0 {
			return fmt.Errorf("no servers available")
		}
		server = serverList[0]
	} else {
		server, err = client.FetchServerByIDContext(ctx, serverID)
		if err != nil {
			return fmt.Errorf("failed to fetch server: %v", err)
		}
	}

	server.Context = client

	log.Printf("[%s] %.2fkm %s by %s", server.ID, server.Distance, server.Name, server.Sponsor)

	if err := server.PingTestContext(ctx, nil); err != nil {
		return fmt.Errorf("ping test failed: %v", err)
	}

	if err := server.DownloadTestContext(ctx); err != nil {
		return fmt.Errorf("download test failed: %v", err)
	}

	if err := server.UploadTestContext(ctx); err != nil {
		return fmt.Errorf("upload test failed: %v", err)
	}

	if !server.CheckResultValid() {
		return fmt.Errorf("invalid test results for server %s", server.Name)
	}

	results := models.TestResults{
		Latency:       float64(server.Latency.Milliseconds()),
		Jitter:        float64(server.Jitter.Microseconds()) / 1000.0,
		DownloadSpeed: server.DLSpeed.Mbps() * 1_000_000,
		UploadSpeed:   server.ULSpeed.Mbps() * 1_000_000,
		TestDuration:  server.TestDuration.Total.Seconds(),
	}

	labels := getLabels(server, config.CLIConfig.Instance)
	metrics.Latency.With(labels).Set(results.Latency)
	if server.Jitter > 0 {
		metrics.Jitter.With(labels).Set(results.Jitter)
	}
	metrics.DownloadSpeed.With(labels).Set(results.DownloadSpeed)
	metrics.UploadSpeed.With(labels).Set(results.UploadSpeed)
	metrics.TestDuration.With(labels).Set(results.TestDuration)

	log.Printf("[%s] Download: %.2fMbps, Upload: %.2fMbps, Latency: %.2fms, Jitter: %.2fms, TestDuration: %.2fs",
		server.ID,
		results.DownloadSpeed/1_000_000,
		results.UploadSpeed/1_000_000,
		results.Latency,
		results.Jitter,
		results.TestDuration,
	)

	return nil
}

func getLabels(server *speedtest.Server, instance string) prometheus.Labels {
	labels := prometheus.Labels{
		"server_id":       server.ID,
		"server_name":     server.Name,
		"server_sponsor":  server.Sponsor,
		"server_distance": fmt.Sprintf("%.2fkm", server.Distance),
	}
	if instance != "" {
		labels["instance"] = instance
	}
	return labels
}
