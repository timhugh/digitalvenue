package db

import "github.com/timhugh/digitalvenue/core"

type MerchantsRepository interface {
	CreateMerchant(merchant core.Merchant) error
	FindMerchantBySquareMerchantId(squareMerchantId string) (core.Merchant, error)
}
