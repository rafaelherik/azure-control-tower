package navigation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewState(t *testing.T) {
	state := NewState()

	assert.NotNil(t, state)
	assert.Equal(t, ViewSubscriptions, state.CurrentView)
	assert.Empty(t, state.SelectedSubscriptionID)
	assert.Empty(t, state.SelectedSubscriptionName)
	assert.Empty(t, state.SelectedResourceGroupName)
	assert.Empty(t, state.SelectedResourceType)
	assert.False(t, state.InDetailsView)
}

func TestNavigateToSubscriptions(t *testing.T) {
	state := &State{
		CurrentView:               ViewResourceGroups,
		SelectedSubscriptionID:    "sub-123",
		SelectedSubscriptionName:  "Test Subscription",
		SelectedResourceGroupName: "test-rg",
		SelectedResourceType:      "Microsoft.Storage/storageAccounts",
		InDetailsView:             true,
	}

	state.NavigateToSubscriptions()

	assert.Equal(t, ViewSubscriptions, state.CurrentView)
	assert.Empty(t, state.SelectedSubscriptionID)
	assert.Empty(t, state.SelectedSubscriptionName)
	assert.Empty(t, state.SelectedResourceGroupName)
	assert.Empty(t, state.SelectedResourceType)
	assert.False(t, state.InDetailsView)
}

func TestNavigateToResourceGroups(t *testing.T) {
	state := NewState()

	state.NavigateToResourceGroups("sub-123", "Test Subscription")

	assert.Equal(t, ViewResourceGroups, state.CurrentView)
	assert.Equal(t, "sub-123", state.SelectedSubscriptionID)
	assert.Equal(t, "Test Subscription", state.SelectedSubscriptionName)
	assert.Empty(t, state.SelectedResourceGroupName)
	assert.Empty(t, state.SelectedResourceType)
	assert.False(t, state.InDetailsView)
}

func TestNavigateToResourceTypes(t *testing.T) {
	state := &State{
		CurrentView:               ViewResourceGroups,
		SelectedSubscriptionID:    "sub-123",
		SelectedSubscriptionName:  "Test Subscription",
		SelectedResourceType:      "old-type",
	}

	state.NavigateToResourceTypes("test-rg")

	assert.Equal(t, ViewResourceTypes, state.CurrentView)
	assert.Equal(t, "test-rg", state.SelectedResourceGroupName)
	assert.Empty(t, state.SelectedResourceType)
	assert.False(t, state.InDetailsView)
	// Subscription info should be preserved
	assert.Equal(t, "sub-123", state.SelectedSubscriptionID)
	assert.Equal(t, "Test Subscription", state.SelectedSubscriptionName)
}

func TestNavigateToResources(t *testing.T) {
	state := NewState()

	state.NavigateToResources("test-rg")

	assert.Equal(t, ViewResources, state.CurrentView)
	assert.Equal(t, "test-rg", state.SelectedResourceGroupName)
	assert.Empty(t, state.SelectedResourceType)
	assert.False(t, state.InDetailsView)
}

func TestNavigateToResourceType(t *testing.T) {
	state := &State{
		CurrentView:               ViewResourceTypes,
		SelectedSubscriptionID:    "sub-123",
		SelectedSubscriptionName:  "Test Subscription",
		SelectedResourceGroupName: "test-rg",
	}

	state.NavigateToResourceType("Microsoft.Storage/storageAccounts")

	assert.Equal(t, ViewResourceType, state.CurrentView)
	assert.Equal(t, "Microsoft.Storage/storageAccounts", state.SelectedResourceType)
	assert.False(t, state.InDetailsView)
	// Context should be preserved
	assert.Equal(t, "sub-123", state.SelectedSubscriptionID)
	assert.Equal(t, "Test Subscription", state.SelectedSubscriptionName)
	assert.Equal(t, "test-rg", state.SelectedResourceGroupName)
}

func TestNavigateToDetails(t *testing.T) {
	state := &State{
		CurrentView:   ViewResourceGroups,
		InDetailsView: false,
	}

	state.NavigateToDetails()

	assert.True(t, state.InDetailsView)
	// Current view should not change
	assert.Equal(t, ViewResourceGroups, state.CurrentView)
}

func TestNavigateBackFromDetails(t *testing.T) {
	state := &State{
		CurrentView:   ViewResourceGroups,
		InDetailsView: true,
	}

	state.NavigateBackFromDetails()

	assert.False(t, state.InDetailsView)
	assert.Equal(t, ViewResourceGroups, state.CurrentView)
}

func TestNavigateToStorageExplorer(t *testing.T) {
	state := &State{
		CurrentView:            ViewResourceType,
		SelectedStorageAccount: "old-account",
		SelectedContainer:      "old-container",
		SelectedBlob:           "old-blob",
		InDetailsView:          true,
	}

	state.NavigateToStorageExplorer("new-storage-account")

	assert.Equal(t, ViewStorageExplorer, state.CurrentView)
	assert.Equal(t, "new-storage-account", state.SelectedStorageAccount)
	assert.Empty(t, state.SelectedContainer)
	assert.Empty(t, state.SelectedBlob)
	assert.False(t, state.InDetailsView)
}

func TestNavigateToBlobs(t *testing.T) {
	state := &State{
		CurrentView:            ViewStorageExplorer,
		SelectedStorageAccount: "test-account",
		SelectedContainer:      "old-container",
		SelectedBlob:           "old-blob",
		BlobPathPrefix:         "old/path/",
	}

	state.NavigateToBlobs("new-container")

	assert.Equal(t, ViewBlobs, state.CurrentView)
	assert.Equal(t, "new-container", state.SelectedContainer)
	assert.Empty(t, state.SelectedBlob)
	assert.Empty(t, state.BlobPathPrefix)
	assert.False(t, state.InDetailsView)
	// Storage account should be preserved
	assert.Equal(t, "test-account", state.SelectedStorageAccount)
}

func TestNavigateBackFromBlobs(t *testing.T) {
	state := &State{
		CurrentView:            ViewBlobs,
		SelectedStorageAccount: "test-account",
		SelectedContainer:      "test-container",
		SelectedBlob:           "test-blob",
		BlobPathPrefix:         "folder/subfolder/",
	}

	state.NavigateBackFromBlobs()

	assert.Equal(t, ViewStorageExplorer, state.CurrentView)
	assert.Empty(t, state.SelectedContainer)
	assert.Empty(t, state.SelectedBlob)
	assert.Empty(t, state.BlobPathPrefix)
	// Storage account should be preserved
	assert.Equal(t, "test-account", state.SelectedStorageAccount)
}

func TestNavigateIntoBlobFolder(t *testing.T) {
	state := &State{
		CurrentView:    ViewBlobs,
		BlobPathPrefix: "folder/",
		SelectedBlob:   "old-blob",
	}

	state.NavigateIntoBlobFolder("folder/subfolder/")

	assert.Equal(t, "folder/subfolder/", state.BlobPathPrefix)
	assert.Empty(t, state.SelectedBlob)
}

func TestNavigateBackFromBlobFolder_WithParent(t *testing.T) {
	state := &State{
		CurrentView:    ViewBlobs,
		BlobPathPrefix: "folder/subfolder/deepfolder/",
	}

	state.NavigateBackFromBlobFolder()

	assert.Equal(t, "folder/subfolder/", state.BlobPathPrefix)
	assert.Empty(t, state.SelectedBlob)
}

func TestNavigateBackFromBlobFolder_ToRoot(t *testing.T) {
	state := &State{
		CurrentView:    ViewBlobs,
		BlobPathPrefix: "folder/",
	}

	state.NavigateBackFromBlobFolder()

	assert.Empty(t, state.BlobPathPrefix)
	assert.Equal(t, ViewBlobs, state.CurrentView)
}

func TestNavigateBackFromBlobFolder_FromRoot(t *testing.T) {
	state := &State{
		CurrentView:            ViewBlobs,
		SelectedStorageAccount: "test-account",
		SelectedContainer:      "test-container",
		BlobPathPrefix:         "",
	}

	state.NavigateBackFromBlobFolder()

	// Should navigate back to storage explorer
	assert.Equal(t, ViewStorageExplorer, state.CurrentView)
	assert.Empty(t, state.SelectedContainer)
	assert.Empty(t, state.BlobPathPrefix)
}

func TestNavigateBackFromBlobFolder_ComplexPath(t *testing.T) {
	tests := []struct {
		name           string
		initialPath    string
		expectedPath   string
		expectedView   ViewType
	}{
		{
			name:         "Three levels deep",
			initialPath:  "a/b/c/",
			expectedPath: "a/b/",
			expectedView: ViewBlobs,
		},
		{
			name:         "Two levels deep",
			initialPath:  "a/b/",
			expectedPath: "a/",
			expectedView: ViewBlobs,
		},
		{
			name:         "One level deep",
			initialPath:  "a/",
			expectedPath: "",
			expectedView: ViewBlobs,
		},
		{
			name:         "Root level",
			initialPath:  "",
			expectedPath: "",
			expectedView: ViewStorageExplorer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &State{
				CurrentView:            ViewBlobs,
				SelectedStorageAccount: "test-account",
				SelectedContainer:      "test-container",
				BlobPathPrefix:         tt.initialPath,
			}

			state.NavigateBackFromBlobFolder()

			assert.Equal(t, tt.expectedView, state.CurrentView)
			assert.Equal(t, tt.expectedPath, state.BlobPathPrefix)
		})
	}
}

func TestNavigateToMenu(t *testing.T) {
	state := &State{
		CurrentView:   ViewResourceGroups,
		InDetailsView: true,
	}

	state.NavigateToMenu()

	assert.Equal(t, ViewMenu, state.CurrentView)
	assert.False(t, state.InDetailsView)
}

func TestNavigateBackFromMenu(t *testing.T) {
	state := &State{
		CurrentView:               ViewMenu,
		SelectedSubscriptionID:    "sub-123",
		SelectedSubscriptionName:  "Test Subscription",
		SelectedResourceGroupName: "test-rg",
	}

	state.NavigateBackFromMenu()

	// Currently goes back to subscriptions (as per implementation)
	assert.Equal(t, ViewSubscriptions, state.CurrentView)
	assert.Empty(t, state.SelectedSubscriptionID)
	assert.Empty(t, state.SelectedSubscriptionName)
	assert.Empty(t, state.SelectedResourceGroupName)
}

func TestViewTypeValues(t *testing.T) {
	// Test that view type constants are unique
	views := map[ViewType]string{
		ViewSubscriptions:   "ViewSubscriptions",
		ViewResourceGroups:  "ViewResourceGroups",
		ViewResourceTypes:   "ViewResourceTypes",
		ViewResources:       "ViewResources",
		ViewResourceType:    "ViewResourceType",
		ViewDetails:         "ViewDetails",
		ViewStorageExplorer: "ViewStorageExplorer",
		ViewBlobs:           "ViewBlobs",
		ViewMenu:            "ViewMenu",
	}

	assert.Len(t, views, 9, "All view types should be unique")
}

func TestNavigationFlow_FullJourney(t *testing.T) {
	state := NewState()

	// Start at subscriptions
	assert.Equal(t, ViewSubscriptions, state.CurrentView)

	// Navigate to resource groups
	state.NavigateToResourceGroups("sub-123", "Test Subscription")
	assert.Equal(t, ViewResourceGroups, state.CurrentView)
	assert.Equal(t, "sub-123", state.SelectedSubscriptionID)

	// Navigate to resource types
	state.NavigateToResourceTypes("test-rg")
	assert.Equal(t, ViewResourceTypes, state.CurrentView)
	assert.Equal(t, "test-rg", state.SelectedResourceGroupName)

	// Navigate to specific resource type
	state.NavigateToResourceType("Microsoft.Storage/storageAccounts")
	assert.Equal(t, ViewResourceType, state.CurrentView)
	assert.Equal(t, "Microsoft.Storage/storageAccounts", state.SelectedResourceType)

	// Navigate to storage explorer
	state.NavigateToStorageExplorer("mystorageaccount")
	assert.Equal(t, ViewStorageExplorer, state.CurrentView)
	assert.Equal(t, "mystorageaccount", state.SelectedStorageAccount)

	// Navigate to blobs
	state.NavigateToBlobs("mycontainer")
	assert.Equal(t, ViewBlobs, state.CurrentView)
	assert.Equal(t, "mycontainer", state.SelectedContainer)

	// Navigate into folder
	state.NavigateIntoBlobFolder("folder1/")
	assert.Equal(t, "folder1/", state.BlobPathPrefix)

	// Navigate into subfolder
	state.NavigateIntoBlobFolder("folder1/subfolder/")
	assert.Equal(t, "folder1/subfolder/", state.BlobPathPrefix)

	// Navigate back to parent folder
	state.NavigateBackFromBlobFolder()
	assert.Equal(t, "folder1/", state.BlobPathPrefix)

	// Navigate back to root
	state.NavigateBackFromBlobFolder()
	assert.Empty(t, state.BlobPathPrefix)

	// Navigate back to storage explorer
	state.NavigateBackFromBlobFolder()
	assert.Equal(t, ViewStorageExplorer, state.CurrentView)

	// Navigate back to blobs
	state.NavigateBackFromBlobs()
	assert.Equal(t, ViewStorageExplorer, state.CurrentView)
}
