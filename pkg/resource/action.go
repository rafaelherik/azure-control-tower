package resource

import (
	"azure-control-tower/internal/models"
)

// ActionContext provides context for action execution
type ActionContext struct {
	SubscriptionID    string
	SubscriptionName  string
	ResourceGroupName string
	App               interface{} // *ui.App - using interface{} to avoid circular dependency
}

// Action represents a resource-specific action that can be performed
type Action struct {
	Key      rune
	Label    string
	Callback func(resource *models.Resource, context *ActionContext) bool
}

