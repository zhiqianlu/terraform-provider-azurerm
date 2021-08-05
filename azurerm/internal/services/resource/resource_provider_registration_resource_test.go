package resource_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// NOTE: the RP's in this test shouldn't be used elsewhere in the provider, as such we're picking ones we're not
// using or third-party RP's for this purpose.

type ResourceProviderRegistrationResource struct {
}

func TestAccResourceProviderRegistration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_provider_registration", "test")
	r := ResourceProviderRegistrationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic("Microsoft.BlockchainTokens"),
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
			Config: r.basic("Raygun.CrashReporting"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(func(data acceptance.TestData) string {
			return r.requiresImport("Raygun.CrashReporting")
		}),
	})
}

func TestAccResourceProviderRegistration_feature(t *testing.T) {
	if os.Getenv("TF_RP_FEATURES_TEST") == "" {
		t.Skip("Skipping since `TF_RP_FEATURES_TEST` is unset")
	}

	data := acceptance.BuildTestData(t, "azurerm_resource_provider_registration", "test")
	r := ResourceProviderRegistrationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		//{
		//	// first ensure the existing/defaulted to on RP is unregistered
		//	Config: r.emptyProvider(),
		//	Check: acceptance.ComposeTestCheckFunc(
		//		data.CheckWithClientWithoutResource(r.ensureResourceProviderIsUnregistered("Microsoft.Compute")),
		//	),
		//},
		{
			Config: r.basic("Microsoft.Compute"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("feature.#").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withFeatures(false, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withFeatures(true, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic("Microsoft.Compute"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (ResourceProviderRegistrationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	name := state.Attributes["name"]
	resp, err := client.Resource.ProvidersClient.Get(ctx, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("Bad: Get on ProvidersClient: %+v", err)
	}

	return utils.Bool(resp.RegistrationState != nil && strings.EqualFold(*resp.RegistrationState, "Registered")), nil
}

func (ResourceProviderRegistrationResource) basic(name string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  skip_provider_registration = true
}

resource "azurerm_resource_provider_registration" "test" {
  name = %q
}
`, name)
}

func (r ResourceProviderRegistrationResource) requiresImport(name string) string {
	template := r.basic(name)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_provider_registration" "import" {
  name = azurerm_resource_provider_registration.test.name
}
`, template)
}

func (r ResourceProviderRegistrationResource) withFeatures(firstEnabled bool, secondEnabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  skip_provider_registration = true
}

resource "azurerm_resource_provider_registration" "test" {
  name = "Microsoft.Compute"

  feature {
    name       = "InGuestPatchVMPreview"
    registered = %t
  }

  feature {
    name       = "InGuestHotPatchVMPreview"
    registered = %t
  }
}
`, firstEnabled, secondEnabled)
}

func (r ResourceProviderRegistrationResource) emptyProvider() string {
	return `
provider "azurerm" {
  features {}
}
`
}

func (r ResourceProviderRegistrationResource) ensureResourceProviderIsUnregistered(resourceProvider string) func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
	return func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
		client := clients.Resource.ResourceProvidersClient
		isRegistered := true // it is by default, so assume so
		resp, err := client.Get(ctx, resourceProvider, "")
		if err != nil {
			return fmt.Errorf("retrieving Resource Provider %q: %+v", resourceProvider, err)
		}
		if resp.RegistrationState != nil {
			isRegistered = strings.EqualFold(*resp.RegistrationState, resource.Registered)
		}
		if isRegistered {
			if _, err := client.Unregister(ctx, resourceProvider); err != nil {
				return fmt.Errorf("unregistering Resource Provider %q: %+v", resourceProvider, err)
			}

			stateConf := &pluginsdk.StateChangeConf{
				Pending: []string{"Processing", "Registered"},
				Target:  []string{"Unregistered"},
				Refresh: func() (interface{}, string, error) {
					resp, err := client.Get(ctx, resourceProvider, "")
					if err != nil {
						return resp, "Failed", err
					}

					if resp.RegistrationState != nil && strings.EqualFold(*resp.RegistrationState, "Unregistered") {
						return resp, "Unregistered", nil
					}

					return resp, "Processing", nil
				},
				MinTimeout: 15 * time.Second,
				Timeout:    30 * time.Minute,
			}
			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for Resource Provider %q to become unregistered: %+v", resourceProvider, err)
			}
		}

		return nil
	}
}
