package recorder

type PayloadInterface interface {
	GetInDev() uint32
	GetOutDev() uint32
	GetData() []byte
}
