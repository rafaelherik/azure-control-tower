package azure

import (
	"context"
	"fmt"

	"azure-control-tower/internal/models"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

// ListResourceGroups returns all resource groups in the specified subscription
func (c *Client) ListResourceGroups(ctx context.Context, subscriptionID string) ([]*models.ResourceGroup, error) {
	client, err := armresources.NewResourceGroupsClient(subscriptionID, c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource groups client: %w", err)
	}

	pager := client.NewListPager(nil)

	var resourceGroups []*models.ResourceGroup
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page: %w", err)
		}

		for _, rg := range page.Value {
			if rg == nil {
				continue
			}

			location := ""
			if rg.Location != nil {
				location = *rg.Location
			}

			resourceGroups = append(resourceGroups, &models.ResourceGroup{
				Name:     *rg.Name,
				Location: location,
				Tags:     rg.Tags,
			})
		}
	}

	return resourceGroups, nil
}
