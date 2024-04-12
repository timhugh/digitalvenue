package core

type QRCode struct {
	Image    []byte
	FileType string
}

type QRCodeStorer interface {
	Save(qr *QRCode) error
}

type QRCodeGenerator interface {
	Encode(data string, size int) (*QRCode, error)
}
