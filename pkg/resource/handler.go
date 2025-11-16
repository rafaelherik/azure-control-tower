package resource

import (
	"azure-control-tower/internal/models"
)

// ColumnConfig defines a table column configuration
type ColumnConfig struct {
	Name       string
	Width      int // 0 means auto-width
	Align      int // 0=left, 1=center, 2=right
	Selectable bool
}

// ResourceHandler defines the interface for resource-type-specific behavior
type ResourceHandler interface {
	// Metadata
	GetResourceType() string
	GetDisplayName() string // e.g., "Storage Accounts"

	// List view configuration
	GetColumns() []ColumnConfig
	GetCellValue(resource *models.Resource, columnIndex int) string

	// Actions
	GetActions() []Action
	CanNavigateToList() bool // Can navigate to resource type list from resource types view
	CanExplore() bool         // Has special exploration view (like storage explorer)

	// Details
	RenderDetails(resource *models.Resource, subscriptionID string) string

	// Navigation
	NavigateToExplore(app interface{}, resource *models.Resource) // For special views (using interface{} to avoid circular dependency)
}

