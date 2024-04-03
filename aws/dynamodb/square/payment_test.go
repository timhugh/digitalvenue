package square

import (
	"github.com/timhugh/digitalvenue/square"
	"testing"
)

func TestPaymentsRepository_ImplementsInterface(t *testing.T) {
	var _ square.PaymentsRepository = &PaymentsRepository{}
}
