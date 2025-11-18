package resource

import (
	"strings"
	"testing"

	"azure-control-tower/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestNewKeyVaultHandler(t *testing.T) {
	handler := NewKeyVaultHandler()
	assert.NotNil(t, handler)
	assert.Equal(t, "Microsoft.KeyVault/vaults", handler.GetResourceType())
}

func TestKeyVaultHandler_GetColumns(t *testing.T) {
	handler := NewKeyVaultHandler()
	columns := handler.GetColumns()
	
	assert.Len(t, columns, 3)
	assert.Equal(t, "Type", columns[0].Name)
	assert.Equal(t, "Name", columns[1].Name)
	assert.Equal(t, "Location", columns[2].Name)
}

func TestKeyVaultHandler_GetCellValue(t *testing.T) {
	handler := NewKeyVaultHandler()
	
	testValue := "test-value"
	resource := &models.Resource{
		Type:     "Microsoft.KeyVault/vaults",
		Name:     "test-keyvault",
		Location: "eastus",
		Tags:     map[string]*string{"env": &testValue},
	}
	
	tests := []struct {
		name        string
		columnIndex int
		expected    string
	}{
		{
			name:        "Column 0 - Key Vault type",
			columnIndex: 0,
			expected:    "vaults",
		},
		{
			name:        "Column 1 - Key Vault name",
			columnIndex: 1,
			expected:    "test-keyvault",
		},
		{
			name:        "Column 2 - Location",
			columnIndex: 2,
			expected:    "eastus",
		},
		{
			name:        "Invalid column index",
			columnIndex: 3,
			expected:    "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.GetCellValue(resource, tt.columnIndex)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestKeyVaultHandler_GetActions(t *testing.T) {
	handler := NewKeyVaultHandler()
	actions := handler.GetActions()
	
	assert.Len(t, actions, 2)
	assert.Equal(t, 'e', actions[0].Key)
	assert.Equal(t, "Explore Key Vault", actions[0].Label)
	assert.Equal(t, 'd', actions[1].Key)
	assert.Equal(t, "Details", actions[1].Label)
}

func TestKeyVaultHandler_CanNavigateToList(t *testing.T) {
	handler := NewKeyVaultHandler()
	assert.True(t, handler.CanNavigateToList())
}

func TestKeyVaultHandler_CanExplore(t *testing.T) {
	handler := NewKeyVaultHandler()
	assert.True(t, handler.CanExplore())
}

func TestKeyVaultHandler_RenderDetails(t *testing.T) {
	handler := NewKeyVaultHandler()
	
	tests := []struct {
		name     string
		resource *models.Resource
		subID    string
		contains []string
	}{
		{
			name: "Key Vault with tags and properties",
			resource: &models.Resource{
				ID:            "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.KeyVault/vaults/test-kv",
				Name:          "test-kv",
				Type:          "Microsoft.KeyVault/vaults",
				Location:      "eastus",
				ResourceGroup: "rg1",
				Tags: map[string]*string{
					"env":     stringPtr("production"),
					"project": stringPtr("my-project"),
				},
				Properties: map[string]interface{}{
					"vaultUri":  "https://test-kv.vault.azure.net/",
					"sku":       "standard",
					"tenantId":  "tenant-123",
				},
			},
			subID: "sub1",
			contains: []string{
				"Key Vault Details",
				"test-kv",
				"Microsoft.KeyVault/vaults",
				"eastus",
				"rg1",
				"sub1",
				"env",
				"production",
				"project",
				"my-project",
				"vaultUri",
				"sku",
				"tenantId",
			},
		},
		{
			name: "Key Vault without tags",
			resource: &models.Resource{
				ID:            "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.KeyVault/vaults/test-kv",
				Name:          "test-kv",
				Type:          "Microsoft.KeyVault/vaults",
				Location:      "westus",
				ResourceGroup: "rg1",
				Tags:          map[string]*string{},
				Properties:    map[string]interface{}{},
			},
			subID: "sub1",
			contains: []string{
				"Key Vault Details",
				"test-kv",
				"Tags:[white] None",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.RenderDetails(tt.resource, tt.subID)
			
			for _, expected := range tt.contains {
				assert.Contains(t, result, expected, "Expected to find '%s' in details", expected)
			}
		})
	}
}

func TestKeyVaultHandler_NavigateToExplore(t *testing.T) {
	handler := NewKeyVaultHandler()
	resource := &models.Resource{
		Name: "test-kv",
		Type: "Microsoft.KeyVault/vaults",
	}
	
	// Should not panic
	handler.NavigateToExplore(nil, resource)
}

func TestKeyVaultHandler_RenderDetails_Formatting(t *testing.T) {
	handler := NewKeyVaultHandler()
	resource := &models.Resource{
		ID:            "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.KeyVault/vaults/test-kv",
		Name:          "test-kv",
		Type:          "Microsoft.KeyVault/vaults",
		Location:      "eastus",
		ResourceGroup: "rg1",
		Tags:          map[string]*string{},
		Properties:    map[string]interface{}{},
	}
	
	result := handler.RenderDetails(resource, "sub1")
	
	// Check formatting
	assert.True(t, strings.HasPrefix(result, "[lightblue::b]Key Vault Details"))
	assert.Contains(t, result, "[lightblue::b]ID:[white]")
	assert.Contains(t, result, "[lightblue::b]Name:[white]")
	assert.Contains(t, result, "[lightblue::b]Type:[white]")
	assert.Contains(t, result, "[lightblue::b]Location:[white]")
	assert.Contains(t, result, "[lightblue::b]Resource Group:[white]")
	assert.Contains(t, result, "[lightblue::b]Subscription ID:[white]")
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
