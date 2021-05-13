package cli

import (
	"context"
	"github.com/schidstorm/nftables-exporter/internal/nftables-exporter/app"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	rootCommand := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			var group int
			if cliGroup, err := cmd.PersistentFlags().GetInt("group"); err != nil {
				return err
			} else {
				group = cliGroup
			}

			var hostname string
			if cliHostname, err := cmd.PersistentFlags().GetString("hostname"); err != nil {
				return err
			} else {
				hostname = cliHostname
			}

			var address string
			if cliAddress, err := cmd.PersistentFlags().GetString("address"); err != nil {
				return err
			} else {
				address = cliAddress
			}

			var metricsPath string
			if cliMetricsPath, err := cmd.PersistentFlags().GetString("metrics-path"); err != nil {
				return err
			} else {
				metricsPath = cliMetricsPath
			}

			applicationContext, cancel := context.WithCancel(context.Background())
			defer cancel()

			httpError := app.RunMetrics(metricsPath, address)
			recorderError := app.RunRecorder(applicationContext, group, hostname)

			signalChannel := make(chan os.Signal)
			signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

			select {
			case err := <- httpError:
				return err
			case err := <- recorderError:
				return err
			case <- signalChannel:
			}

			return nil
		},
	}

	var hostname = "unknown"
	if osHostname, err := os.Hostname(); err == nil {
		hostname = osHostname
	}

	rootCommand.PersistentFlags().Int("group", 0, "netfilter group number")
	rootCommand.PersistentFlags().String("hostname", hostname, "hostname passed as metric label")
	rootCommand.PersistentFlags().String("address", ":2112", "listen address for metrics http server")
	rootCommand.PersistentFlags().String("metrics-path", "/metrics", "http path for metrics")

	if err := rootCommand.Execute(); err != nil {
		logrus.Error(err)
	}
}
