package app

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func RunMetrics(metricsPath, address string) error {
	http.Handle(metricsPath, promhttp.Handler())
	return http.ListenAndServe(address, nil)
}
