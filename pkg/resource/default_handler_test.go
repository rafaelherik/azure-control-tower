package resource

import (
	"azure-control-tower/internal/models"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultHandler(t *testing.T) {
	handler := NewDefaultHandler()

	assert.NotNil(t, handler)
	assert.Empty(t, handler.GetResourceType())
	assert.Equal(t, "Resources", handler.GetDisplayName())
}

func TestDefaultHandler_GetColumns(t *testing.T) {
	handler := NewDefaultHandler()

	columns := handler.GetColumns()

	assert.Len(t, columns, 3)
	assert.Equal(t, "Type", columns[0].Name)
	assert.Equal(t, "Name", columns[1].Name)
	assert.Equal(t, "Location", columns[2].Name)
	assert.Equal(t, AlignLeft, columns[0].Align)
	assert.Equal(t, AlignLeft, columns[1].Align)
	assert.Equal(t, AlignLeft, columns[2].Align)
}

func TestDefaultHandler_GetCellValue(t *testing.T) {
	handler := NewDefaultHandler()

	tests := []struct {
		name        string
		resource    *models.Resource
		columnIndex int
		expected    string
	}{
		{
			name: "Column 0 - Simple type",
			resource: &models.Resource{
				Type:     "Microsoft.Storage/storageAccounts",
				Name:     "mystorageaccount",
				Location: "eastus",
			},
			columnIndex: 0,
			expected:    "storageAccounts",
		},
		{
			name: "Column 0 - Type without slash",
			resource: &models.Resource{
				Type:     "GenericResource",
				Name:     "myresource",
				Location: "westus",
			},
			columnIndex: 0,
			expected:    "GenericResource",
		},
		{
			name: "Column 0 - Nested type",
			resource: &models.Resource{
				Type:     "Microsoft.Compute/virtualMachines/extensions",
				Name:     "myextension",
				Location: "centralus",
			},
			columnIndex: 0,
			expected:    "extensions",
		},
		{
			name: "Column 0 - Type ending with slash",
			resource: &models.Resource{
				Type:     "Microsoft.Network/virtualNetworks/",
				Name:     "myvnet",
				Location: "northeurope",
			},
			columnIndex: 0,
			expected:    "Microsoft.Network/virtualNetworks/",
		},
		{
			name: "Column 1 - Resource name",
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
				Location: "eastus",
			},
			columnIndex: 2,
			expected:    "eastus",
		},
		{
			name: "Invalid column index",
			resource: &models.Resource{
				Type:     "Microsoft.Storage/storageAccounts",
				Name:     "mystorageaccount",
				Location: "eastus",
			},
			columnIndex: 99,
			expected:    "",
		},
		{
			name: "Negative column index",
			resource: &models.Resource{
				Type:     "Microsoft.Storage/storageAccounts",
				Name:     "mystorageaccount",
				Location: "eastus",
			},
			columnIndex: -1,
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

func TestDefaultHandler_GetActions(t *testing.T) {
	handler := NewDefaultHandler()

	actions := handler.GetActions()

	assert.Len(t, actions, 1)
	assert.Equal(t, 'd', actions[0].Key)
	assert.Equal(t, "Details", actions[0].Label)
	assert.NotNil(t, actions[0].Callback)
}

func TestDefaultHandler_CanNavigateToList(t *testing.T) {
	handler := NewDefaultHandler()

	assert.True(t, handler.CanNavigateToList())
}

func TestDefaultHandler_CanExplore(t *testing.T) {
	handler := NewDefaultHandler()

	assert.False(t, handler.CanExplore())
}

func TestDefaultHandler_RenderDetails(t *testing.T) {
	handler := NewDefaultHandler()

	t.Run("Resource with tags and properties", func(t *testing.T) {
		tag1 := "production"
		tag2 := "us-east"
		resource := &models.Resource{
			ID:            "/subscriptions/sub-123/resourceGroups/my-rg/providers/Microsoft.Storage/storageAccounts/myaccount",
			Name:          "myaccount",
			Type:          "Microsoft.Storage/storageAccounts",
			Location:      "eastus",
			ResourceGroup: "my-rg",
			Tags: map[string]*string{
				"environment": &tag1,
				"region":      &tag2,
			},
			Properties: map[string]interface{}{
				"provisioningState": "Succeeded",
				"primaryLocation":   "eastus",
			},
		}

		result := handler.RenderDetails(resource, "sub-123")

		assert.Contains(t, result, "Resource Details")
		assert.Contains(t, result, "ID:")
		assert.Contains(t, result, resource.ID)
		assert.Contains(t, result, "Name:")
		assert.Contains(t, result, "myaccount")
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
		assert.Contains(t, result, "region:")
		assert.Contains(t, result, "us-east")
		assert.Contains(t, result, "Properties:")
		assert.Contains(t, result, "provisioningState:")
		assert.Contains(t, result, "Succeeded")
	})

	t.Run("Resource without tags", func(t *testing.T) {
		resource := &models.Resource{
			ID:            "/subscriptions/sub-123/resourceGroups/my-rg/providers/Microsoft.Storage/storageAccounts/myaccount",
			Name:          "myaccount",
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

	t.Run("Resource with nil tag value", func(t *testing.T) {
		resource := &models.Resource{
			ID:            "/subscriptions/sub-123/resourceGroups/my-rg/providers/Microsoft.Storage/storageAccounts/myaccount",
			Name:          "myaccount",
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
		// Should handle nil value gracefully
		assert.NotContains(t, result, "<nil>")
	})
}

func TestDefaultHandler_NavigateToExplore(t *testing.T) {
	handler := NewDefaultHandler()
	resource := &models.Resource{
		Name: "test-resource",
		Type: "Microsoft.Storage/storageAccounts",
	}

	// Should not panic
	assert.NotPanics(t, func() {
		handler.NavigateToExplore(nil, resource)
	})
}

func TestAlignmentConstants(t *testing.T) {
	assert.Equal(t, 0, AlignLeft)
	assert.Equal(t, 1, AlignCenter)
	assert.Equal(t, 2, AlignRight)
}

func TestDefaultHandler_GetCellValue_StripProvider(t *testing.T) {
	handler := NewDefaultHandler()

	tests := []struct {
		resourceType string
		expected     string
	}{
		{"Microsoft.Storage/storageAccounts", "storageAccounts"},
		{"Microsoft.Compute/virtualMachines", "virtualMachines"},
		{"Microsoft.Network/virtualNetworks", "virtualNetworks"},
		{"Microsoft.Web/sites", "sites"},
		{"Microsoft.Sql/servers/databases", "databases"},
		{"NoSlash", "NoSlash"},
		{"", ""},
		{"/", "/"},
		{"Provider/", "Provider/"},
	}

	for _, tt := range tests {
		t.Run(tt.resourceType, func(t *testing.T) {
			resource := &models.Resource{Type: tt.resourceType}
			result := handler.GetCellValue(resource, 0)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultHandler_RenderDetails_Formatting(t *testing.T) {
	handler := NewDefaultHandler()
	resource := &models.Resource{
		ID:            "/subscriptions/sub-123/resourceGroups/my-rg/providers/Microsoft.Storage/storageAccounts/myaccount",
		Name:          "myaccount",
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
}

func BenchmarkDefaultHandler_GetCellValue(b *testing.B) {
	handler := NewDefaultHandler()
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

func BenchmarkDefaultHandler_RenderDetails(b *testing.B) {
	handler := NewDefaultHandler()
	tag1 := "production"
	resource := &models.Resource{
		ID:            "/subscriptions/sub-123/resourceGroups/my-rg/providers/Microsoft.Storage/storageAccounts/myaccount",
		Name:          "myaccount",
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
