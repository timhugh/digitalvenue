package square

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/core"
	"github.com/timhugh/digitalvenue/util/logger"
	"strconv"
)

const (
	OrderIDKey    = "SquareOrderID"
	ItemIDKey     = "SquareItemID"
	MerchantIDKey = "SquareMerchantID"
	PaymentIDKey  = "SquarePaymentID"
)

type Order struct {
	SquareOrderID    string      `json:"id"`
	SquareCustomerID string      `json:"customer_id"`
	SquareLocationID string      `json:"location_id"`
	OrderItems       []OrderItem `json:"line_items"`
}

type OrderItem struct {
	ItemID          string `json:"uid"`
	Name            string `json:"name"`
	Quantity        string `json:"quantity"`
	CatalogObjectID string `json:"catalog_object_id"`
}

type OrderBuilder struct {
	log       *logger.ContextLogger
	squareAPI APIClient
}

func NewOrderBuilder(
	log *logger.ContextLogger,
	squareAPI APIClient,
) *OrderBuilder {
	return &OrderBuilder{
		squareAPI: squareAPI,
		log:       log,
	}
}

func (b *OrderBuilder) BuildOrder(squareOrder *Order, merchant *Merchant, squarePaymentID string, customerID string) (*core.Order, error) {
	order := core.Order{
		ID:         squareOrder.SquareOrderID, // We use the SquareOrderID as the ID for orders that come from Square to make duplicate detection easier
		TenantID:   merchant.TenantID,
		CustomerID: customerID,
		Meta: map[string]string{
			MerchantIDKey: merchant.ID,
			PaymentIDKey:  squarePaymentID,
			OrderIDKey:    squareOrder.SquareOrderID,
			CustomerIDKey: squareOrder.SquareCustomerID,
		},
	}

	var items []core.OrderItem
	for _, item := range squareOrder.OrderItems {
		newItems, err := b.buildItems(&item, merchant)
		if err != nil {
			if errors.Is(err, NonTicketableItemError{}) {
				continue
			}

			return nil, errors.Wrap(err, "failure while building order item")
		}
		items = append(items, newItems...)
	}

	order.Items = items

	return &order, nil
}

type NonTicketableItemError struct {
	error
}

func (b *OrderBuilder) buildItems(squareOrderItem *OrderItem, merchant *Merchant) ([]core.OrderItem, error) {
	var items []core.OrderItem

	var catalogItemVariation CatalogItemVariation
	err := b.squareAPI.GetCatalogObject(squareOrderItem.CatalogObjectID, merchant.SquareAPIToken, &catalogItemVariation)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get catalog item variation")
	}

	var catalogItem CatalogItem
	err = b.squareAPI.GetCatalogObject(catalogItemVariation.ItemVariationData.ItemID, merchant.SquareAPIToken, &catalogItem)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get catalog item")
	}

	isTicketable := false
	for _, itemCategory := range catalogItem.ItemData.Categories {
		for _, ticketableCategory := range merchant.TicketableCategories {
			if itemCategory.ID == ticketableCategory {
				isTicketable = true
			}
		}
	}
	if !isTicketable {
		return nil, NonTicketableItemError{}
	}

	quantity, err := strconv.Atoi(squareOrderItem.Quantity)
	if err != nil {
		return nil, errors.Wrap(err, "invalid item quantity")
	}

	for i := 0; i < quantity; i++ {
		items = append(items, core.OrderItem{
			Name: squareOrderItem.Name,
			ID:   fmt.Sprintf("%s-%d", squareOrderItem.ItemID, i+1),
			Meta: map[string]string{
				ItemIDKey: squareOrderItem.ItemID,
			},
		})
	}

	return items, nil
}
