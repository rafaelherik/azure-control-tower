package ui

import (
	"context"
	"strings"

	"azure-control-tower/pkg/resource"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// MenuRowData wraps resource type info for display
type MenuRowData struct {
	ResourceType string
	DisplayName  string
}

// MenuView displays a menu of supported resource types
type MenuView struct {
	*TableView
	registry          *resource.Registry
	resourceTypes     []string
	onSelect         func(resourceType string)
	subscriptionID   string
	subscriptionName  string
	resourceGroupName string
}

// NewMenuView creates a new menu view
func NewMenuView(registry *resource.Registry) *MenuView {
	mv := &MenuView{
		registry: registry,
	}

	// Create table configuration
	config := &TableConfig{
		Title: "Resource Types Menu",
		Columns: []ColumnConfig{
			{Name: "Resource Type", Align: tview.AlignLeft},
		},
		RowActions: []RowAction{},
		OnSelect: func(rowIndex int, data interface{}) {
			if rowData, ok := data.(*MenuRowData); ok && mv.onSelect != nil {
				mv.onSelect(rowData.ResourceType)
			}
		},
		GetCellValue: func(data interface{}, columnIndex int) string {
			rowData, ok := data.(*MenuRowData)
			if !ok {
				return ""
			}
			switch columnIndex {
			case 0:
				return rowData.DisplayName
			default:
				return ""
			}
		},
	}

	mv.TableView = NewTableView(config)
	return mv
}

// LoadResourceTypes loads supported resource types into the menu
func (mv *MenuView) LoadResourceTypes(ctx context.Context, subscriptionID, subscriptionName, resourceGroupName string) error {
	mv.subscriptionID = subscriptionID
	mv.subscriptionName = subscriptionName
	mv.resourceGroupName = resourceGroupName

	// Get supported resource types from registry
	mv.resourceTypes = mv.registry.GetSupportedResourceTypes()

	// Convert to interface{} slice with row data
	data := make([]interface{}, len(mv.resourceTypes))
	for i, resourceType := range mv.resourceTypes {
		// Get display name from handler
		handler, err := mv.registry.GetHandler(resourceType)
		displayName := resourceType
		if err == nil && handler != nil {
			displayName = handler.GetDisplayName()
		} else {
			// Fallback: strip provider prefix
			if idx := strings.LastIndex(resourceType, "/"); idx >= 0 && idx < len(resourceType)-1 {
				displayName = resourceType[idx+1:]
			}
		}

		data[i] = &MenuRowData{
			ResourceType: resourceType,
			DisplayName:  displayName,
		}
	}

	mv.LoadData(data)
	return nil
}

// SetOnSelect sets the callback for when a resource type is selected
func (mv *MenuView) SetOnSelect(callback func(resourceType string)) {
	mv.onSelect = callback
}

// HandleKey handles key events for this view
func (mv *MenuView) HandleKey(event *tcell.EventKey) *tcell.EventKey {
	// Let TableView handle keys (including Enter for selection)
	if handled := mv.TableView.HandleKey(event); handled != event {
		return handled
	}
	return event
}


