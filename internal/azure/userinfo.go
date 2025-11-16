package azure

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"azure-control-tower/internal/models"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// GetUserInfo extracts user and tenant information from the Azure token
func (c *Client) GetUserInfo(ctx context.Context) (*models.UserInfo, error) {
	opts := policy.TokenRequestOptions{
		Scopes: []string{"https://management.azure.com/.default"},
	}

	token, err := c.credential.GetToken(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	// Parse JWT token
	parts := strings.Split(token.Token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	// Decode the payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode token payload: %w", err)
	}

	// Parse JSON claims
	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse token claims: %w", err)
	}

	userInfo := &models.UserInfo{}

	// Extract tenant ID
	if tid, ok := claims["tid"].(string); ok {
		userInfo.TenantID = tid
	}

	// Extract user name (prefer "name", fallback to "preferred_username")
	if name, ok := claims["name"].(string); ok && name != "" {
		userInfo.Name = name
	} else if preferredUsername, ok := claims["preferred_username"].(string); ok {
		userInfo.Name = preferredUsername
	}

	// Extract email (prefer "upn", fallback to "email", then "preferred_username")
	if upn, ok := claims["upn"].(string); ok && upn != "" {
		userInfo.Email = upn
	} else if email, ok := claims["email"].(string); ok && email != "" {
		userInfo.Email = email
	} else if preferredUsername, ok := claims["preferred_username"].(string); ok {
		userInfo.Email = preferredUsername
	}

	return userInfo, nil
}
