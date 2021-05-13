package app

import (
	"context"
	"github.com/schidstorm/nftables-exporter/internal/nftables-exporter/recorder"
)

func RunRecorder(ctx context.Context, group int, hostname string) chan error {
	recorderError := make(chan error)
	go func() {
		recorderError <- recorder.Packet(ctx, group, hostname)
	}()
	return recorderError
}
