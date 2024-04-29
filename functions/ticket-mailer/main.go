package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/timhugh/digitalvenue/util/core"
	"github.com/timhugh/digitalvenue/util/logger"
	"gopkg.in/gomail.v2"
	"os"
	"strings"
)

func main() {
	log := logger.Default().AddParam("service", "ticket-generator")
	env, err := core.RequireEnv("ENVIRONMENT")
	if err != nil {
		log.AddParam("error", err.Error()).Fatal("Failed to determine application environment")
		os.Exit(1)
	}
	log.AddParam("environment", env)

	handler, err := initializeHandler(log)
	if err != nil {
		log.AddParam("error", err.Error()).Fatal("Failed to initialize handler")
	}
	lambda.Start(handler.Handle)
}

type TicketMailerHandler struct {
	log          *logger.ContextLogger
	tenantRepo   core.TenantRepository
	orderRepo    core.OrderRepository
	customerRepo core.CustomerRepository
}

func NewTicketMailerHandler(log *logger.ContextLogger, tenantRepo core.TenantRepository, orderRepo core.OrderRepository, customerRepo core.CustomerRepository) *TicketMailerHandler {
	return &TicketMailerHandler{
		log:          log,
		tenantRepo:   tenantRepo,
		orderRepo:    orderRepo,
		customerRepo: customerRepo,
	}
}

func (h *TicketMailerHandler) Handle(event events.SQSEvent) (events.SQSEventResponse, error) {
	var retryFailures []events.SQSBatchItemFailure

	for _, record := range event.Records {
		log := h.log.Sub().AddParam("messageID", record.MessageId)

		messageParams := strings.Split(record.Body, ":")
		if len(messageParams) != 2 {
			log.Error("Received invalid event: '%s'", record.Body)
			continue // not retryable
		}
		tenantID := messageParams[0]
		orderID := messageParams[1]

		tenant, err := h.tenantRepo.GetTenant(tenantID)
		if err != nil {
			log.AddParam("error", err.Error()).Error("Failed to get tenant. Queueing for retry")
			retryFailures = append(retryFailures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})
			continue
		}

		order, err := h.orderRepo.GetOrder(tenantID, orderID)
		if err != nil {
			log.AddParam("error", err.Error()).Error("Failed to get order. Queueing for retry")
			retryFailures = append(retryFailures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})
			continue
		}

		customer, err := h.customerRepo.GetCustomer(tenant.TenantID, order.CustomerID)
		if err != nil {
			log.AddParam("error", err.Error()).Error("Failed to get customer. Queueing for retry")
			retryFailures = append(retryFailures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})
			continue
		}

		email, err := buildEmail(order, customer, tenant)
		if err != nil {
			log.AddParam("error", err.Error()).Error("Failed to build email. Queueing for retry")
			retryFailures = append(retryFailures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})
			continue
		}

		dialer := gomail.NewDialer(tenant.SMTPHost, tenant.SMTPPort, tenant.SMTPUser, tenant.SMTPPassword)
		if err := dialer.DialAndSend(email); err != nil {
			log.AddParam("error", err.Error()).Error("Failed to send email. Queueing for retry")
			retryFailures = append(retryFailures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})
			continue
		}
	}

	response := events.SQSEventResponse{}
	if len(retryFailures) > 0 {
		response.BatchItemFailures = retryFailures
	}
	return response, nil
}

func buildEmail(order *core.Order, customer *core.Customer, tenant *core.Tenant) (*gomail.Message, error) {
	m := gomail.NewMessage()
	m.SetHeader("From", tenant.SMTPUser)
	m.SetHeader("To", customer.Email)
	m.SetHeader("Subject", fmt.Sprintf("Your tickets for %s", tenant.Name))
	m.SetBody("text/html", "Here's them tickets")
	return m, nil
}
