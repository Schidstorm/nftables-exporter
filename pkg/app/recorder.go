package app

import (
	"context"
	nftables_exporter "github.com/schidstorm/nftables-exporter/pkg/nftables-exporter"
	"github.com/schidstorm/nftables-exporter/pkg/recorder"
)

func RunRecorder(ctx context.Context, group int, hostname string) error {
	return recorder.Packet(ctx, group, hostname, nftables_exporter.NewQueue())
}
