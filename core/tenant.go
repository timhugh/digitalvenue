package core

type Tenant struct {
	TenantID string
	Name     string
	Meta     TenantMeta
}

type TenantMeta struct {
	SquareMerchantID string
}
