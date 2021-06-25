package recorder

type packetMock struct {
	data []byte
}

func (p *packetMock) GetInDev() uint32 {
	return 1
}

func (p *packetMock) GetOutDev() uint32 {
	return 1
}

func (p *packetMock) GetData() []byte {
	return p.data
}
