package aadb2c

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/aadb2c/sdk/2019-01-01-preview/tenants"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/aadb2c/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AadB2cDirectoryModel struct {
	BillingType           string            `tfschema:"billing_type"`
	CountryCode           string            `tfschema:"country_code"`
	DisplayName           string            `tfschema:"display_name"`
	EffectiveStartDate    string            `tfschema:"effective_start_date"`
	DataResidencyLocation string            `tfschema:"data_residency_location"`
	Name                  string            `tfschema:"name"`
	ResourceGroup         string            `tfschema:"resource_group_name"`
	Sku                   string            `tfschema:"sku_name"`
	Tags                  map[string]string `tfschema:"tags"`
}

var _ sdk.Resource = AadB2cDirectoryResource{}
var _ sdk.ResourceWithUpdate = AadB2cDirectoryResource{}

type AadB2cDirectoryResource struct{}

func (r AadB2cDirectoryResource) ResourceType() string {
	return "azurerm_aadb2c_directory"
}

func (r AadB2cDirectoryResource) ModelObject() interface{} {
	return &AadB2cDirectoryModel{}
}

func (r AadB2cDirectoryResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.B2CDirectoryID
}

func (r AadB2cDirectoryResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"data_residency_location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(tenants.LocationAsiaPacific),
				string(tenants.LocationAustralia),
				string(tenants.LocationEurope),
				string(tenants.LocationGlobal),
				string(tenants.LocationUnitedStates),
			}, false),
		},

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(tenants.SkuNamePremiumP1),
				string(tenants.SkuNamePremiumP2),
				string(tenants.SkuNameStandard),
			}, false),
		},

		"tags": tags.Schema(),
	}
}

func (r AadB2cDirectoryResource) Attributes() map[string]*pluginsdk.Schema {
	return nil
}

func (r AadB2cDirectoryResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AadB2c.TenantsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AadB2cDirectoryModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := tenants.NewB2CDirectoryID(subscriptionId, model.ResourceGroup, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && existing.HttpResponse.StatusCode != http.StatusNotFound {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			metadata.Logger.Infof("Creating %s", id)

			properties := tenants.CreateTenant{
				Location: tenants.Location(model.DataResidencyLocation),
				Properties: tenants.CreateTenantProperties{
					CountryCode: model.CountryCode,
					DisplayName: model.DisplayName,
				},
				Sku: tenants.Sku{
					Name: tenants.SkuName(model.Sku),
					Tier: tenants.SkuTierA0,
				},
				Tags: &model.Tags,
			}

			if err := client.CreateThenPoll(ctx, id, properties); err != nil {
				return err
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r AadB2cDirectoryResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AadB2c.TenantsClient

			id, err := tenants.ParseB2CDirectoryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state AadB2cDirectoryModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			metadata.Logger.Infof("Updating %s", id)

			properties := tenants.UpdateTenant{
				Properties: tenants.UpdateTenantProperties{
					BillingConfig: &tenants.BillingConfig{
						BillingType: (*tenants.BillingType)(&state.BillingType),
					},
				},
				Sku: tenants.Sku{
					Name: tenants.SkuName(state.Sku),
					Tier: tenants.SkuTierA0,
				},
				Tags: &state.Tags,
			}

			if _, err := client.Update(ctx, *id, properties); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r AadB2cDirectoryResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AadB2c.TenantsClient

			id, err := tenants.ParseB2CDirectoryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Reading %s", id)
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if resp.HttpResponse.StatusCode == http.StatusNotFound {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := AadB2cDirectoryModel{
				Name:          id.Name,
				ResourceGroup: id.ResourceGroup,
			}

			if model.Location != nil {
				state.DataResidencyLocation = string(*model.Location)
			}

			if model.Sku != nil {
				state.Sku = string(model.Sku.Name)
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			if properties := model.Properties; properties != nil {
				if billingConfig := properties.BillingConfig; billingConfig != nil {
					if billingConfig.BillingType != nil {
						state.BillingType = string(*billingConfig.BillingType)
					}
					if billingConfig.EffectiveStartDateUtc != nil {
						state.EffectiveStartDate = *billingConfig.EffectiveStartDateUtc
					}
				}

				if properties.CountryCode != nil {
					state.CountryCode = *properties.CountryCode
				}
				if properties.DisplayName != nil {
					state.DisplayName = *properties.DisplayName
				}
			}

			return nil
		},
	}
}

func (r AadB2cDirectoryResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AadB2c.TenantsClient

			id, err := tenants.ParseB2CDirectoryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Deleting %s", id)

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
