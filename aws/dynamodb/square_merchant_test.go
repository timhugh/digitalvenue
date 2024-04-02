package dynamodb

import (
	"github.com/timhugh/digitalvenue/square/db"
	"testing"
)

func TestMerchantsRepository_ImplementsInterface(t *testing.T) {
	var _ db.SquareMerchantsRepository = &SquareMerchantsRepository{}
}
