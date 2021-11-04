package aadb2c_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/aadb2c/sdk/2019-01-01-preview/tenants"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AadB2cDirectoryResource struct{}

func TestAccAadB2cDirectoryResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_aadb2c_directory", "test")
	r := AadB2cDirectoryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r AadB2cDirectoryResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := tenants.ParseB2CDirectoryID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.AadB2c.TenantsClient.Get(ctx, *id)
	if err != nil {
		if resp.HttpResponse.StatusCode == http.StatusNotFound {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r AadB2cDirectoryResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_aadb2c_directory" "test" {
  name                    = "acctest%[1]d"
  resource_group_name     = azurerm_resource_group.test.name
  data_residency_location = "europe"
  sku_name                = "Standard"
}
`, data.RandomInteger, data.Locations.Primary)
}
