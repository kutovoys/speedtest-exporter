package speed

import (
	"fmt"
	"log"
	"speedtest-exporter/config"
	"speedtest-exporter/metrics"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/showwin/speedtest-go/speedtest"
)

func RunTests(serverIDs []string) error {
	var servers []*speedtest.Server

	if len(serverIDs) == 0 || serverIDs[0] == "0" {
		serverList, err := speedtest.FetchServers()
		if err != nil {
			return fmt.Errorf("failed to fetch server list: %v", err)
		}
		servers = append(servers, serverList[0])
	} else {
		for _, id := range serverIDs {
			id = strings.TrimSpace(id)
			if id == "" {
				continue
			}
			server, err := speedtest.FetchServerByID(id)
			if err != nil {
				log.Printf("Failed to fetch server with ID %s: %v", id, err)
				continue
			}
			servers = append(servers, server)
		}
	}

	var wg sync.WaitGroup
	for _, server := range servers {
		wg.Add(1)
		go func(srv *speedtest.Server) {
			defer wg.Done()
			if err := RunPingTest(srv); err != nil {
				log.Printf("Ping test failed for server %s: %v", srv.Name, err)
			}
		}(server)
	}
	wg.Wait()

	for _, server := range servers {
		if err := RunSpeedTest(server); err != nil {
			log.Printf("Speed test failed for server %s: %v", server.Name, err)
		}
		time.Sleep(5 * time.Second)
	}

	return nil
}

func RunPingTest(server *speedtest.Server) error {
	log.Printf("Starting ping test for server %s", server.Name)
	server.PingTest(nil)

	labels := getLabels(server, config.CLIConfig.Instance)

	latencyMs := float64(server.Latency.Milliseconds())
	metrics.Latency.With(labels).Set(latencyMs)
	log.Printf("Latency measured for server %s: %.2f ms", server.Name, latencyMs)

	if server.Jitter > 0 {
		jitterMs := float64(server.Jitter.Microseconds()) / 1000.0
		metrics.Jitter.With(labels).Set(jitterMs)
		log.Printf("Jitter measured for server %s: %.2f ms", server.Name, jitterMs)
	} else {
		log.Printf("No jitter data available for server %s", server.Name)
	}

	return nil
}

func RunSpeedTest(server *speedtest.Server) error {
	log.Printf("Starting speed tests for server %s", server.Name)

	log.Printf("Starting download test for server %s", server.Name)
	err := server.DownloadTest()
	if err != nil {
		return fmt.Errorf("download test failed: %v", err)
	}
	log.Printf("Download test completed for server %s. Speed: %.2f Mbps", server.Name, server.DLSpeed.Mbps())

	log.Printf("Starting upload test for server %s", server.Name)
	err = server.UploadTest()
	if err != nil {
		return fmt.Errorf("upload test failed: %v", err)
	}
	log.Printf("Upload test completed for server %s. Speed: %.2f Mbps", server.Name, server.ULSpeed.Mbps())

	labels := getLabels(server, config.CLIConfig.Instance)

	metrics.DownloadSpeed.With(labels).Set(server.DLSpeed.Mbps() * 1_000_000)
	metrics.UploadSpeed.With(labels).Set(server.ULSpeed.Mbps() * 1_000_000)

	return nil
}

func getLabels(server *speedtest.Server, instance string) prometheus.Labels {
	labels := prometheus.Labels{
		"server_id":      server.ID,
		"server_name":    server.Name,
		"server_sponsor": server.Sponsor,
	}
	if instance != "" {
		labels["instance"] = instance
	}
	return labels
}
