package resource_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// NOTE: this can be moved up a level when all the others are

type ResourceProviderRegistrationResource struct {
}

func TestAccResourceProviderRegistration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_provider_registration", "test")
	r := ResourceProviderRegistrationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "Microsoft.BlockchainTokens"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceProviderRegistration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_provider_registration", "test")
	r := ResourceProviderRegistrationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "Wandisco.Fusion"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(func(data acceptance.TestData) string {
			return r.requiresImport(data, "Wandisco.Fusion")
		}),
	})
}

func TestAccResourceProviderRegistration_feature(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_provider_registration", "test")
	r := ResourceProviderRegistrationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "Microsoft.IoTCentral"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("feature.#").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.feature(data, "Microsoft.IoTCentral", "scootaloo"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("feature.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.multiFeatures(data, "Microsoft.IoTCentral", "scootaloo", "metering"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("feature.#").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.feature(data, "Microsoft.IoTCentral", "metering"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("feature.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.noFeatures(data, "Microsoft.IoTCentral"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("feature.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func (ResourceProviderRegistrationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	name := state.Attributes["name"]
	client.Resource.ProvidersClient.BaseClient.SubscriptionID = os.Getenv("ARM_SUBSCRIPTION_ID_ALT")
	resp, err := client.Resource.ProvidersClient.Get(ctx, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("Bad: Get on ProvidersClient: %+v", err)
	}

	return utils.Bool(resp.RegistrationState != nil && strings.EqualFold(*resp.RegistrationState, "Registered")), nil
}

func (ResourceProviderRegistrationResource) basic(data acceptance.TestData, name string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  skip_provider_registration = true
  subscription_id            = "%s"
}

resource "azurerm_resource_provider_registration" "test" {
  name = %q
}
`, data.Client().SubscriptionIDAlt, name)
}

func (r ResourceProviderRegistrationResource) requiresImport(data acceptance.TestData, name string) string {
	template := r.basic(data, name)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_provider_registration" "import" {
  name = azurerm_resource_provider_registration.test.name
}
`, template)
}

func (ResourceProviderRegistrationResource) noFeatures(data acceptance.TestData, name string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  skip_provider_registration = true
  subscription_id            = "%s"
}

resource "azurerm_resource_provider_registration" "test" {
  name    = %q
  feature = []
}
`, data.Client().SubscriptionIDAlt, name)
}

func (ResourceProviderRegistrationResource) feature(data acceptance.TestData, name string, feature string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  skip_provider_registration = true
  subscription_id            = "%s"
}

resource "azurerm_resource_provider_registration" "test" {
  name = %q
  feature {
    name = %q
  }
}
`, data.Client().SubscriptionIDAlt, name, feature)
}

func (ResourceProviderRegistrationResource) multiFeatures(data acceptance.TestData, name string, feature1 string, feature2 string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  skip_provider_registration = true
  subscription_id            = "%s"
}

resource "azurerm_resource_provider_registration" "test" {
  name = %q
  feature {
    name = %q
  }
  feature {
    name = %q
  }
}
`, data.Client().SubscriptionIDAlt, name, feature1, feature2)
}
