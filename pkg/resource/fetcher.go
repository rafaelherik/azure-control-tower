package resource

import (
	"context"
)

// Fetcher defines the interface for fetching resources
type Fetcher interface {
	FetchResources(ctx context.Context, subscriptionID string, resourceType string) ([]interface{}, error)
}

