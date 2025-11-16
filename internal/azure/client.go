package azure

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
)

// Client wraps Azure SDK clients
type Client struct {
	SubscriptionsClient *armsubscriptions.Client
	credential          azcore.TokenCredential
}

// NewClient creates a new Azure client wrapper
func NewClient(credential azcore.TokenCredential) (*Client, error) {
	subscriptionsClient, err := armsubscriptions.NewClient(credential, nil)
	if err != nil {
		return nil, err
	}

	return &Client{
		SubscriptionsClient: subscriptionsClient,
		credential:          credential,
	}, nil
}

