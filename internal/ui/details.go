package ui

import (
	"fmt"
	"strings"

	"azure-control-tower/internal/models"
	"azure-control-tower/pkg/resource"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// DetailsView displays detailed information about a resource
type DetailsView struct {
	*tview.TextView
	registry *resource.Registry
	onBack   func()
	theme    *Theme
}

// NewDetailsView creates a new details view
func NewDetailsView(registry *resource.Registry) *DetailsView {
	theme := DefaultTheme()

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetWordWrap(true)

	dv := &DetailsView{
		TextView: textView,
		registry: registry,
		theme:    theme,
	}

	textView.SetBorder(true).
		SetBorderColor(theme.Border).
		SetTitle("Details (Press ESC to go back)")

	return dv
}

// ShowSubscriptionDetails displays subscription details
func (dv *DetailsView) ShowSubscriptionDetails(sub *models.Subscription) {
	var content strings.Builder
	content.WriteString("[lightblue::b]Subscription Details[white]\n\n")
	content.WriteString(fmt.Sprintf("[lightblue::b]ID:[white] %s\n", sub.ID))
	content.WriteString(fmt.Sprintf("[lightblue::b]Name:[white] %s\n", sub.Name))
	content.WriteString(fmt.Sprintf("[lightblue::b]Display Name:[white] %s\n", sub.DisplayName))
	content.WriteString(fmt.Sprintf("[lightblue::b]State:[white] %s\n", sub.State))
	content.WriteString(fmt.Sprintf("[lightblue::b]Tenant ID:[white] %s\n", sub.TenantID))

	dv.SetText(content.String())
}

// ShowResourceGroupDetails displays resource group details
func (dv *DetailsView) ShowResourceGroupDetails(rg *models.ResourceGroup, subscriptionID string) {
	var content strings.Builder
	content.WriteString("[lightblue::b]Resource Group Details[white]\n\n")
	content.WriteString(fmt.Sprintf("[lightblue::b]Name:[white] %s\n", rg.Name))
	content.WriteString(fmt.Sprintf("[lightblue::b]Location:[white] %s\n", rg.Location))
	content.WriteString(fmt.Sprintf("[lightblue::b]Subscription ID:[white] %s\n", subscriptionID))

	if len(rg.Tags) > 0 {
		content.WriteString("\n[lightblue::b]Tags:[white]\n")
		for key, value := range rg.Tags {
			val := ""
			if value != nil {
				val = *value
			}
			content.WriteString(fmt.Sprintf("  [lightblue::b]%s:[white] %s\n", key, val))
		}
	} else {
		content.WriteString("\n[lightblue::b]Tags:[white] None\n")
	}

	dv.SetText(content.String())
}

// ShowResourceDetails displays resource details
func (dv *DetailsView) ShowResourceDetails(resource *models.Resource, subscriptionID string) {
	// Try to use handler's RenderDetails method
	handler := dv.registry.GetHandlerOrDefault(resource.Type)
	if handler != nil {
		content := handler.RenderDetails(resource, subscriptionID)
		dv.SetText(content)
		return
	}

	// Fallback to default rendering if no handler
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

	dv.SetText(content.String())
}

// ShowContainerDetails displays container details
func (dv *DetailsView) ShowContainerDetails(container *models.Container, storageAccountName string) {
	var content strings.Builder
	content.WriteString("[lightblue::b]Container Details[white]\n\n")
	content.WriteString(fmt.Sprintf("[lightblue::b]Storage Account:[white] %s\n", storageAccountName))
	content.WriteString(fmt.Sprintf("[lightblue::b]Name:[white] %s\n", container.Name))
	content.WriteString(fmt.Sprintf("[lightblue::b]Public Access:[white] %s\n", getPublicAccessDisplay(container.PublicAccess)))
	content.WriteString(fmt.Sprintf("[lightblue::b]Last Modified:[white] %s\n", container.LastModified.Format("2006-01-02 15:04:05")))
	content.WriteString(fmt.Sprintf("[lightblue::b]ETag:[white] %s\n", container.ETag))

	if len(container.Metadata) > 0 {
		content.WriteString("\n[lightblue::b]Metadata:[white]\n")
		for key, value := range container.Metadata {
			content.WriteString(fmt.Sprintf("  [lightblue::b]%s:[white] %s\n", key, value))
		}
	} else {
		content.WriteString("\n[lightblue::b]Metadata:[white] None\n")
	}

	dv.SetText(content.String())
}

// getPublicAccessDisplay returns a display string for public access
func getPublicAccessDisplay(publicAccess string) string {
	if publicAccess == "" {
		return "Private"
	}
	return publicAccess
}

// ShowBlobDetails displays blob details
func (dv *DetailsView) ShowBlobDetails(blob *models.Blob, storageAccountName, containerName string) {
	var content strings.Builder
	content.WriteString("[lightblue::b]Blob Details[white]\n\n")
	content.WriteString(fmt.Sprintf("[lightblue::b]Storage Account:[white] %s\n", storageAccountName))
	content.WriteString(fmt.Sprintf("[lightblue::b]Container:[white] %s\n", containerName))
	content.WriteString(fmt.Sprintf("[lightblue::b]Name:[white] %s\n", blob.Name))
	content.WriteString(fmt.Sprintf("[lightblue::b]Size:[white] %s\n", formatBlobSize(blob.Size)))
	content.WriteString(fmt.Sprintf("[lightblue::b]Content Type:[white] %s\n", blob.ContentType))
	content.WriteString(fmt.Sprintf("[lightblue::b]Last Modified:[white] %s\n", blob.LastModified.Format("2006-01-02 15:04:05")))
	content.WriteString(fmt.Sprintf("[lightblue::b]ETag:[white] %s\n", blob.ETag))

	if len(blob.Metadata) > 0 {
		content.WriteString("\n[lightblue::b]Metadata:[white]\n")
		for key, value := range blob.Metadata {
			content.WriteString(fmt.Sprintf("  [lightblue::b]%s:[white] %s\n", key, value))
		}
	} else {
		content.WriteString("\n[lightblue::b]Metadata:[white] None\n")
	}

	dv.SetText(content.String())
}

// formatBlobSize formats a blob size in bytes to a human-readable string
func formatBlobSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// SetOnBack sets the callback for when ESC is pressed
func (dv *DetailsView) SetOnBack(callback func()) {
	dv.onBack = callback
}

// GetKeyBindings returns key bindings for this view
func (dv *DetailsView) GetKeyBindings() map[tcell.Key]func() {
	return map[tcell.Key]func(){
		tcell.KeyEscape: func() {
			if dv.onBack != nil {
				dv.onBack()
			}
		},
	}
}
