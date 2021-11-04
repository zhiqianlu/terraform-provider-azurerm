package tenants

type BillingType string

const (
	BillingTypeAuths BillingType = "auths"
	BillingTypeMAU   BillingType = "mau"
)

type Location string

const (
	LocationAsiaPacific  Location = "asiapacific"
	LocationAustralia    Location = "australia"
	LocationEurope       Location = "europe"
	LocationGlobal       Location = "global,unitedstates,europe,asiapacific,australia"
	LocationUnitedStates Location = "unitedstates"
)

type SkuName string

const (
	SkuNamePremiumP1 SkuName = "PremiumP1"
	SkuNamePremiumP2 SkuName = "PremiumP2"
	SkuNameStandard  SkuName = "Standard"
)

type SkuTier string

const (
	SkuTierA0 SkuTier = "A0"
)
