package core

type QRCode struct {
	TenantID    string
	OrderID     string
	OrderItemID string
	Image       []byte
	FileType    string
}

type QRCodeStore interface {
	Save(qr *QRCode) (string, error)
}

type QRCodeGenerator interface {
	Encode(data string, size int) (*QRCode, error)
}
