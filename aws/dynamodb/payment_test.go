package dynamodb

import (
	"github.com/timhugh/digitalvenue/core/db"
	"testing"
)

func TestPaymentsRepository_ImplementsInterface(t *testing.T) {
	var _ db.PaymentsRepository = &PaymentsRepository{}
}
