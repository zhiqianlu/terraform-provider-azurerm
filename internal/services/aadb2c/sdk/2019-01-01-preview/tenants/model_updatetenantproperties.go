package tenants

type UpdateTenantProperties struct {
	BillingConfig *BillingConfig `json:"billingConfig,omitempty"`
	TenantId      *string        `json:"tenantId,omitempty"`
}
