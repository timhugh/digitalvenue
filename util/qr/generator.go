package qr

import (
	"github.com/skip2/go-qrcode"
	"github.com/timhugh/digitalvenue/util/core"
)

const (
	recoveryLevel = qrcode.Medium
	fileType      = "png"
)

type Generator struct {
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Encode(data string, size int) (*core.QRCode, error) {
	png, err := qrcode.Encode(data, recoveryLevel, size)
	return &core.QRCode{
		Image:    png,
		FileType: fileType,
	}, err
}
