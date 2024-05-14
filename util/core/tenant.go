package core

type Tenant struct {
	TenantID string
	Name     string

	EmailsEnabled   bool
	SMTPUser        string
	SMTPPassword    string
	SMTPHost        string
	SMTPPort        int
	SMTPFromAddress string

	Meta map[string]string
}

type TenantRepository interface {
	GetTenant(tenantID string) (*Tenant, error)
}
