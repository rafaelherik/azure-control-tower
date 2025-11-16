package ui

import (
	"context"

	"azure-control-tower/internal/models"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ResourceGroupRowData wraps ResourceGroup with subscription info for display
type ResourceGroupRowData struct {
	ResourceGroup   *models.ResourceGroup
	SubscriptionID  string
	SubscriptionName string
}

// ResourceGroupsView displays a table of resource groups
type ResourceGroupsView struct {
	*TableView
	resourceGroups   []*models.ResourceGroup
	subscriptionID    string
	subscriptionName  string
	onSelect         func(resourceGroup *models.ResourceGroup)
	onShowDetails    func(resourceGroup *models.ResourceGroup)
}

// NewResourceGroupsView creates a new resource groups view
func NewResourceGroupsView() *ResourceGroupsView {
	rgv := &ResourceGroupsView{}

	// Create table configuration
	config := &TableConfig{
		Title: "", // Will be set in LoadResourceGroups
		Columns: []ColumnConfig{
			{Name: "Name", Align: tview.AlignLeft},
			{Name: "Location", Align: tview.AlignLeft},
			{Name: "Subscription ID", Align: tview.AlignLeft},
		},
		RowActions: []RowAction{
			{
				Rune:  'd',
				Label: "Details",
				Callback: func(rowIndex int, data interface{}) bool {
					if rowData, ok := data.(*ResourceGroupRowData); ok && rgv.onShowDetails != nil {
						rgv.onShowDetails(rowData.ResourceGroup)
						return true
					}
					return false
				},
			},
		},
		OnSelect: func(rowIndex int, data interface{}) {
			if rowData, ok := data.(*ResourceGroupRowData); ok && rgv.onSelect != nil {
				rgv.onSelect(rowData.ResourceGroup)
			}
		},
		GetCellValue: func(data interface{}, columnIndex int) string {
			rowData, ok := data.(*ResourceGroupRowData)
			if !ok {
				return ""
			}
			switch columnIndex {
			case 0:
				return rowData.ResourceGroup.Name
			case 1:
				return rowData.ResourceGroup.Location
			case 2:
				return rowData.SubscriptionID
			default:
				return ""
			}
		},
	}

	rgv.TableView = NewTableView(config)
	return rgv
}

// LoadResourceGroups loads resource groups into the view
func (rgv *ResourceGroupsView) LoadResourceGroups(ctx context.Context, resourceGroups []*models.ResourceGroup, subscriptionID, subscriptionName string) error {
	rgv.resourceGroups = resourceGroups
	rgv.subscriptionID = subscriptionID
	rgv.subscriptionName = subscriptionName

	// Update title
	rgv.SetTitle("")

	// Convert to interface{} slice with row data
	data := make([]interface{}, len(resourceGroups))
	for i, rg := range resourceGroups {
		data[i] = &ResourceGroupRowData{
			ResourceGroup:    rg,
			SubscriptionID:    subscriptionID,
			SubscriptionName:  subscriptionName,
		}
	}

	rgv.LoadData(data)
	return nil
}

// SetOnSelect sets the callback for when a resource group is selected
func (rgv *ResourceGroupsView) SetOnSelect(callback func(*models.ResourceGroup)) {
	rgv.onSelect = callback
}

// SetOnShowDetails sets the callback for when details are requested (d key)
func (rgv *ResourceGroupsView) SetOnShowDetails(callback func(*models.ResourceGroup)) {
	rgv.onShowDetails = callback
}

// GetSubscriptionID returns the current subscription ID
func (rgv *ResourceGroupsView) GetSubscriptionID() string {
	return rgv.subscriptionID
}

// HandleKey handles key events for this view
func (rgv *ResourceGroupsView) HandleKey(event *tcell.EventKey) *tcell.EventKey {
	// Let TableView handle row actions first
	if handled := rgv.TableView.HandleKey(event); handled != event {
		return handled
	}
	return event
}

