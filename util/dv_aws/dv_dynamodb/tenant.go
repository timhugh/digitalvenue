package dv_dynamodb

import (
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/core"
	"maps"
)

type tenant struct {
	PK   string
	SK   string
	Name string

	EmailsEnabled   bool
	SMTPAccount     string
	SMTPPassword    string
	SMTPHost        string
	SMTPPort        int
	SMTPFromAddress string

	Meta map[string]string
}

func (repo *Repository) GetTenant(tenantID string) (*core.Tenant, error) {
	tenantKey := PrefixID("Tenant", tenantID)
	key := map[string]string{
		"PK": tenantKey,
		"SK": tenantKey,
	}

	item := tenant{}
	err := repo.get("Tenant", key, &item)
	if err != nil {
		if errors.Is(err, ItemNotFoundException{}) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to get Tenant")
	}

	return &core.Tenant{
		TenantID: tenantID,
		Name:     item.Name,

		EmailsEnabled:   item.EmailsEnabled,
		SMTPUser:        item.SMTPAccount,
		SMTPPassword:    item.SMTPPassword,
		SMTPHost:        item.SMTPHost,
		SMTPPort:        item.SMTPPort,
		SMTPFromAddress: item.SMTPFromAddress,

		Meta: maps.Clone(item.Meta),
	}, nil
}
