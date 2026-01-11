package domain

// TenantProvisionRequest represents the payload to provision a new tenant
type TenantProvisionRequest struct {
	ID             string `json:"id"`
	Subdomain      string `json:"subdomain"`
	ComposeContent string `json:"compose_content"`
	EnvContent     string `json:"env_content"`
}

// TenantActionRequest represents a simple action on an existing tenant
type TenantActionRequest struct {
	ID string `json:"id"`
}

// TenantStatus represents the current state of a tenant's containers
type TenantStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Uptime string `json:"uptime"`
}
