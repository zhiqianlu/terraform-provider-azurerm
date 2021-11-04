package tenants

type CreateTenant struct {
	Location   Location               `json:"location"`
	Properties CreateTenantProperties `json:"properties"`
	Sku        Sku                    `json:"sku"`
	Tags       *map[string]string     `json:"tags,omitempty"`
}
