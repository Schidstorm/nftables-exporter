package recorder

import "github.com/chifflier/nflog-go/nflog"

type Payload struct {
	nflogPayload *nflog.Payload
	inDev uint32
	outDev uint32
}

func (p *Payload) GetInDev() uint32 {
	return p.inDev
}

func (p *Payload) GetOutDev() uint32 {
	return p.outDev
}

func (p *Payload) GetData() []byte {
	return p.nflogPayload.Data
}
