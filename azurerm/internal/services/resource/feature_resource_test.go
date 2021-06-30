package resource_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type FeatureResource struct {
}

func TestAccFeatureResource_basic(t *testing.T) {
	if true {
		t.Skip("Skipping due to acceptance test doesn't support having features toggled currently")
	}
	data := acceptance.BuildTestData(t, "azurerm_feature", "test")
	r := FeatureResource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccFeatureResource_requiresImport(t *testing.T) {
	if true {
		t.Skip("Skipping due to acceptance test doesn't support having features toggled currently")
	}
	data := acceptance.BuildTestData(t, "azurerm_feature", "test")
	r := FeatureResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.requiresImportBasic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r FeatureResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FeatureID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Resource.FeaturesClient.Get(ctx, id.ProviderNamespace, id.Name)
	if err != nil {
		return nil, fmt.Errorf("error retrieving feature %q: %+v", id, err)
	}
	if resp.Properties != nil && resp.Properties.State != nil {
		if strings.EqualFold(*resp.Properties.State, "Pending") {
			return nil, fmt.Errorf("feature (%q) which requires manual approval can not be managed by terraform", id)
		}
		if !strings.EqualFold(*resp.Properties.State, "Registered") {
			return utils.Bool(false), nil
		}
	}
	return utils.Bool(true), nil
}

func (r FeatureResource) basic() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_feature" "test" {
  name               = "AutoApproveFeature"
  provider_namespace = "Microsoft.CognitiveServices"
}
`
}

func (r FeatureResource) requiresImportBasic() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_feature" "test" {
  name               = "AllowManagedDisksReplaceOSDisk"
  provider_namespace = "Microsoft.Compute"
}
`
}

func (r FeatureResource) requiresImport(data acceptance.TestData) string {
	config := r.requiresImportBasic()
	return fmt.Sprintf(`
%s

resource "azurerm_feature" "import" {
  name               = azurerm_feature.test.name
  provider_namespace = azurerm_feature.test.provider_namespace
}
`, config)
}
