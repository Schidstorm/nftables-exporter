package nftables_exporter

import "github.com/chifflier/nflog-go/nflog"

type QueueInterface interface {
	Initialize(group int, handler nflog.Callback) error
	Start() error
	Stop()
}
