package azure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitResourceID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected []string
	}{
		{
			name:     "Valid Azure resource ID",
			id:       "/subscriptions/sub-123/resourceGroups/my-rg/providers/Microsoft.Storage/storageAccounts/myaccount",
			expected: []string{"subscriptions", "sub-123", "resourceGroups", "my-rg", "providers", "Microsoft.Storage", "storageAccounts", "myaccount"},
		},
		{
			name:     "Resource ID with nested types",
			id:       "/subscriptions/sub-123/resourceGroups/my-rg/providers/Microsoft.Compute/virtualMachines/myvm/extensions/ext1",
			expected: []string{"subscriptions", "sub-123", "resourceGroups", "my-rg", "providers", "Microsoft.Compute", "virtualMachines", "myvm", "extensions", "ext1"},
		},
		{
			name:     "Empty string",
			id:       "",
			expected: []string{},
		},
		{
			name:     "No leading slash",
			id:       "subscriptions/sub-123",
			expected: []string{},
		},
		{
			name:     "Just a slash",
			id:       "/",
			expected: []string{},
		},
		{
			name:     "Single segment",
			id:       "/subscriptions",
			expected: []string{"subscriptions"},
		},
		{
			name:     "Multiple consecutive slashes",
			id:       "/subscriptions//sub-123",
			expected: []string{"subscriptions", "sub-123"},
		},
		{
			name:     "Trailing slash",
			id:       "/subscriptions/sub-123/",
			expected: []string{"subscriptions", "sub-123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitResourceID(tt.id)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractResourceGroupFromID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected string
	}{
		{
			name:     "Valid storage account ID",
			id:       "/subscriptions/sub-123/resourceGroups/my-rg/providers/Microsoft.Storage/storageAccounts/myaccount",
			expected: "my-rg",
		},
		{
			name:     "Valid VM ID",
			id:       "/subscriptions/sub-456/resourceGroups/production-rg/providers/Microsoft.Compute/virtualMachines/myvm",
			expected: "production-rg",
		},
		{
			name:     "Resource group with special characters",
			id:       "/subscriptions/sub-789/resourceGroups/my-rg-2024/providers/Microsoft.Network/virtualNetworks/vnet1",
			expected: "my-rg-2024",
		},
		{
			name:     "Empty ID",
			id:       "",
			expected: "",
		},
		{
			name:     "Malformed ID - no resourceGroups segment",
			id:       "/subscriptions/sub-123/providers/Microsoft.Storage/storageAccounts/myaccount",
			expected: "",
		},
		{
			name:     "Malformed ID - wrong order",
			id:       "/resourceGroups/my-rg/subscriptions/sub-123",
			expected: "",
		},
		{
			name:     "Malformed ID - too short",
			id:       "/subscriptions/sub-123/resourceGroups",
			expected: "",
		},
		{
			name:     "Nested resource ID",
			id:       "/subscriptions/sub-123/resourceGroups/my-rg/providers/Microsoft.Compute/virtualMachines/myvm/extensions/ext1",
			expected: "my-rg",
		},
		{
			name:     "Resource group ID itself",
			id:       "/subscriptions/sub-123/resourceGroups/my-rg",
			expected: "my-rg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractResourceGroupFromID(tt.id)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractResourceGroupFromID_EdgeCases(t *testing.T) {
	t.Run("Resource ID with very long resource group name", func(t *testing.T) {
		id := "/subscriptions/sub-123/resourceGroups/this-is-a-very-long-resource-group-name-that-should-still-work/providers/Microsoft.Storage/storageAccounts/myaccount"
		result := extractResourceGroupFromID(id)
		assert.Equal(t, "this-is-a-very-long-resource-group-name-that-should-still-work", result)
	})

	t.Run("Resource ID with single character resource group", func(t *testing.T) {
		id := "/subscriptions/sub-123/resourceGroups/a/providers/Microsoft.Storage/storageAccounts/myaccount"
		result := extractResourceGroupFromID(id)
		assert.Equal(t, "a", result)
	})

	t.Run("Resource ID with numeric resource group", func(t *testing.T) {
		id := "/subscriptions/sub-123/resourceGroups/123/providers/Microsoft.Storage/storageAccounts/myaccount"
		result := extractResourceGroupFromID(id)
		assert.Equal(t, "123", result)
	})
}

func TestSplitResourceID_Performance(t *testing.T) {
	// Test with a very long resource ID to ensure no performance issues
	longID := "/subscriptions/very-long-subscription-id-12345678-1234-1234-1234-123456789012/resourceGroups/very-long-resource-group-name-with-many-characters/providers/Microsoft.Storage/storageAccounts/verylongstorageaccountname/blobServices/default/containers/mycontainer"

	result := splitResourceID(longID)

	assert.NotEmpty(t, result)
	assert.Contains(t, result, "subscriptions")
	assert.Contains(t, result, "resourceGroups")
	assert.Contains(t, result, "very-long-resource-group-name-with-many-characters")
}

func BenchmarkSplitResourceID(b *testing.B) {
	id := "/subscriptions/sub-123/resourceGroups/my-rg/providers/Microsoft.Storage/storageAccounts/myaccount"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = splitResourceID(id)
	}
}

func BenchmarkExtractResourceGroupFromID(b *testing.B) {
	id := "/subscriptions/sub-123/resourceGroups/my-rg/providers/Microsoft.Storage/storageAccounts/myaccount"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = extractResourceGroupFromID(id)
	}
}
