package dynamodb

import (
	"github.com/timhugh/digitalvenue/db"
	"testing"
)

func TestMerchantsRepository_ImplementsInterface(t *testing.T) {
	var _ db.MerchantsRepository = &MerchantsRepository{}
}
