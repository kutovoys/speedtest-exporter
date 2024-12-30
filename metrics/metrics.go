package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	labelNames = []string{"server_id", "server_name", "server_sponsor"}
)

func UpdateLabelNames(instance string) {
	if instance != "" {
		labelNames = append(labelNames, "instance")
	}
}

var (
	DownloadSpeed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "speedtest_download_bits_per_second",
			Help: "Download speed in bits per second",
		},
		labelNames,
	)

	UploadSpeed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "speedtest_upload_bits_per_second",
			Help: "Upload speed in bits per second",
		},
		labelNames,
	)

	Latency = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "speedtest_latency",
			Help: "Latency in ms",
		},
		labelNames,
	)

	Jitter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "speedtest_jitter",
			Help: "Jitter in ms",
		},
		labelNames,
	)
)
