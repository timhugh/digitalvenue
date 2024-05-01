package core

type Template struct {
	Key      string
	TenantID string
	Body     string
}

type TemplateStore interface {
	Save(template *Template) error
	Get(tenantID string, templateKey string) (*Template, error)
}
