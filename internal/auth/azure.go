package auth

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

// NewAzureAuth creates a new Azure credential using Azure CLI default credentials
func NewAzureAuth() (*azidentity.DefaultAzureCredential, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure credential: %w", err)
	}

	// Verify credentials by getting a token
	opts := policy.TokenRequestOptions{
		Scopes: []string{"https://management.azure.com/.default"},
	}
	_, err = cred.GetToken(context.Background(), opts)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Azure: %w. Please run 'az login'", err)
	}

	return cred, nil
}

