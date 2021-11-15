package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/communication/sdk/2020-08-20/communicationservice"
)

type Client struct {
	ServiceClient *communicationservice.CommunicationServiceClient
}

func NewClient(o *common.ClientOptions) *Client {
	serviceClient := communicationservice.NewCommunicationServiceClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&serviceClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ServiceClient: &serviceClient,
	}
}
