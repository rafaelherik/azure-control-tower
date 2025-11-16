package azure

import (
	"context"
	"fmt"

	"azure-control-tower/internal/models"
)

// ListSubscriptions returns all subscriptions accessible to the current user
func (c *Client) ListSubscriptions(ctx context.Context) ([]*models.Subscription, error) {
	pager := c.SubscriptionsClient.NewListPager(nil)

	var subscriptions []*models.Subscription
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page: %w", err)
		}

		for _, sub := range page.Value {
			if sub == nil {
				continue
			}

			state := "Unknown"
			if sub.State != nil {
				state = string(*sub.State)
			}

			displayName := ""
			if sub.DisplayName != nil {
				displayName = *sub.DisplayName
			}

			tenantID := ""
			if sub.TenantID != nil {
				tenantID = *sub.TenantID
			}

			subscriptions = append(subscriptions, &models.Subscription{
				ID:          *sub.SubscriptionID,
				Name:        displayName,
				State:       state,
				DisplayName: displayName,
				TenantID:    tenantID,
			})
		}
	}

	return subscriptions, nil
}
