package dynamodb

import (
	"github.com/timhugh/digitalvenue/db"
	"testing"
)

func TestPaymentsRepository_ImplementsInterface(t *testing.T) {
	var _ db.PaymentsRepository = &PaymentsRepository{}
}
