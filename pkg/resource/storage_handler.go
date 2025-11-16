package resource

import (
	"fmt"
	"strings"

	"azure-control-tower/internal/models"
)

const (
	storageAccountType = "Microsoft.Storage/storageAccounts"
)

// StorageHandler provides behavior for storage account resources
type StorageHandler struct {
	resourceType string
}

// NewStorageHandler creates a new storage account handler
func NewStorageHandler() *StorageHandler {
	return &StorageHandler{
		resourceType: storageAccountType,
	}
}

// GetResourceType returns the resource type
func (h *StorageHandler) GetResourceType() string {
	return h.resourceType
}

// GetDisplayName returns a display name for the resource type
func (h *StorageHandler) GetDisplayName() string {
	return "Storage Accounts"
}

// GetColumns returns the column configuration for storage accounts
func (h *StorageHandler) GetColumns() []ColumnConfig {
	return []ColumnConfig{
		{Name: "Type", Align: AlignLeft},
		{Name: "Name", Align: AlignLeft},
		{Name: "Location", Align: AlignLeft},
	}
}

// GetCellValue extracts cell values from a storage account resource
func (h *StorageHandler) GetCellValue(resource *models.Resource, columnIndex int) string {
	switch columnIndex {
	case 0:
		// Strip provider prefix
		resourceType := resource.Type
		if idx := strings.LastIndex(resourceType, "/"); idx >= 0 && idx < len(resourceType)-1 {
			return resourceType[idx+1:]
		}
		return resourceType
	case 1:
		return resource.Name
	case 2:
		return resource.Location
	default:
		return ""
	}
}

// GetActions returns the actions available for storage accounts
func (h *StorageHandler) GetActions() []Action {
	return []Action{
		{
			Key:   'e',
			Label: "Explore Storage",
			Callback: func(resource *models.Resource, context *ActionContext) bool {
				// Action will be handled by the UI layer
				return true
			},
		},
		{
			Key:   'd',
			Label: "Details",
			Callback: func(resource *models.Resource, context *ActionContext) bool {
				// Action will be handled by the UI layer
				return true
			},
		},
	}
}

// CanNavigateToList returns true as storage accounts can be navigated to from resource types view
func (h *StorageHandler) CanNavigateToList() bool {
	return true
}

// CanExplore returns true as storage accounts support exploration
func (h *StorageHandler) CanExplore() bool {
	return true
}

// RenderDetails renders the details view for a storage account resource
func (h *StorageHandler) RenderDetails(resource *models.Resource, subscriptionID string) string {
	// Use the same format as default handler for now
	// This could be extended with storage-specific details
	var content strings.Builder
	content.WriteString("[lightblue::b]Storage Account Details[white]\n\n")
	content.WriteString("[lightblue::b]ID:[white] " + resource.ID + "\n")
	content.WriteString("[lightblue::b]Name:[white] " + resource.Name + "\n")
	content.WriteString("[lightblue::b]Type:[white] " + resource.Type + "\n")
	content.WriteString("[lightblue::b]Location:[white] " + resource.Location + "\n")
	content.WriteString("[lightblue::b]Resource Group:[white] " + resource.ResourceGroup + "\n")
	content.WriteString("[lightblue::b]Subscription ID:[white] " + subscriptionID + "\n")

	if len(resource.Tags) > 0 {
		content.WriteString("\n[lightblue::b]Tags:[white]\n")
		for key, value := range resource.Tags {
			val := ""
			if value != nil {
				val = *value
			}
			content.WriteString("  [lightblue::b]" + key + ":[white] " + val + "\n")
		}
	} else {
		content.WriteString("\n[lightblue::b]Tags:[white] None\n")
	}

	if len(resource.Properties) > 0 {
		content.WriteString("\n[lightblue::b]Properties:[white]\n")
		for key, value := range resource.Properties {
			content.WriteString("  [lightblue::b]" + key + ":[white] " + strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%v", value), "map[", ""), "]", "")) + "\n")
		}
	}

	return content.String()
}

// NavigateToExplore navigates to the storage explorer view
func (h *StorageHandler) NavigateToExplore(app interface{}, resource *models.Resource) {
	// This will be called by the UI layer to navigate to storage explorer
	// The actual navigation logic is in app.go
}

