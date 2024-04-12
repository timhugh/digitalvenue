package core

type Ticket struct {
	ID             string
	TenantID       string
	OrderID        string
	CustomerID     string
	Name           string
	QRCodeBase64   []byte
	QRCodeFileType string
	QRCodeURL      string
}
