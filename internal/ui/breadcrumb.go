package ui

import (
	"strings"

	"azure-control-tower/internal/navigation"

	"github.com/rivo/tview"
)

// stripProviderPrefix removes the provider prefix from a resource type
// e.g., "Microsoft.Storage/storageAccounts" -> "storageAccounts"
func stripProviderPrefix(resourceType string) string {
	if idx := strings.LastIndex(resourceType, "/"); idx >= 0 && idx < len(resourceType)-1 {
		return resourceType[idx+1:]
	}
	return resourceType
}

// BreadcrumbView displays the navigation breadcrumb
type BreadcrumbView struct {
	*tview.TextView
	theme *Theme
}

// NewBreadcrumbView creates a new breadcrumb view
func NewBreadcrumbView() *BreadcrumbView {
	theme := DefaultTheme()

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false).
		SetTextAlign(tview.AlignLeft)

	bv := &BreadcrumbView{
		TextView: textView,
		theme:    theme,
	}

	// Set border with theme colors
	textView.SetBorder(true).
		SetBorderColor(theme.Border)

	return bv
}

// Update updates the breadcrumb based on navigation state
func (bv *BreadcrumbView) Update(navState *navigation.State) {
	var breadcrumb strings.Builder

	// Build breadcrumb path based on current view
	switch {
	case navState.InDetailsView:
		// Details view: show full path
		breadcrumb.WriteString("[white]Subscriptions[white]")
		if navState.SelectedSubscriptionName != "" {
			breadcrumb.WriteString(" [gray]>[white] ")
			breadcrumb.WriteString("[white]" + navState.SelectedSubscriptionName + "[white]")
		}
		if navState.SelectedResourceGroupName != "" {
			breadcrumb.WriteString(" [gray]>[white] ")
			breadcrumb.WriteString("[white]" + navState.SelectedResourceGroupName + "[white]")
		}
		breadcrumb.WriteString(" [gray]>[white] ")
		breadcrumb.WriteString("[lightblue::b]Details[white]")
	case navState.CurrentView == navigation.ViewResourceType:
		// Resource type filtered view
		breadcrumb.WriteString("[white]Subscriptions[white]")
		if navState.SelectedSubscriptionName != "" {
			breadcrumb.WriteString(" [gray]>[white] ")
			breadcrumb.WriteString("[white]" + navState.SelectedSubscriptionName + "[white]")
		}
		if navState.SelectedResourceGroupName != "" {
			breadcrumb.WriteString(" [gray]>[white] ")
			breadcrumb.WriteString("[white]" + navState.SelectedResourceGroupName + "[white]")
		}
		breadcrumb.WriteString(" [gray]>[white] ")
		breadcrumb.WriteString("[lightblue::b]" + stripProviderPrefix(navState.SelectedResourceType) + "[white]")
	case navState.CurrentView == navigation.ViewResources:
		// Resources view
		breadcrumb.WriteString("[white]Subscriptions[white]")
		if navState.SelectedSubscriptionName != "" {
			breadcrumb.WriteString(" [gray]>[white] ")
			breadcrumb.WriteString("[white]" + navState.SelectedSubscriptionName + "[white]")
		}
		if navState.SelectedResourceGroupName != "" {
			breadcrumb.WriteString(" [gray]>[white] ")
			breadcrumb.WriteString("[lightblue::b]" + navState.SelectedResourceGroupName + "[white]")
		}
	case navState.CurrentView == navigation.ViewResourceTypes:
		// Resource types view
		breadcrumb.WriteString("[white]Subscriptions[white]")
		if navState.SelectedSubscriptionName != "" {
			breadcrumb.WriteString(" [gray]>[white] ")
			breadcrumb.WriteString("[white]" + navState.SelectedSubscriptionName + "[white]")
		}
		if navState.SelectedResourceGroupName != "" {
			breadcrumb.WriteString(" [gray]>[white] ")
			breadcrumb.WriteString("[lightblue::b]" + navState.SelectedResourceGroupName + "[white]")
		}
	case navState.CurrentView == navigation.ViewResourceGroups:
		// Resource groups view
		breadcrumb.WriteString("[white]Subscriptions[white]")
		if navState.SelectedSubscriptionName != "" {
			breadcrumb.WriteString(" [gray]>[white] ")
			breadcrumb.WriteString("[lightblue::b]" + navState.SelectedSubscriptionName + "[white]")
		}
	case navState.CurrentView == navigation.ViewSubscriptions:
		// Subscriptions view
		breadcrumb.WriteString("[lightblue::b]Subscriptions[white]")
	default:
		breadcrumb.WriteString("[white]Home[white]")
	}

	bv.SetText(breadcrumb.String())
}
