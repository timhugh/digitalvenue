package persistence

import "github.com/timhugh/digitalvenue/core"

type MerchantRepo interface {
	CreateMerchant(merchant *core.Merchant) error
	FindMerchantBySquareMerchantId(squareMerchantId string) (*core.Merchant, error)
}
