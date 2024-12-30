package config

import (
	"speedtest-exporter/models"

	"github.com/alecthomas/kong"
)

var CLIConfig models.CLI

func Parse(version, commit string) {
	ctx := kong.Parse(&CLIConfig,
		kong.Name("speedtest-exporter"),
		kong.Description("A Prometheus exporter for speedtest metrics."),
		kong.Vars{
			"version": version,
			"commit":  commit,
		},
	)
	_ = ctx
}
