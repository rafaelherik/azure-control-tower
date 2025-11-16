package ui

import (
	"context"

	"azure-control-tower/internal/models"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// SubscriptionsView displays a table of subscriptions
type SubscriptionsView struct {
	*TableView
	subscriptions []*models.Subscription
	onSelect      func(subscription *models.Subscription)
	onShowDetails func(subscription *models.Subscription)
}

// NewSubscriptionsView creates a new subscriptions view
func NewSubscriptionsView() *SubscriptionsView {
	sv := &SubscriptionsView{}

	// Create table configuration
	config := &TableConfig{
		Title: "",
		Columns: []ColumnConfig{
			{Name: "ID", Align: tview.AlignLeft},
			{Name: "Name", Align: tview.AlignLeft},
			{Name: "Tenant ID", Align: tview.AlignLeft},
		},
		RowActions: []RowAction{
			{
				Key:   tcell.KeyEnter,
				Label: "Select",
				Callback: func(rowIndex int, data interface{}) bool {
					if sub, ok := data.(*models.Subscription); ok && sv.onSelect != nil {
						sv.onSelect(sub)
						return true
					}
					return false
				},
			},
			{
				Rune:  'd',
				Label: "Details",
				Callback: func(rowIndex int, data interface{}) bool {
					if sub, ok := data.(*models.Subscription); ok && sv.onShowDetails != nil {
						sv.onShowDetails(sub)
						return true
					}
					return false
				},
			},
		},
		OnSelect: func(rowIndex int, data interface{}) {
			if sub, ok := data.(*models.Subscription); ok && sv.onSelect != nil {
				sv.onSelect(sub)
			}
		},
		GetCellValue: func(data interface{}, columnIndex int) string {
			sub, ok := data.(*models.Subscription)
			if !ok {
				return ""
			}
			switch columnIndex {
			case 0:
				return sub.ID
			case 1:
				name := sub.DisplayName
				if name == "" {
					name = sub.Name
				}
				if name == "" {
					name = sub.ID
				}
				return name
			case 2:
				return sub.TenantID
			default:
				return ""
			}
		},
	}

	sv.TableView = NewTableView(config)
	return sv
}

// LoadSubscriptions loads subscriptions into the view
func (sv *SubscriptionsView) LoadSubscriptions(ctx context.Context, subscriptions []*models.Subscription) error {
	sv.subscriptions = subscriptions

	// Convert to interface{} slice
	data := make([]interface{}, len(subscriptions))
	for i, sub := range subscriptions {
		data[i] = sub
	}

	sv.LoadData(data)
	return nil
}

// SetOnSelect sets the callback for when a subscription is selected (Enter key)
func (sv *SubscriptionsView) SetOnSelect(callback func(*models.Subscription)) {
	sv.onSelect = callback
}

// SetOnShowDetails sets the callback for when details are requested (d key)
func (sv *SubscriptionsView) SetOnShowDetails(callback func(*models.Subscription)) {
	sv.onShowDetails = callback
}

// HandleKey handles key events for this view
func (sv *SubscriptionsView) HandleKey(event *tcell.EventKey) *tcell.EventKey {
	// Let TableView handle row actions first
	if handled := sv.TableView.HandleKey(event); handled != event {
		return handled
	}
	return event
}
