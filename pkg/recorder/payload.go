package recorder

import "github.com/chifflier/nflog-go/nflog"

type Payload struct {
	nflogPayload *nflog.Payload
}

func (p *Payload) GetInDev() uint32 {
	return p.nflogPayload.GetInDev()
}

func (p *Payload) GetOutDev() uint32 {
	return p.nflogPayload.GetOutDev()
}

func (p *Payload) GetData() []byte {
	return p.nflogPayload.Data
}
