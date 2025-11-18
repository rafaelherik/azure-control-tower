package ui

import (
	"fmt"
	"strings"

	"azure-control-tower/internal/models"
	"azure-control-tower/internal/navigation"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	azctLogo = `    █████╗ ███████╗ ██████╗████████╗
   ██╔══██╗╚══███╔╝██╔════╝╚══██╔══╝
   ███████║  ███╔╝ ██║        ██║   
   ██╔══██║ ███╔╝  ██║        ██║   
   ██║  ██║███████╗╚██████╗   ██║   
   ╚═╝  ╚═╝╚══════╝ ╚═════╝   ╚═╝`
)

// HeaderView displays user and tenant information
type HeaderView struct {
	*tview.Flex
	logoView               *tview.TextView
	actionsView            *tview.TextView
	userInfoView           *tview.TextView
	separator1             *tview.Box
	separator2             *tview.Box
	userInfo               *models.UserInfo
	selectedSubscription   string
	selectedSubscriptionID string
	theme                  *Theme
}

// NewHeaderView creates a new header view
func NewHeaderView() *HeaderView {
	// Create user info view (left) - tenant, user, subscriptions
	userInfoView := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false).
		SetTextAlign(tview.AlignLeft)

	theme := DefaultTheme()

	// Create separator 1 (vertical line between user info and actions)
	separator1 := tview.NewBox().
		SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
			// Draw vertical line
			for i := 0; i < height; i++ {
				screen.SetContent(x, y+i, '│', nil, tcell.StyleDefault.Foreground(theme.Secondary))
			}
			return x + width, y, 0, height
		})

	// Create actions view (middle) - keyboard shortcuts in 2 columns, left-aligned
	actionsView := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false).
		SetTextAlign(tview.AlignLeft)

	// Create separator 2 (vertical line between actions and logo)
	separator2 := tview.NewBox().
		SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
			// Draw vertical line
			for i := 0; i < height; i++ {
				screen.SetContent(x, y+i, '│', nil, tcell.StyleDefault.Foreground(theme.Secondary))
			}
			return x + width, y, 0, height
		})

	// Create logo view (right) - azct ASCII art
	logoView := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false).
		SetText(azctLogo)

	// Create main flex container (horizontal)
	// UserInfo: flexible, Separator1: 1 char, Actions: flexible, Separator2: 1 char, Logo: fixed width
	flex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(userInfoView, 0, 1, false).
		AddItem(separator1, 1, 0, false).
		AddItem(actionsView, 0, 1, false).
		AddItem(separator2, 1, 0, false).
		AddItem(logoView, 35, 0, false)

	// Set border with colors from theme
	flex.SetBorder(true).
		SetBorderColor(theme.Border).
		SetTitle("[yellow]Azure Control Tower[white]").
		SetTitleColor(tcell.ColorYellow)

	hv := &HeaderView{
		Flex:         flex,
		logoView:     logoView,
		actionsView:  actionsView,
		userInfoView: userInfoView,
		separator1:   separator1,
		separator2:   separator2,
		theme:        theme,
	}

	// Initialize actions content with empty state (will be updated when navigation state is available)
	hv.updateActions(nil)

	return hv
}

// UpdateActions updates the actions based on the current navigation state
func (hv *HeaderView) UpdateActions(navState *navigation.State) {
	hv.updateActions(navState)
}

// updateActions updates the actions/keyboard shortcuts display in 2 columns
func (hv *HeaderView) updateActions(navState *navigation.State) {
	var actionLines []string
	actionLines = append(actionLines, "[cyan::b]Actions:[white]")

	// Determine available actions based on navigation state
	if navState == nil {
		// Default actions when state is not available
		actionLines = append(actionLines, "[yellow]/[white] - Filter    [yellow]Enter[white] - Select")
		actionLines = append(actionLines, "[yellow]q[white] - Quit      [yellow]d[white] - Details")
		actionLines = append(actionLines, "[yellow]Esc[white] - Back")
		hv.actionsView.SetText(strings.Join(actionLines, "\n"))
		return
	}

	// Build actions based on current view
	var actions []string

	// Filter action - available in all table views, not in details view
	if !navState.InDetailsView {
		actions = append(actions, "[yellow]/[white] - Filter")
	}

	// Menu action - available in all table views, not in details view
	if !navState.InDetailsView {
		actions = append(actions, "[yellow]M[white] - Menu")
	}

	// Enter/Select action - available in subscriptions, resource groups, resource types, storage explorer, blobs, key vault views
	if !navState.InDetailsView {
		switch navState.CurrentView {
		case navigation.ViewSubscriptions, navigation.ViewResourceGroups, navigation.ViewResourceTypes,
			navigation.ViewStorageExplorer, navigation.ViewBlobs,
			navigation.ViewKeyVaultExplorer, navigation.ViewKeyVaultSecrets, navigation.ViewKeyVaultKeys, navigation.ViewKeyVaultCertificates:
			actions = append(actions, "[yellow]Enter[white] - Select")
		}
	}

	// Explore action (E) - available in resources and resource type views for storage accounts and Key Vaults
	if !navState.InDetailsView {
		if navState.CurrentView == navigation.ViewResources ||
			(navState.CurrentView == navigation.ViewResourceType && 
				(navState.SelectedResourceType == "Microsoft.Storage/storageAccounts" || 
				 navState.SelectedResourceType == "Microsoft.KeyVault/vaults")) {
			actions = append(actions, "[yellow]E[white] - Explore")
		}
	}

	// View secret value action (V) - available in Key Vault secrets view
	if !navState.InDetailsView && navState.CurrentView == navigation.ViewKeyVaultSecrets {
		actions = append(actions, "[yellow]V[white] - View Value")
	}

	// Details action (d) - available in subscriptions, resource groups, resources, resource type, storage explorer, blobs, and Key Vault views
	// Not available in resource types view or details view
	if !navState.InDetailsView {
		switch navState.CurrentView {
		case navigation.ViewSubscriptions, navigation.ViewResourceGroups, navigation.ViewResources,
			navigation.ViewResourceType, navigation.ViewStorageExplorer, navigation.ViewBlobs,
			navigation.ViewKeyVaultSecrets, navigation.ViewKeyVaultKeys, navigation.ViewKeyVaultCertificates:
			actions = append(actions, "[yellow]d[white] - Details")
		}
	}

	// Back action (Esc) - available when not at root (subscriptions view)
	if navState.CurrentView != navigation.ViewSubscriptions {
		actions = append(actions, "[yellow::b]Esc[white] - Back")
	}

	// Quit action - always available
	actions = append(actions, "[yellow::b]q[white] - Quit")

	// Format actions in 2 columns, left-aligned
	// Distribute actions across lines (2 per line)
	for i := 0; i < len(actions); i += 2 {
		if i+1 < len(actions) {
			actionLines = append(actionLines, fmt.Sprintf("%s    %s", actions[i], actions[i+1]))
		} else {
			actionLines = append(actionLines, actions[i])
		}
	}

	hv.actionsView.SetText(strings.Join(actionLines, "\n"))
}

// UpdateUserInfo updates the user information displayed in the header
func (hv *HeaderView) UpdateUserInfo(userInfo *models.UserInfo) {
	hv.userInfo = userInfo
	hv.updateContent()
}

// UpdateSelectedSubscription updates the selected subscription in the header
func (hv *HeaderView) UpdateSelectedSubscription(subscriptionName, subscriptionID string) {
	hv.selectedSubscription = subscriptionName
	hv.selectedSubscriptionID = subscriptionID
	hv.updateContent()
}

// updateContent refreshes the header content
func (hv *HeaderView) updateContent() {
	if hv.userInfo == nil {
		hv.userInfoView.SetText("[yellow]Loading user information...[white]")
		return
	}

	var content strings.Builder

	// Format with colors: labels in bold light blue, values in white
	// Using tview color tags: [lightblue:bold] for bold light blue
	content.WriteString(fmt.Sprintf("[lightblue::b]Tenant:[white] %s\n", hv.userInfo.TenantID))

	subscriptionText := "None"
	if hv.selectedSubscription != "" {
		if hv.selectedSubscriptionID != "" {
			subscriptionText = fmt.Sprintf("%s (%s)", hv.selectedSubscription, hv.selectedSubscriptionID)
		} else {
			subscriptionText = hv.selectedSubscription
		}
	}
	content.WriteString(fmt.Sprintf("[lightblue::b]Subscription:[white] %s\n", subscriptionText))
	content.WriteString(fmt.Sprintf("[lightblue::b]User:[white] %s", hv.userInfo.Email))

	hv.userInfoView.SetText(content.String())
}
