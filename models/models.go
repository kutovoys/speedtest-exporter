package models

import (
	"fmt"

	"github.com/alecthomas/kong"
)

type CLI struct {
	Port             string      `name:"metrics-port" help:"Port to listen on" default:"9090" env:"METRICS_PORT"`
	ProtectedMetrics bool        `name:"metrics-protected" help:"Whether metrics are protected by basic auth" default:"false" env:"METRICS_PROTECTED"`
	MetricsUsername  string      `name:"metrics-username" help:"Username for metrics if protected by basic auth" default:"metricsUser" env:"METRICS_USERNAME"`
	MetricsPassword  string      `name:"metrics-password" help:"Password for metrics if protected by basic auth" default:"MetricsVeryHardPassword" env:"METRICS_PASSWORD"`
	UpdateInterval   int         `name:"update-interval" help:"Interval for metrics update in minutes" default:"3" env:"UPDATE_INTERVAL"`
	ServerIDs        string      `name:"server-ids" help:"Comma-separated list of speedtest server IDs" default:"0" env:"SERVER_IDS"`
	Instance         string      `name:"instance" help:"Instance label for metrics" env:"INSTANCE"`
	Version          VersionFlag `name:"version" help:"Print version information and quit"`
}

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println("Speedtest Exporter")
	fmt.Printf("Version:\t %s\n", vars["version"])
	fmt.Printf("Commit:\t %s\n", vars["commit"])
	fmt.Printf("GitHub: https://github.com/kutovoys/speedtest-exporter\n")
	app.Exit(0)
	return nil
}
