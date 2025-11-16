package models

// Subscription represents an Azure subscription
type Subscription struct {
	ID          string
	Name        string
	State       string
	DisplayName string
	TenantID    string
}

// ResourceGroup represents an Azure resource group
type ResourceGroup struct {
	Name     string
	Location string
	Tags     map[string]*string
}

// Resource represents a generic Azure resource
type Resource struct {
	ID            string
	Name          string
	Type          string
	Location      string
	ResourceGroup string
	Tags          map[string]*string
	Properties    map[string]interface{} // Generic properties
}

// ResourceTypeSummary represents a summary of resources by type
type ResourceTypeSummary struct {
	Type  string
	Count int
}
