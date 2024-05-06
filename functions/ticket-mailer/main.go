package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/core"
	"github.com/timhugh/digitalvenue/util/logger"
	"gopkg.in/gomail.v2"
	"os"
	"path"
	"strings"
	"text/template"
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
	log                     *logger.ContextLogger
	tenantRepo              core.TenantRepository
	orderRepo               core.OrderRepository
	customerRepo            core.CustomerRepository
	ticketRepo              core.TicketRepository
	templateStore           core.TemplateStore
	tenantFilesBucketName   string
	tenantFilesBucketRegion string
}

func NewTicketMailerHandler(
	log *logger.ContextLogger,
	tenantRepo core.TenantRepository,
	orderRepo core.OrderRepository,
	customerRepo core.CustomerRepository,
	templateStore core.TemplateStore,
	ticketRepo core.TicketRepository,
) (*TicketMailerHandler, error) {
	tenantFilesBucketName, err := core.RequireEnv("S3_TENANT_FILES_BUCKET_NAME")
	if err != nil {
		return nil, err
	}

	tenantFilesBucketRegion, err := core.RequireEnv("AWS_REGION")
	if err != nil {
		return nil, err
	}

	return &TicketMailerHandler{
		log:                     log,
		tenantRepo:              tenantRepo,
		orderRepo:               orderRepo,
		customerRepo:            customerRepo,
		templateStore:           templateStore,
		ticketRepo:              ticketRepo,
		tenantFilesBucketName:   tenantFilesBucketName,
		tenantFilesBucketRegion: tenantFilesBucketRegion,
	}, nil
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
		log.AddParams(map[string]interface{}{
			"tenantID": tenantID,
			"orderID":  orderID,
		})

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

		tickets, err := h.ticketRepo.GetTickets(tenant.TenantID, order.ID)
		if err != nil {
			log.AddParam("error", err.Error()).Error("Failed to get tickets. Queueing for retry")
			retryFailures = append(retryFailures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})
			continue
		}

		emailTemplate, err := h.templateStore.Get(tenant.TenantID, "ticketEmail.html")
		if err != nil {
			log.AddParam("error", err.Error()).Error("Failed to get email template. Queueing for retry")
			retryFailures = append(retryFailures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})
			continue
		}

		email, err := buildEmail(emailTemplate, order, customer, tenant, tickets, h.tenantFilesBucketName, h.tenantFilesBucketRegion)
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

type emailTemplateParams struct {
	FileBucketURL string
	Tickets       []emailTicketTemplateParams
}

type emailTicketTemplateParams struct {
	Name           string
	QrCodeImageURL string
}

func buildEmail(emailTemplate *core.Template, order *core.Order, customer *core.Customer, tenant *core.Tenant, tickets []*core.Ticket, tenantFileBucketName string, tenantFileBucketRegion string) (*gomail.Message, error) {
	ticketTemplate, err := template.New("ticketEmail").Parse(emailTemplate.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse email template")
	}

	var emailBody strings.Builder
	tenantBucketURL := path.Join(fmt.Sprintf("https://s3-%s.amazonaws.com", tenantFileBucketRegion), tenantFileBucketName, tenant.TenantID)
	err = ticketTemplate.Execute(&emailBody, buildEmailTemplateParams(order, customer, tenant, tickets, tenantBucketURL))
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute email template")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", tenant.SMTPUser)
	m.SetHeader("To", customer.Email)
	m.SetHeader("Subject", fmt.Sprintf("Your tickets for %s", tenant.Name))
	m.SetBody("text/html", emailBody.String())
	return m, nil
}

func buildEmailTemplateParams(order *core.Order, customer *core.Customer, tenant *core.Tenant, tickets []*core.Ticket, tenantFileBucketURL string) *emailTemplateParams {
	params := emailTemplateParams{
		FileBucketURL: tenantFileBucketURL,
	}
	for _, ticket := range tickets {
		params.Tickets = append(params.Tickets, emailTicketTemplateParams{
			Name:           ticket.Name,
			QrCodeImageURL: ticket.QRCodeURL,
		})
	}
	return &params
}
