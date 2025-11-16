package ui

import (
	"context"
	"fmt"
	"strings"

	"azure-control-tower/internal/models"
	"azure-control-tower/pkg/resource"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ResourceTypesView displays a table of resource types with counts
type ResourceTypesView struct {
	*TableView
	resourceTypes     []*models.ResourceTypeSummary
	subscriptionID    string
	subscriptionName  string
	resourceGroupName string
	registry          *resource.Registry
	onSelect          func(resourceType *models.ResourceTypeSummary)
}

// NewResourceTypesView creates a new resource types view
func NewResourceTypesView(registry *resource.Registry) *ResourceTypesView {
	rtv := &ResourceTypesView{
		registry: registry,
	}

	// Create table configuration
	config := &TableConfig{
		Title: "", // Will be set in LoadResourceTypes
		Columns: []ColumnConfig{
			{Name: "Resource Type", Align: tview.AlignLeft},
			{Name: "Count", Align: tview.AlignRight},
		},
		RowActions: []RowAction{},
		OnSelect: func(rowIndex int, data interface{}) {
			if rowData, ok := data.(*models.ResourceTypeSummary); ok && rtv.onSelect != nil {
				// Check if handler supports navigation to list
				handler := rtv.registry.GetHandlerOrDefault(rowData.Type)
				if handler != nil && handler.CanNavigateToList() {
					rtv.onSelect(rowData)
				}
			}
		},
		GetCellValue: func(data interface{}, columnIndex int) string {
			rowData, ok := data.(*models.ResourceTypeSummary)
			if !ok {
				return ""
			}
			switch columnIndex {
			case 0:
				// Strip provider prefix (e.g., "Microsoft.Storage/storageAccounts" -> "storageAccounts")
				resourceType := rowData.Type
				if idx := strings.LastIndex(resourceType, "/"); idx >= 0 && idx < len(resourceType)-1 {
					return resourceType[idx+1:]
				}
				return resourceType
			case 1:
				return fmt.Sprintf("%d", rowData.Count)
			default:
				return ""
			}
		},
	}

	rtv.TableView = NewTableView(config)
	return rtv
}

// LoadResourceTypes loads resource types into the view
func (rtv *ResourceTypesView) LoadResourceTypes(ctx context.Context, resourceTypes []*models.ResourceTypeSummary, subscriptionID, subscriptionName, resourceGroupName string) error {
	rtv.resourceTypes = resourceTypes
	rtv.subscriptionID = subscriptionID
	rtv.subscriptionName = subscriptionName
	rtv.resourceGroupName = resourceGroupName

	// Update title
	rtv.SetTitle("")

	// Convert to interface{} slice
	data := make([]interface{}, len(resourceTypes))
	for i, rt := range resourceTypes {
		data[i] = rt
	}

	rtv.LoadData(data)

	return nil
}

// RenderData overrides TableView's RenderData to apply styling after rendering
func (rtv *ResourceTypesView) RenderData() {
	rtv.TableView.RenderData()
	rtv.applyStyling()
}

// SetFilter overrides TableView's SetFilter to reapply styling after filtering
func (rtv *ResourceTypesView) SetFilter(filterText string) {
	rtv.TableView.SetFilter(filterText)
	rtv.applyStyling()
}

// applyStyling applies visual styling to indicate which types are navigable
func (rtv *ResourceTypesView) applyStyling() {
	// Iterate through all rows and style based on navigability
	for row := 1; row < rtv.GetRowCount(); row++ {
		// Get the data index for this display row
		displayRowIndex := row - 1
		if displayRowIndex >= 0 && displayRowIndex < len(rtv.filteredIndices) {
			dataIndex := rtv.filteredIndices[displayRowIndex]
			if dataIndex >= 0 && dataIndex < len(rtv.data) {
				if summary, ok := rtv.data[dataIndex].(*models.ResourceTypeSummary); ok {
					// Check if handler supports navigation
					handler := rtv.registry.GetHandlerOrDefault(summary.Type)
					canNavigate := handler != nil && handler.CanNavigateToList()
					// Gray out non-navigable types
					if !canNavigate {
						for col := 0; col < len(rtv.config.Columns); col++ {
							if cell := rtv.GetCell(row, col); cell != nil {
								cell.SetTextColor(tcell.ColorGray)
							}
						}
					}
				}
			}
		}
	}
}

// SetOnSelect sets the callback for when a resource type is selected
func (rtv *ResourceTypesView) SetOnSelect(callback func(*models.ResourceTypeSummary)) {
	rtv.onSelect = callback
}

// GetSubscriptionID returns the current subscription ID
func (rtv *ResourceTypesView) GetSubscriptionID() string {
	return rtv.subscriptionID
}

// GetResourceGroupName returns the current resource group name
func (rtv *ResourceTypesView) GetResourceGroupName() string {
	return rtv.resourceGroupName
}

// HandleKey handles key events for this view
func (rtv *ResourceTypesView) HandleKey(event *tcell.EventKey) *tcell.EventKey {
	// Handle Enter key - only allow for navigable resource types
	if event.Key() == tcell.KeyEnter {
		row, _ := rtv.GetSelection()
		if row > 0 {
			selectedData := rtv.GetSelectedData()
			if selectedData != nil {
				if summary, ok := selectedData.(*models.ResourceTypeSummary); ok {
					// Check if handler supports navigation
					handler := rtv.registry.GetHandlerOrDefault(summary.Type)
					if handler != nil && handler.CanNavigateToList() && rtv.onSelect != nil {
						rtv.onSelect(summary)
						return nil
					}
				}
			}
		}
		// Consume the event even if not navigable to prevent default behavior
		return nil
	}

	// Let TableView handle other keys
	if handled := rtv.TableView.HandleKey(event); handled != event {
		return handled
	}
	return event
}
