package resource

import (
	"azure-control-tower/internal/models"
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock handler for testing
type mockHandler struct {
	resourceType       string
	canNavigateToList  bool
	canExplore         bool
}

func (m *mockHandler) GetResourceType() string {
	return m.resourceType
}

func (m *mockHandler) GetDisplayName() string {
	return "Mock Handler"
}

func (m *mockHandler) GetColumns() []ColumnConfig {
	return []ColumnConfig{}
}

func (m *mockHandler) GetCellValue(resource *models.Resource, columnIndex int) string {
	return ""
}

func (m *mockHandler) GetActions() []Action {
	return []Action{}
}

func (m *mockHandler) CanNavigateToList() bool {
	return m.canNavigateToList
}

func (m *mockHandler) CanExplore() bool {
	return m.canExplore
}

func (m *mockHandler) RenderDetails(resource *models.Resource, subscriptionID string) string {
	return ""
}

func (m *mockHandler) NavigateToExplore(app interface{}, resource *models.Resource) {
}

// Mock fetcher for testing
type mockFetcher struct {
	name string
}

func (m *mockFetcher) FetchResources(ctx context.Context, subscriptionID string, resourceType string) ([]interface{}, error) {
	return nil, nil
}

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()

	assert.NotNil(t, registry)
	assert.NotNil(t, registry.fetchers)
	assert.NotNil(t, registry.handlers)
	assert.Equal(t, 0, len(registry.fetchers))
	assert.Equal(t, 0, len(registry.handlers))
}

func TestRegisterAndGetFetcher(t *testing.T) {
	registry := NewRegistry()
	fetcher := &mockFetcher{name: "test-fetcher"}

	// Register fetcher
	registry.Register("Microsoft.Storage/storageAccounts", fetcher)

	// Get fetcher
	retrieved, err := registry.GetFetcher("Microsoft.Storage/storageAccounts")
	require.NoError(t, err)
	assert.Equal(t, fetcher, retrieved)
}

func TestGetFetcher_NotFound(t *testing.T) {
	registry := NewRegistry()

	fetcher, err := registry.GetFetcher("non-existent-type")

	assert.Error(t, err)
	assert.Nil(t, fetcher)
	assert.Contains(t, err.Error(), "not registered")
}

func TestListResourceTypes(t *testing.T) {
	registry := NewRegistry()

	// Initially empty
	types := registry.ListResourceTypes()
	assert.Empty(t, types)

	// Register some fetchers
	registry.Register("Microsoft.Storage/storageAccounts", &mockFetcher{})
	registry.Register("Microsoft.Compute/virtualMachines", &mockFetcher{})
	registry.Register("Microsoft.Network/virtualNetworks", &mockFetcher{})

	// List all types
	types = registry.ListResourceTypes()
	assert.Len(t, types, 3)
	assert.Contains(t, types, "Microsoft.Storage/storageAccounts")
	assert.Contains(t, types, "Microsoft.Compute/virtualMachines")
	assert.Contains(t, types, "Microsoft.Network/virtualNetworks")
}

func TestRegisterHandler(t *testing.T) {
	registry := NewRegistry()
	handler := &mockHandler{
		resourceType:      "Microsoft.Storage/storageAccounts",
		canNavigateToList: true,
	}

	registry.RegisterHandler(handler)

	retrieved, err := registry.GetHandler("Microsoft.Storage/storageAccounts")
	require.NoError(t, err)
	assert.Equal(t, handler, retrieved)
}

func TestGetHandler_NotFound(t *testing.T) {
	registry := NewRegistry()

	handler, err := registry.GetHandler("non-existent-type")

	assert.Error(t, err)
	assert.Nil(t, handler)
	assert.Contains(t, err.Error(), "handler not registered")
}

func TestGetHandlerOrDefault_WithHandler(t *testing.T) {
	registry := NewRegistry()
	handler := &mockHandler{
		resourceType:      "Microsoft.Storage/storageAccounts",
		canNavigateToList: true,
	}

	registry.RegisterHandler(handler)

	retrieved := registry.GetHandlerOrDefault("Microsoft.Storage/storageAccounts")
	assert.Equal(t, handler, retrieved)
}

func TestGetHandlerOrDefault_WithDefaultHandler(t *testing.T) {
	registry := NewRegistry()
	defaultHandler := &mockHandler{
		resourceType:      "",
		canNavigateToList: false,
	}

	// Register default handler (empty string key)
	registry.RegisterHandler(defaultHandler)

	// Request non-existent handler should return default
	retrieved := registry.GetHandlerOrDefault("non-existent-type")
	assert.Equal(t, defaultHandler, retrieved)
}

func TestGetHandlerOrDefault_NoHandlerNoDefault(t *testing.T) {
	registry := NewRegistry()

	retrieved := registry.GetHandlerOrDefault("non-existent-type")
	assert.Nil(t, retrieved)
}

func TestGetSupportedResourceTypes(t *testing.T) {
	registry := NewRegistry()

	// Initially empty
	types := registry.GetSupportedResourceTypes()
	assert.Empty(t, types)

	// Register handlers with different CanNavigateToList values
	handler1 := &mockHandler{
		resourceType:      "Microsoft.Storage/storageAccounts",
		canNavigateToList: true,
	}
	handler2 := &mockHandler{
		resourceType:      "Microsoft.Compute/virtualMachines",
		canNavigateToList: false,
	}
	handler3 := &mockHandler{
		resourceType:      "Microsoft.Network/virtualNetworks",
		canNavigateToList: true,
	}

	registry.RegisterHandler(handler1)
	registry.RegisterHandler(handler2)
	registry.RegisterHandler(handler3)

	// Only handlers with CanNavigateToList = true should be returned
	types = registry.GetSupportedResourceTypes()
	assert.Len(t, types, 2)
	assert.Contains(t, types, "Microsoft.Storage/storageAccounts")
	assert.Contains(t, types, "Microsoft.Network/virtualNetworks")
	assert.NotContains(t, types, "Microsoft.Compute/virtualMachines")
}

func TestRegistryConcurrentAccess(t *testing.T) {
	registry := NewRegistry()
	var wg sync.WaitGroup

	// Concurrent writes
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			handler := &mockHandler{
				resourceType:      "type-" + string(rune(index)),
				canNavigateToList: true,
			}
			registry.RegisterHandler(handler)
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = registry.GetSupportedResourceTypes()
		}()
	}

	wg.Wait()

	// Verify registry state is consistent
	types := registry.GetSupportedResourceTypes()
	assert.NotEmpty(t, types)
}

func TestRegistryHandlerReplacement(t *testing.T) {
	registry := NewRegistry()

	// Register first handler
	handler1 := &mockHandler{
		resourceType:      "Microsoft.Storage/storageAccounts",
		canNavigateToList: true,
	}
	registry.RegisterHandler(handler1)

	// Register second handler with same type (should replace)
	handler2 := &mockHandler{
		resourceType:      "Microsoft.Storage/storageAccounts",
		canNavigateToList: false,
	}
	registry.RegisterHandler(handler2)

	// Should get the second handler
	retrieved, err := registry.GetHandler("Microsoft.Storage/storageAccounts")
	require.NoError(t, err)
	assert.Equal(t, handler2, retrieved)
	assert.False(t, retrieved.CanNavigateToList())
}
