package tenants

import "fmt"

const defaultApiVersion = "2019-01-01-preview"

func userAgent() string {
	return fmt.Sprintf("pandora/tenants/%s", defaultApiVersion)
}
