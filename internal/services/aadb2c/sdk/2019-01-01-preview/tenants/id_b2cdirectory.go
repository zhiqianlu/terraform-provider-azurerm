package tenants

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type B2CDirectoryId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewB2CDirectoryID(subscriptionId, resourceGroup, name string) B2CDirectoryId {
	return B2CDirectoryId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id B2CDirectoryId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "B2 C Directory", segmentsStr)
}

func (id B2CDirectoryId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AzureActiveDirectory/b2cDirectories/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// ParseB2CDirectoryID parses a B2CDirectory ID into an B2CDirectoryId struct
func ParseB2CDirectoryID(input string) (*B2CDirectoryId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := B2CDirectoryId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Name, err = id.PopSegment("b2cDirectories"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseB2CDirectoryIDInsensitively parses an B2CDirectory ID into an B2CDirectoryId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseB2CDirectoryID method should be used instead for validation etc.
func ParseB2CDirectoryIDInsensitively(input string) (*B2CDirectoryId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := B2CDirectoryId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'b2cDirectories' segment
	b2cDirectoriesKey := "b2cDirectories"
	for key := range id.Path {
		if strings.EqualFold(key, b2cDirectoriesKey) {
			b2cDirectoriesKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(b2cDirectoriesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
