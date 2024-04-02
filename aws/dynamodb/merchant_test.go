package dynamodb

import (
	"github.com/timhugh/digitalvenue/core/db"
	"testing"
)

func TestMerchantsRepository_ImplementsInterface(t *testing.T) {
	var _ db.MerchantsRepository = &MerchantsRepository{}
}
