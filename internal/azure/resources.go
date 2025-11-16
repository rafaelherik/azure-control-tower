package azure

import (
	"context"
	"fmt"

	"azure-control-tower/internal/models"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

// ListResources returns all resources in the specified subscription, optionally filtered by resource type
func (c *Client) ListResources(ctx context.Context, subscriptionID string, resourceType string) ([]*models.Resource, error) {
	client, err := armresources.NewClient(subscriptionID, c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create resources client: %w", err)
	}

	var options *armresources.ClientListOptions
	if resourceType != "" {
		// Filter by resource type if specified
		filter := fmt.Sprintf("resourceType eq '%s'", resourceType)
		options = &armresources.ClientListOptions{
			Filter: &filter,
		}
	}

	pager := client.NewListPager(options)

	var resources []*models.Resource
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page: %w", err)
		}

		for _, resource := range page.Value {
			if resource == nil {
				continue
			}

			// Extract resource group from ID
			resourceGroup := extractResourceGroupFromID(*resource.ID)

			resourceType := ""
			if resource.Type != nil {
				resourceType = *resource.Type
			}

			name := ""
			if resource.Name != nil {
				name = *resource.Name
			}

			location := ""
			if resource.Location != nil {
				location = *resource.Location
			}

			res := &models.Resource{
				ID:            *resource.ID,
				Name:          name,
				Type:          resourceType,
				Location:      location,
				ResourceGroup: resourceGroup,
				Tags:          resource.Tags,
				Properties:    make(map[string]interface{}),
			}

			// Convert properties if available
			if resource.Properties != nil {
				// Properties is a map[string]interface{} in the SDK
				// We'll store it as-is for now
				if props, ok := resource.Properties.(map[string]interface{}); ok {
					res.Properties = props
				}
			}

			resources = append(resources, res)
		}
	}

	return resources, nil
}

// ListResourcesByResourceGroup returns all resources in a specific resource group
func (c *Client) ListResourcesByResourceGroup(ctx context.Context, subscriptionID, resourceGroupName string, resourceType string) ([]*models.Resource, error) {
	client, err := armresources.NewClient(subscriptionID, c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create resources client: %w", err)
	}

	var options *armresources.ClientListByResourceGroupOptions
	if resourceType != "" {
		filter := fmt.Sprintf("resourceType eq '%s'", resourceType)
		options = &armresources.ClientListByResourceGroupOptions{
			Filter: &filter,
		}
	}

	pager := client.NewListByResourceGroupPager(resourceGroupName, options)

	var resources []*models.Resource
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page: %w", err)
		}

		for _, resource := range page.Value {
			if resource == nil {
				continue
			}

			resourceType := ""
			if resource.Type != nil {
				resourceType = *resource.Type
			}

			name := ""
			if resource.Name != nil {
				name = *resource.Name
			}

			location := ""
			if resource.Location != nil {
				location = *resource.Location
			}

			res := &models.Resource{
				ID:            *resource.ID,
				Name:          name,
				Type:          resourceType,
				Location:      location,
				ResourceGroup: resourceGroupName,
				Tags:          resource.Tags,
				Properties:    make(map[string]interface{}),
			}

			if resource.Properties != nil {
				if props, ok := resource.Properties.(map[string]interface{}); ok {
					res.Properties = props
				}
			}

			resources = append(resources, res)
		}
	}

	return resources, nil
}

// GetResourceTypeCounts returns resource type summaries with counts for a resource group
func (c *Client) GetResourceTypeCounts(ctx context.Context, subscriptionID, resourceGroupName string) ([]*models.ResourceTypeSummary, error) {
	// Get all resources in the resource group
	resources, err := c.ListResourcesByResourceGroup(ctx, subscriptionID, resourceGroupName, "")
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %w", err)
	}

	// Aggregate by resource type
	typeCounts := make(map[string]int)
	for _, resource := range resources {
		if resource.Type != "" {
			typeCounts[resource.Type]++
		}
	}

	// Convert to ResourceTypeSummary slice
	summaries := make([]*models.ResourceTypeSummary, 0, len(typeCounts))
	for resourceType, count := range typeCounts {
		summaries = append(summaries, &models.ResourceTypeSummary{
			Type:  resourceType,
			Count: count,
		})
	}

	return summaries, nil
}

// extractResourceGroupFromID extracts the resource group name from a resource ID
// Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/{provider}/{type}/{name}
func extractResourceGroupFromID(id string) string {
	parts := splitResourceID(id)
	if len(parts) >= 4 && parts[0] == "subscriptions" && parts[2] == "resourceGroups" {
		return parts[3]
	}
	return ""
}

// splitResourceID splits a resource ID into its components
func splitResourceID(id string) []string {
	if len(id) == 0 || id[0] != '/' {
		return []string{}
	}

	parts := []string{}
	current := ""
	for i := 1; i < len(id); i++ {
		if id[i] == '/' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(id[i])
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
