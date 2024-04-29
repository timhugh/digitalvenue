package test

import (
	"github.com/timhugh/digitalvenue/util/core"
)

const (
	TenantID           = "testTenant"
	TenantName         = "Test Tenant"
	EmailsEnabled      = false
	TenantSMTPAccount  = "tenant@email.com"
	TenantSMTPPassword = "besttenantever"
	TenantSMTPHost     = "smtp.email.com"
	TenantSMTPPort     = 465

	CustomerID        = "customerID"
	CustomerName      = "Tim Heuett"
	CustomerFirstName = "Tim"
	CustomerLastName  = "Heuett"
	CustomerEmail     = "info@timheuett.com"
	CustomerPhone     = "+12062062062"

	OrderID   = "orderID"
	ItemID1   = "itemID1"
	ItemName1 = "Item 1"
	ItemID2   = "itemID2"
	ItemName2 = "Item 2"
	ItemID3   = "itemID3"

	QRCodeBucket = "qr-code-bucket"
	QRCodeImage  = "image"
	QRCodeType   = "png"
)

func NewTenant() *core.Tenant {
	return &core.Tenant{
		TenantID: TenantID,
		Name:     TenantName,

		EmailsEnabled: EmailsEnabled,
		SMTPUser:      TenantSMTPAccount,
		SMTPPassword:  TenantSMTPPassword,
		SMTPHost:      TenantSMTPHost,
		SMTPPort:      TenantSMTPPort,

		Meta: map[string]string{},
	}
}

func NewOrder() *core.Order {
	return &core.Order{
		ID:         OrderID,
		TenantID:   TenantID,
		CustomerID: CustomerID,
		Meta:       map[string]string{},
		Items: []core.OrderItem{
			{
				ID:   ItemID1,
				Name: ItemName1,
				Meta: map[string]string{},
			},
			{
				ID:   ItemID2,
				Name: ItemName2,
				Meta: map[string]string{},
			},
			{
				ID:   ItemID3,
				Name: ItemName2,
				Meta: map[string]string{},
			},
		},
	}
}

func NewCustomer() *core.Customer {
	return &core.Customer{
		TenantID: TenantID,
		ID:       CustomerID,
		Name:     CustomerName,
		Email:    CustomerEmail,
		Phone:    CustomerPhone,
		Meta:     map[string]string{},
	}
}

func NewQRCode() *core.QRCode {
	return &core.QRCode{
		TenantID:    TenantID,
		OrderID:     OrderID,
		OrderItemID: ItemID1,
		Image:       []byte(QRCodeImage),
		FileType:    QRCodeType,
	}
}
