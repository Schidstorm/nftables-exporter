package main

import (
	"context"
	"github.com/schidstorm/nftables-exporter/pkg/app"
	"github.com/schidstorm/nftables-exporter/pkg/cli"
	"github.com/schidstorm/nftables-exporter/pkg/config"
)

func main() {

	cli.Run(func(cfg config.Config, applicationContext context.Context) chan error {
		errorChannel := make(chan error, 1)

		go func() {
			errorChannel <- app.RunMetrics(cfg.MetricsPath, cfg.Address)
		}()

		go func() {
			errorChannel <- app.RunRecorder(applicationContext, cfg.Group, cfg.Hostname)
		}()

		return errorChannel
	})

}
