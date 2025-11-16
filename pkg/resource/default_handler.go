package resource

import (
	"fmt"
	"strings"

	"azure-control-tower/internal/models"
)

const (
	AlignLeft   = 0
	AlignCenter = 1
	AlignRight  = 2
)

// DefaultHandler provides default behavior for generic resources
type DefaultHandler struct {
	resourceType string
}

// NewDefaultHandler creates a new default handler
func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{
		resourceType: "", // Empty string indicates default handler
	}
}

// GetResourceType returns the resource type (empty for default)
func (h *DefaultHandler) GetResourceType() string {
	return h.resourceType
}

// GetDisplayName returns a display name for the resource type
func (h *DefaultHandler) GetDisplayName() string {
	return "Resources"
}

// GetColumns returns the default column configuration
func (h *DefaultHandler) GetColumns() []ColumnConfig {
	return []ColumnConfig{
		{Name: "Type", Align: AlignLeft},
		{Name: "Name", Align: AlignLeft},
		{Name: "Location", Align: AlignLeft},
	}
}

// GetCellValue extracts cell values from a resource
func (h *DefaultHandler) GetCellValue(resource *models.Resource, columnIndex int) string {
	switch columnIndex {
	case 0:
		// Strip provider prefix (e.g., "Microsoft.Storage/storageAccounts" -> "storageAccounts")
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

// GetActions returns the default actions available for resources
func (h *DefaultHandler) GetActions() []Action {
	return []Action{
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

// CanNavigateToList returns true if this resource type can be navigated to from resource types view
func (h *DefaultHandler) CanNavigateToList() bool {
	return true // Default handler allows navigation to list
}

// CanExplore returns false as default handler doesn't support exploration
func (h *DefaultHandler) CanExplore() bool {
	return false
}

// RenderDetails renders the details view for a resource
func (h *DefaultHandler) RenderDetails(resource *models.Resource, subscriptionID string) string {
	var content strings.Builder
	content.WriteString("[lightblue::b]Resource Details[white]\n\n")
	content.WriteString(fmt.Sprintf("[lightblue::b]ID:[white] %s\n", resource.ID))
	content.WriteString(fmt.Sprintf("[lightblue::b]Name:[white] %s\n", resource.Name))
	content.WriteString(fmt.Sprintf("[lightblue::b]Type:[white] %s\n", resource.Type))
	content.WriteString(fmt.Sprintf("[lightblue::b]Location:[white] %s\n", resource.Location))
	content.WriteString(fmt.Sprintf("[lightblue::b]Resource Group:[white] %s\n", resource.ResourceGroup))
	content.WriteString(fmt.Sprintf("[lightblue::b]Subscription ID:[white] %s\n", subscriptionID))

	if len(resource.Tags) > 0 {
		content.WriteString("\n[lightblue::b]Tags:[white]\n")
		for key, value := range resource.Tags {
			val := ""
			if value != nil {
				val = *value
			}
			content.WriteString(fmt.Sprintf("  [lightblue::b]%s:[white] %s\n", key, val))
		}
	} else {
		content.WriteString("\n[lightblue::b]Tags:[white] None\n")
	}

	if len(resource.Properties) > 0 {
		content.WriteString("\n[lightblue::b]Properties:[white]\n")
		for key, value := range resource.Properties {
			content.WriteString(fmt.Sprintf("  [lightblue::b]%s:[white] %v\n", key, value))
		}
	}

	return content.String()
}

// NavigateToExplore is a no-op for default handler
func (h *DefaultHandler) NavigateToExplore(app interface{}, resource *models.Resource) {
	// Default handler doesn't support exploration
}
