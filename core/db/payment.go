package db

import "github.com/timhugh/digitalvenue/core"

type PaymentsRepository interface {
	CreatePayment(payment core.Payment) error
}
