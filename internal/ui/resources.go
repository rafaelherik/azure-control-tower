package ui

import (
	"context"
	"strings"

	"azure-control-tower/internal/models"
	"azure-control-tower/pkg/resource"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ResourceRowData wraps Resource with context info for display
type ResourceRowData struct {
	Resource          *models.Resource
	SubscriptionID    string
	SubscriptionName  string
	ResourceGroupName string
}

// ResourcesView displays a table of resources in a resource group
type ResourcesView struct {
	*TableView
	resources          []*models.Resource
	subscriptionID     string
	subscriptionName   string
	resourceGroupName  string
	registry           *resource.Registry
	onShowResourceType func(resourceType string)
	onShowDetails      func(resource *models.Resource)
	onExploreStorage   func(resource *models.Resource)
}

// NewResourcesView creates a new resources view
func NewResourcesView(registry *resource.Registry) *ResourcesView {
	rv := &ResourcesView{
		registry: registry,
	}

	// Create table configuration with default values
	// Will be updated dynamically based on resource types in LoadResources
	config := &TableConfig{
		Title: "", // Will be set in LoadResources
		Columns: []ColumnConfig{
			{Name: "Type", Align: tview.AlignLeft},
			{Name: "Name", Align: tview.AlignLeft},
			{Name: "Location", Align: tview.AlignLeft},
		},
		RowActions: []RowAction{},
		OnSelect: func(rowIndex int, data interface{}) {
			// Enter key on a resource - could show details or do nothing
		},
		GetCellValue: func(data interface{}, columnIndex int) string {
			rowData, ok := data.(*ResourceRowData)
			if !ok {
				return ""
			}
			// Get handler for this resource type
			handler := rv.registry.GetHandlerOrDefault(rowData.Resource.Type)
			if handler != nil {
				return handler.GetCellValue(rowData.Resource, columnIndex)
			}
			// Fallback to default behavior
			switch columnIndex {
			case 0:
				resourceType := rowData.Resource.Type
				if idx := strings.LastIndex(resourceType, "/"); idx >= 0 && idx < len(resourceType)-1 {
					return resourceType[idx+1:]
				}
				return resourceType
			case 1:
				return rowData.Resource.Name
			case 2:
				return rowData.Resource.Location
			default:
				return ""
			}
		},
	}

	rv.TableView = NewTableView(config)
	return rv
}

// LoadResources loads resources into the view
func (rv *ResourcesView) LoadResources(ctx context.Context, resources []*models.Resource, subscriptionID, subscriptionName, resourceGroupName string) error {
	rv.resources = resources
	rv.subscriptionID = subscriptionID
	rv.subscriptionName = subscriptionName
	rv.resourceGroupName = resourceGroupName

	// Update title
	rv.SetTitle("")

	// Determine columns and actions based on resource types in the list
	// Use the first resource's handler, or default if empty
	var handler resource.ResourceHandler
	if len(resources) > 0 {
		handler = rv.registry.GetHandlerOrDefault(resources[0].Type)
	} else {
		handler = rv.registry.GetHandlerOrDefault("")
	}

	if handler != nil {
		// Update columns from handler - convert resource.ColumnConfig to ui.ColumnConfig
		resourceColumns := handler.GetColumns()
		columns := make([]ColumnConfig, len(resourceColumns))
		for i, col := range resourceColumns {
			columns[i] = ColumnConfig{
				Name:       col.Name,
				Width:      col.Width,
				Align:      col.Align,
				Selectable: col.Selectable,
			}
		}
		rv.config.Columns = columns

		// Build row actions from handler actions
		actions := handler.GetActions()
		rowActions := make([]RowAction, 0, len(actions)+1) // +1 for filter by type

		// Add filter by type action
		rowActions = append(rowActions, RowAction{
			Rune:  't',
			Label: "Filter by Type",
			Callback: func(rowIndex int, data interface{}) bool {
				if rowData, ok := data.(*ResourceRowData); ok && rv.onShowResourceType != nil {
					rv.onShowResourceType(rowData.Resource.Type)
					return true
				}
				return false
			},
		})

		// Add handler actions
		for _, action := range actions {
			actionCopy := action // Capture loop variable
			rowActions = append(rowActions, RowAction{
				Rune:  actionCopy.Key,
				Label: actionCopy.Label,
				Callback: func(rowIndex int, data interface{}) bool {
					if rowData, ok := data.(*ResourceRowData); ok {
						actionContext := &resource.ActionContext{
							SubscriptionID:    subscriptionID,
							SubscriptionName:  subscriptionName,
							ResourceGroupName: resourceGroupName,
						}
						if actionCopy.Callback != nil {
							if actionCopy.Callback(rowData.Resource, actionContext) {
								// Route to appropriate callback based on action key
								switch actionCopy.Key {
								case 'e', 'E':
									if rv.onExploreStorage != nil && handler.CanExplore() {
										rv.onExploreStorage(rowData.Resource)
									}
								case 'd', 'D':
									if rv.onShowDetails != nil {
										rv.onShowDetails(rowData.Resource)
									}
								}
								return true
							}
						}
					}
					return false
				},
			})
		}

		rv.config.RowActions = rowActions
		rv.SetConfig(rv.config)
	}

	// Convert to interface{} slice with row data
	data := make([]interface{}, len(resources))
	for i, res := range resources {
		data[i] = &ResourceRowData{
			Resource:          res,
			SubscriptionID:    subscriptionID,
			SubscriptionName:  subscriptionName,
			ResourceGroupName: resourceGroupName,
		}
	}

	rv.LoadData(data)
	return nil
}

// SetOnShowResourceType sets the callback for when resource type filter is requested (t key)
func (rv *ResourcesView) SetOnShowResourceType(callback func(string)) {
	rv.onShowResourceType = callback
}

// SetOnShowDetails sets the callback for when details are requested (d key)
func (rv *ResourcesView) SetOnShowDetails(callback func(*models.Resource)) {
	rv.onShowDetails = callback
}

// SetOnExploreStorage sets the callback for when storage exploration is requested (e key)
func (rv *ResourcesView) SetOnExploreStorage(callback func(*models.Resource)) {
	rv.onExploreStorage = callback
}

// GetSubscriptionID returns the current subscription ID
func (rv *ResourcesView) GetSubscriptionID() string {
	return rv.subscriptionID
}

// GetResourceGroupName returns the current resource group name
func (rv *ResourcesView) GetResourceGroupName() string {
	return rv.resourceGroupName
}

// HandleKey handles key events for this view
func (rv *ResourcesView) HandleKey(event *tcell.EventKey) *tcell.EventKey {
	// Handle 't' and 'T' for resource type filter
	if event.Key() == tcell.KeyRune && (event.Rune() == 't' || event.Rune() == 'T') {
		row, _ := rv.GetSelection()
		if row > 0 {
			selectedData := rv.GetSelectedData()
			if selectedData != nil {
				if rowData, ok := selectedData.(*ResourceRowData); ok && rv.onShowResourceType != nil {
					rv.onShowResourceType(rowData.Resource.Type)
					return nil
				}
			}
		}
	}

	// Handle handler-specific actions (e.g., 'e' for explore, 'd' for details)
	if event.Key() == tcell.KeyRune {
		row, _ := rv.GetSelection()
		if row > 0 {
			selectedData := rv.GetSelectedData()
			if selectedData != nil {
				if rowData, ok := selectedData.(*ResourceRowData); ok {
					handler := rv.registry.GetHandlerOrDefault(rowData.Resource.Type)
					if handler != nil {
						actions := handler.GetActions()
						for _, action := range actions {
							if event.Rune() == action.Key || event.Rune() == action.Key-32 { // Handle both lowercase and uppercase
								actionContext := &resource.ActionContext{
									SubscriptionID:    rv.subscriptionID,
									SubscriptionName:  rv.subscriptionName,
									ResourceGroupName: rv.resourceGroupName,
								}
								if action.Callback != nil && action.Callback(rowData.Resource, actionContext) {
									// Route to appropriate callback
									switch action.Key {
									case 'e', 'E':
										if rv.onExploreStorage != nil && handler.CanExplore() {
											rv.onExploreStorage(rowData.Resource)
										}
									case 'd', 'D':
										if rv.onShowDetails != nil {
											rv.onShowDetails(rowData.Resource)
										}
									}
									return nil
								}
							}
						}
					}
				}
			}
		}
	}

	// Let TableView handle row actions first
	if handled := rv.TableView.HandleKey(event); handled != event {
		return handled
	}
	return event
}
