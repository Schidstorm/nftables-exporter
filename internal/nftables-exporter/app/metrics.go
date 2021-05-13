package app

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func RunMetrics(metricsPath, address string) chan error {
	httpError := make(chan error)
	go func() {
		http.Handle(metricsPath, promhttp.Handler())
		httpError <- http.ListenAndServe(address, nil)
	}()

	return httpError
}