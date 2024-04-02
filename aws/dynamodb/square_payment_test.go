package dynamodb

import (
	"github.com/timhugh/digitalvenue/square/db"
	"testing"
)

func TestPaymentsRepository_ImplementsInterface(t *testing.T) {
	var _ db.SquarePaymentsRepository = &SquarePaymentsRepository{}
}
