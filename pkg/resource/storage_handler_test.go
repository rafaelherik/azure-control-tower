package resource

import (
	"azure-control-tower/internal/models"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStorageHandler(t *testing.T) {
	handler := NewStorageHandler()

	assert.NotNil(t, handler)
	assert.Equal(t, "Microsoft.Storage/storageAccounts", handler.GetResourceType())
	assert.Equal(t, "Storage Accounts", handler.GetDisplayName())
}

func TestStorageHandler_GetColumns(t *testing.T) {
	handler := NewStorageHandler()

	columns := handler.GetColumns()

	assert.Len(t, columns, 3)
	assert.Equal(t, "Type", columns[0].Name)
	assert.Equal(t, "Name", columns[1].Name)
	assert.Equal(t, "Location", columns[2].Name)
}

func TestStorageHandler_GetCellValue(t *testing.T) {
	handler := NewStorageHandler()

	tests := []struct {
		name        string
		resource    *models.Resource
		columnIndex int
		expected    string
	}{
		{
			name: "Column 0 - Storage account type",
			resource: &models.Resource{
				Type:     "Microsoft.Storage/storageAccounts",
				Name:     "mystorageaccount",
				Location: "eastus",
			},
			columnIndex: 0,
			expected:    "storageAccounts",
		},
		{
			name: "Column 1 - Storage account name",
			resource: &models.Resource{
				Type:     "Microsoft.Storage/storageAccounts",
				Name:     "mystorageaccount",
				Location: "eastus",
			},
			columnIndex: 1,
			expected:    "mystorageaccount",
		},
		{
			name: "Column 2 - Location",
			resource: &models.Resource{
				Type:     "Microsoft.Storage/storageAccounts",
				Name:     "mystorageaccount",
				Location: "westus2",
			},
			columnIndex: 2,
			expected:    "westus2",
		},
		{
			name: "Invalid column index",
			resource: &models.Resource{
				Type:     "Microsoft.Storage/storageAccounts",
				Name:     "mystorageaccount",
				Location: "eastus",
			},
			columnIndex: 5,
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.GetCellValue(tt.resource, tt.columnIndex)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStorageHandler_GetActions(t *testing.T) {
	handler := NewStorageHandler()

	actions := handler.GetActions()

	assert.Len(t, actions, 2)

	// First action should be explore
	assert.Equal(t, 'e', actions[0].Key)
	assert.Equal(t, "Explore Storage", actions[0].Label)
	assert.NotNil(t, actions[0].Callback)

	// Second action should be details
	assert.Equal(t, 'd', actions[1].Key)
	assert.Equal(t, "Details", actions[1].Label)
	assert.NotNil(t, actions[1].Callback)
}

func TestStorageHandler_CanNavigateToList(t *testing.T) {
	handler := NewStorageHandler()

	assert.True(t, handler.CanNavigateToList())
}

func TestStorageHandler_CanExplore(t *testing.T) {
	handler := NewStorageHandler()

	assert.True(t, handler.CanExplore())
}

func TestStorageHandler_RenderDetails(t *testing.T) {
	handler := NewStorageHandler()

	t.Run("Storage account with tags and properties", func(t *testing.T) {
		tag1 := "production"
		tag2 := "critical"
		resource := &models.Resource{
			ID:            "/subscriptions/sub-123/resourceGroups/my-rg/providers/Microsoft.Storage/storageAccounts/mystorageaccount",
			Name:          "mystorageaccount",
			Type:          "Microsoft.Storage/storageAccounts",
			Location:      "eastus",
			ResourceGroup: "my-rg",
			Tags: map[string]*string{
				"environment": &tag1,
				"tier":        &tag2,
			},
			Properties: map[string]interface{}{
				"provisioningState":      "Succeeded",
				"primaryLocation":        "eastus",
				"statusOfPrimary":        "available",
				"supportsHttpsTrafficOnly": true,
			},
		}

		result := handler.RenderDetails(resource, "sub-123")

		assert.Contains(t, result, "Storage Account Details")
		assert.Contains(t, result, "ID:")
		assert.Contains(t, result, resource.ID)
		assert.Contains(t, result, "Name:")
		assert.Contains(t, result, "mystorageaccount")
		assert.Contains(t, result, "Type:")
		assert.Contains(t, result, "Microsoft.Storage/storageAccounts")
		assert.Contains(t, result, "Location:")
		assert.Contains(t, result, "eastus")
		assert.Contains(t, result, "Resource Group:")
		assert.Contains(t, result, "my-rg")
		assert.Contains(t, result, "Subscription ID:")
		assert.Contains(t, result, "sub-123")
		assert.Contains(t, result, "Tags:")
		assert.Contains(t, result, "environment:")
		assert.Contains(t, result, "production")
		assert.Contains(t, result, "Properties:")
		assert.Contains(t, result, "provisioningState:")
	})

	t.Run("Storage account without tags", func(t *testing.T) {
		resource := &models.Resource{
			ID:            "/subscriptions/sub-123/resourceGroups/my-rg/providers/Microsoft.Storage/storageAccounts/mystorageaccount",
			Name:          "mystorageaccount",
			Type:          "Microsoft.Storage/storageAccounts",
			Location:      "eastus",
			ResourceGroup: "my-rg",
			Tags:          map[string]*string{},
			Properties:    map[string]interface{}{},
		}

		result := handler.RenderDetails(resource, "sub-123")

		assert.Contains(t, result, "Tags:")
		assert.Contains(t, result, "None")
	})

	t.Run("Storage account with nil tag value", func(t *testing.T) {
		resource := &models.Resource{
			ID:            "/subscriptions/sub-123/resourceGroups/my-rg/providers/Microsoft.Storage/storageAccounts/mystorageaccount",
			Name:          "mystorageaccount",
			Type:          "Microsoft.Storage/storageAccounts",
			Location:      "eastus",
			ResourceGroup: "my-rg",
			Tags: map[string]*string{
				"empty-tag": nil,
			},
			Properties: map[string]interface{}{},
		}

		result := handler.RenderDetails(resource, "sub-123")

		assert.Contains(t, result, "empty-tag:")
		// Should handle nil value gracefully (empty string)
		lines := strings.Split(result, "\n")
		found := false
		for _, line := range lines {
			if strings.Contains(line, "empty-tag:") {
				found = true
				// After the tag name, there should just be whitespace or empty string
				parts := strings.Split(line, "empty-tag:")
				assert.Len(t, parts, 2)
			}
		}
		assert.True(t, found, "Should find empty-tag in output")
	})

	t.Run("Storage account with complex properties", func(t *testing.T) {
		resource := &models.Resource{
			ID:            "/subscriptions/sub-123/resourceGroups/my-rg/providers/Microsoft.Storage/storageAccounts/mystorageaccount",
			Name:          "mystorageaccount",
			Type:          "Microsoft.Storage/storageAccounts",
			Location:      "eastus",
			ResourceGroup: "my-rg",
			Tags:          map[string]*string{},
			Properties: map[string]interface{}{
				"encryption": map[string]interface{}{
					"services": map[string]interface{}{
						"blob": map[string]interface{}{
							"enabled": true,
						},
					},
				},
			},
		}

		result := handler.RenderDetails(resource, "sub-123")

		assert.Contains(t, result, "Properties:")
		assert.Contains(t, result, "encryption:")
		// The implementation flattens maps, so we just verify it doesn't panic
		assert.NotEmpty(t, result)
	})
}

func TestStorageHandler_NavigateToExplore(t *testing.T) {
	handler := NewStorageHandler()
	resource := &models.Resource{
		Name: "mystorageaccount",
		Type: "Microsoft.Storage/storageAccounts",
	}

	// Should not panic
	assert.NotPanics(t, func() {
		handler.NavigateToExplore(nil, resource)
	})
}

func TestStorageHandler_RenderDetails_Formatting(t *testing.T) {
	handler := NewStorageHandler()
	resource := &models.Resource{
		ID:            "/subscriptions/sub-123/resourceGroups/my-rg/providers/Microsoft.Storage/storageAccounts/mystorageaccount",
		Name:          "mystorageaccount",
		Type:          "Microsoft.Storage/storageAccounts",
		Location:      "eastus",
		ResourceGroup: "my-rg",
		Tags:          map[string]*string{},
		Properties:    map[string]interface{}{},
	}

	result := handler.RenderDetails(resource, "sub-123")

	// Check for proper formatting with tview tags
	assert.True(t, strings.Contains(result, "[lightblue::b]"))
	assert.True(t, strings.Contains(result, "[white]"))
	assert.True(t, strings.Contains(result, "\n"))
	// Specific to storage handler
	assert.Contains(t, result, "Storage Account Details")
}

func TestStorageHandler_CompareWithDefault(t *testing.T) {
	storageHandler := NewStorageHandler()
	defaultHandler := NewDefaultHandler()

	// Both should have same column structure
	assert.Equal(t, len(defaultHandler.GetColumns()), len(storageHandler.GetColumns()))

	// Storage handler should support exploration
	assert.True(t, storageHandler.CanExplore())
	assert.False(t, defaultHandler.CanExplore())

	// Storage handler should have more actions
	assert.Greater(t, len(storageHandler.GetActions()), len(defaultHandler.GetActions()))
}

func BenchmarkStorageHandler_GetCellValue(b *testing.B) {
	handler := NewStorageHandler()
	resource := &models.Resource{
		Type:     "Microsoft.Storage/storageAccounts",
		Name:     "mystorageaccount",
		Location: "eastus",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = handler.GetCellValue(resource, 0)
		_ = handler.GetCellValue(resource, 1)
		_ = handler.GetCellValue(resource, 2)
	}
}

func BenchmarkStorageHandler_RenderDetails(b *testing.B) {
	handler := NewStorageHandler()
	tag1 := "production"
	resource := &models.Resource{
		ID:            "/subscriptions/sub-123/resourceGroups/my-rg/providers/Microsoft.Storage/storageAccounts/mystorageaccount",
		Name:          "mystorageaccount",
		Type:          "Microsoft.Storage/storageAccounts",
		Location:      "eastus",
		ResourceGroup: "my-rg",
		Tags: map[string]*string{
			"environment": &tag1,
		},
		Properties: map[string]interface{}{
			"provisioningState": "Succeeded",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = handler.RenderDetails(resource, "sub-123")
	}
}
