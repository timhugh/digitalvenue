package square

import (
	"github.com/timhugh/digitalvenue/square"
	"testing"
)

func TestMerchantsRepository_ImplementsInterface(t *testing.T) {
	var _ square.MerchantsRepository = &MerchantsRepository{}
}
