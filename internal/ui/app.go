package ui

import (
	"context"
	"fmt"
	"strings"

	"azure-control-tower/internal/azure"
	"azure-control-tower/internal/models"
	"azure-control-tower/internal/navigation"
	"azure-control-tower/pkg/resource"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App wraps the tview application with navigation and state management
type App struct {
	*tview.Application
	azureClient         *azure.Client
	registry            *resource.Registry
	navState            *navigation.State
	headerView          *HeaderView
	breadcrumbView      *BreadcrumbView
	viewTitleView       *ViewTitleView
	footerView          *FooterView
	subscriptionsView   *SubscriptionsView
	resourceGroupsView  *ResourceGroupsView
	resourceTypesView   *ResourceTypesView
	resourcesView       *ResourcesView
	detailsView               *DetailsView
	storageExplorerView       *StorageExplorerView
	blobsView                 *BlobsView
	keyVaultExplorerView      *KeyVaultExplorerView
	keyVaultSecretsView       *KeyVaultSecretsView
	keyVaultKeysView          *KeyVaultKeysView
	keyVaultCertificatesView  *KeyVaultCertificatesView
	menuView                  *MenuView
	filterMode          *FilterMode
	mainFlex            *tview.Flex
	currentView         tview.Primitive
	userInfo            *models.UserInfo
}

// NewApp creates a new application instance
func NewApp(azureClient *azure.Client, registry *resource.Registry) *App {
	app := tview.NewApplication()

	navState := navigation.NewState()
	headerView := NewHeaderView()
	breadcrumbView := NewBreadcrumbView()
	viewTitleView := NewViewTitleView()
	footerView := NewFooterView()
	subscriptionsView := NewSubscriptionsView()
	resourceGroupsView := NewResourceGroupsView()
	resourceTypesView := NewResourceTypesView(registry)
	resourcesView := NewResourcesView(registry)
	detailsView := NewDetailsView(registry)
	storageExplorerView := NewStorageExplorerView()
	blobsView := NewBlobsView()
	keyVaultExplorerView := NewKeyVaultExplorerView()
	keyVaultSecretsView := NewKeyVaultSecretsView()
	keyVaultKeysView := NewKeyVaultKeysView()
	keyVaultCertificatesView := NewKeyVaultCertificatesView()
	menuView := NewMenuView(registry)
	filterMode := NewFilterMode(app)

	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow)

	a := &App{
		Application:         app,
		azureClient:         azureClient,
		registry:            registry,
		navState:            navState,
		headerView:          headerView,
		breadcrumbView:      breadcrumbView,
		viewTitleView:       viewTitleView,
		footerView:          footerView,
		subscriptionsView:   subscriptionsView,
		resourceGroupsView:  resourceGroupsView,
		resourceTypesView:   resourceTypesView,
		resourcesView:       resourcesView,
		detailsView:              detailsView,
		storageExplorerView:      storageExplorerView,
		blobsView:                blobsView,
		keyVaultExplorerView:     keyVaultExplorerView,
		keyVaultSecretsView:      keyVaultSecretsView,
		keyVaultKeysView:         keyVaultKeysView,
		keyVaultCertificatesView: keyVaultCertificatesView,
		menuView:                 menuView,
		filterMode:          filterMode,
		mainFlex:            mainFlex,
		currentView:         subscriptionsView,
	}

	// Set up subscriptions view callbacks
	subscriptionsView.SetOnSelect(func(sub *models.Subscription) {
		a.navigateToResourceGroups(sub.ID, sub.DisplayName)
	})

	subscriptionsView.SetOnShowDetails(func(sub *models.Subscription) {
		a.showSubscriptionDetails(sub)
	})

	// Set up resource groups view callbacks
	resourceGroupsView.SetOnSelect(func(rg *models.ResourceGroup) {
		a.navigateToResourceTypes(rg.Name)
	})
	resourceGroupsView.SetOnShowDetails(func(rg *models.ResourceGroup) {
		a.showResourceGroupDetails(rg)
	})

	// Set up resource types view callbacks
	resourceTypesView.SetOnSelect(func(summary *models.ResourceTypeSummary) {
		a.navigateToResourceType(summary.Type)
	})

	// Set up menu view callbacks
	menuView.SetOnSelect(func(resourceType string) {
		a.navigateToResourceTypeFromMenu(resourceType)
	})

	// Set up resources view callbacks
	resourcesView.SetOnShowResourceType(func(resourceType string) {
		a.navigateToResourceType(resourceType)
	})
	resourcesView.SetOnShowDetails(func(resource *models.Resource) {
		a.showResourceDetails(resource)
	})
	resourcesView.SetOnExploreStorage(func(resource *models.Resource) {
		// Route to appropriate explorer based on resource type
		switch resource.Type {
		case "Microsoft.Storage/storageAccounts":
			a.navigateToStorageExplorer(resource)
		case "Microsoft.KeyVault/vaults":
			a.navigateToKeyVaultExplorer(resource)
		default:
			// Try to determine if it's an explorable resource type
			handler := registry.GetHandlerOrDefault(resource.Type)
			if handler != nil && handler.CanExplore() {
				// For now, just navigate to storage explorer as fallback
				// In the future, we can add more handlers
				a.navigateToStorageExplorer(resource)
			}
		}
	})

	// Set up storage explorer view callbacks
	storageExplorerView.SetOnSelect(func(container *models.Container) {
		a.navigateToBlobs(container.Name)
	})
	storageExplorerView.SetOnShowDetails(func(container *models.Container) {
		a.showContainerDetails(container)
	})

	// Set up blobs view callbacks
	blobsView.SetOnShowDetails(func(blob *models.Blob) {
		a.showBlobDetails(blob)
	})
	blobsView.SetOnNavigateFolder(func(folderPath string) {
		a.navigateIntoBlobFolder(folderPath)
	})

	// Set up Key Vault explorer view callbacks
	keyVaultExplorerView.SetOnSelect(func(itemType string) {
		a.navigateToKeyVaultItemType(itemType)
	})

	// Set up Key Vault secrets view callbacks
	keyVaultSecretsView.SetOnShowDetails(func(secret *models.Secret) {
		a.showSecretDetails(secret)
	})
	keyVaultSecretsView.SetOnViewValue(func(secret *models.Secret) {
		a.viewSecretValue(secret)
	})

	// Set up Key Vault keys view callbacks
	keyVaultKeysView.SetOnShowDetails(func(key *models.Key) {
		a.showKeyDetails(key)
	})

	// Set up Key Vault certificates view callbacks
	keyVaultCertificatesView.SetOnShowDetails(func(cert *models.Certificate) {
		a.showCertificateDetails(cert)
	})

	// Set up details view callback
	detailsView.SetOnBack(func() {
		a.navigateBackFromDetails()
	})

	// Set up filter mode
	filterMode.SetOnFilter(func(filterText string) {
		a.applyFilter(filterText)
	})

	filterMode.SetOnCancel(func() {
		a.clearFilter()
		a.updateLayout()
		app.SetFocus(a.currentView)
	})

	// Set up key bindings
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if filterMode.IsVisible() {
			// Let filter mode handle its own keys
			return event
		}

		// Handle details view navigation
		if navState.InDetailsView {
			if event.Key() == tcell.KeyEscape {
				a.navigateBackFromDetails()
				return nil
			}
			// Let details view handle its own keys
			return event
		}

		// Handle key bindings for table views
		switch navState.CurrentView {
		case navigation.ViewSubscriptions:
			if handled := subscriptionsView.HandleKey(event); handled != event {
				return handled
			}
		case navigation.ViewResourceGroups:
			if handled := resourceGroupsView.HandleKey(event); handled != event {
				return handled
			}
		case navigation.ViewResourceTypes:
			if handled := resourceTypesView.HandleKey(event); handled != event {
				return handled
			}
		case navigation.ViewResources:
			if handled := resourcesView.HandleKey(event); handled != event {
				return handled
			}
		case navigation.ViewResourceType:
			if handled := resourcesView.HandleKey(event); handled != event {
				return handled
			}
		case navigation.ViewStorageExplorer:
			if handled := storageExplorerView.HandleKey(event); handled != event {
				return handled
			}
		case navigation.ViewBlobs:
			if handled := blobsView.HandleKey(event); handled != event {
				return handled
			}
		case navigation.ViewKeyVaultExplorer:
			if handled := keyVaultExplorerView.HandleKey(event); handled != event {
				return handled
			}
		case navigation.ViewKeyVaultSecrets:
			if handled := keyVaultSecretsView.HandleKey(event); handled != event {
				return handled
			}
		case navigation.ViewKeyVaultKeys:
			if handled := keyVaultKeysView.HandleKey(event); handled != event {
				return handled
			}
		case navigation.ViewKeyVaultCertificates:
			if handled := keyVaultCertificatesView.HandleKey(event); handled != event {
				return handled
			}
		case navigation.ViewMenu:
			if handled := menuView.HandleKey(event); handled != event {
				return handled
			}
		}

		switch event.Key() {
		case tcell.KeyEscape:
			switch navState.CurrentView {
			case navigation.ViewMenu:
				// Go back from menu to previous view
				a.navigateBackFromMenu()
				return nil
			case navigation.ViewBlobs:
				// Go back to storage explorer
				a.navigateBackFromBlobs()
				return nil
			case navigation.ViewKeyVaultSecrets, navigation.ViewKeyVaultKeys, navigation.ViewKeyVaultCertificates:
				// Go back to Key Vault explorer
				a.navigateBackToKeyVaultExplorer()
				return nil
			case navigation.ViewKeyVaultExplorer:
				// Go back to resource type view
				a.navigateBackToResourceType()
				return nil
			case navigation.ViewStorageExplorer:
				// Go back to resource type view
				a.navigateBackToResourceType()
				return nil
			case navigation.ViewResourceType:
				// Go back to resource types view
				a.navigateBackToResourceTypes()
				return nil
			case navigation.ViewResourceTypes:
				// Go back to resource groups view
				a.navigateBackToResourceGroups()
				return nil
			case navigation.ViewResourceGroups:
				a.navigateToSubscriptions()
				return nil
			}
		case tcell.KeyRune:
			switch event.Rune() {
			case '/':
				// Only allow filtering in table views
				if !navState.InDetailsView {
					filterMode.Show()
					a.updateLayout()
					return nil
				}
			case 'm', 'M':
				// Open menu (only when not in details view)
				if !navState.InDetailsView {
					a.navigateToMenu()
					return nil
				}
			case 'q':
				app.Stop()
				return nil
			}
		}

		return event
	})

	// Initial layout - use updateLayout to ensure header is visible
	a.updateLayout()

	return a
}

// updateLayout updates the main layout based on current state
func (a *App) updateLayout() {
	a.mainFlex.Clear()

	// Add header (always visible at top) - reduced height for smaller logo
	a.mainFlex.AddItem(a.headerView, 12, 0, false)

	// Update header actions based on navigation state
	a.headerView.UpdateActions(a.navState)

	// Update breadcrumb
	a.breadcrumbView.Update(a.navState)

	// Add breadcrumb (between header and content)
	a.mainFlex.AddItem(a.breadcrumbView, 1, 0, false)

	// Add view title (between breadcrumb and content) - only for table views
	if !a.navState.InDetailsView {
		a.updateViewTitle()
		a.mainFlex.AddItem(a.viewTitleView, 1, 0, false)
	}

	// Add filter mode if visible (between view title and content)
	if a.filterMode.IsVisible() {
		a.mainFlex.AddItem(a.filterMode.GetInputField(), 1, 0, true)
	}

	// Add main content view (details, subscriptions, resource groups, resources, resource type, storage explorer, or blobs)
	if a.navState.InDetailsView {
		a.mainFlex.AddItem(a.detailsView, 0, 1, true)
		a.currentView = a.detailsView
		a.updateFooterWithActions(0, 0, false, "ESC: back, q: quit") // No count for details view
	} else if a.navState.CurrentView == navigation.ViewSubscriptions {
		a.mainFlex.AddItem(a.subscriptionsView, 0, 1, true)
		a.currentView = a.subscriptionsView
		a.updateFooterForTableView(a.subscriptionsView.TableView)
	} else if a.navState.CurrentView == navigation.ViewResourceGroups {
		a.mainFlex.AddItem(a.resourceGroupsView, 0, 1, true)
		a.currentView = a.resourceGroupsView
		a.updateFooterForTableView(a.resourceGroupsView.TableView)
	} else if a.navState.CurrentView == navigation.ViewResourceTypes {
		a.mainFlex.AddItem(a.resourceTypesView, 0, 1, true)
		a.currentView = a.resourceTypesView
		a.updateFooterForTableView(a.resourceTypesView.TableView)
	} else if a.navState.CurrentView == navigation.ViewResources || a.navState.CurrentView == navigation.ViewResourceType {
		a.mainFlex.AddItem(a.resourcesView, 0, 1, true)
		a.currentView = a.resourcesView
		a.updateFooterForTableView(a.resourcesView.TableView)
	} else if a.navState.CurrentView == navigation.ViewStorageExplorer {
		a.mainFlex.AddItem(a.storageExplorerView, 0, 1, true)
		a.currentView = a.storageExplorerView
		a.updateFooterForTableView(a.storageExplorerView.TableView)
	} else if a.navState.CurrentView == navigation.ViewBlobs {
		a.mainFlex.AddItem(a.blobsView, 0, 1, true)
		a.currentView = a.blobsView
		a.updateFooterForTableView(a.blobsView.TableView)
	} else if a.navState.CurrentView == navigation.ViewKeyVaultExplorer {
		a.mainFlex.AddItem(a.keyVaultExplorerView, 0, 1, true)
		a.currentView = a.keyVaultExplorerView
		a.updateFooterForTableView(a.keyVaultExplorerView.TableView)
	} else if a.navState.CurrentView == navigation.ViewKeyVaultSecrets {
		a.mainFlex.AddItem(a.keyVaultSecretsView, 0, 1, true)
		a.currentView = a.keyVaultSecretsView
		a.updateFooterForTableView(a.keyVaultSecretsView.TableView)
	} else if a.navState.CurrentView == navigation.ViewKeyVaultKeys {
		a.mainFlex.AddItem(a.keyVaultKeysView, 0, 1, true)
		a.currentView = a.keyVaultKeysView
		a.updateFooterForTableView(a.keyVaultKeysView.TableView)
	} else if a.navState.CurrentView == navigation.ViewKeyVaultCertificates {
		a.mainFlex.AddItem(a.keyVaultCertificatesView, 0, 1, true)
		a.currentView = a.keyVaultCertificatesView
		a.updateFooterForTableView(a.keyVaultCertificatesView.TableView)
	} else if a.navState.CurrentView == navigation.ViewMenu {
		a.mainFlex.AddItem(a.menuView, 0, 1, true)
		a.currentView = a.menuView
		a.updateFooterForTableView(a.menuView.TableView)
	}

	// Add footer at the bottom
	a.mainFlex.AddItem(a.footerView, 1, 0, false)

	a.SetRoot(a.mainFlex, true)
}

// updateFooterForTableView updates the footer based on a table view
func (a *App) updateFooterForTableView(tableView *TableView) {
	totalCount := len(tableView.data)
	filteredCount := tableView.GetDataRowCount()
	hasFilter := tableView.GetFilter() != ""

	// Get action keys based on current view
	var actions string
	switch a.navState.CurrentView {
	case navigation.ViewSubscriptions:
		actions = "Enter: view Resource Groups, d: details, ESC: back, /: filter, q: quit"
	case navigation.ViewResourceGroups:
		actions = "Enter: view Resource List, d: details, ESC: back, /: filter, q: quit"
	case navigation.ViewResourceTypes:
		actions = "Enter: view storage accounts, ESC: back, /: filter, q: quit"
	case navigation.ViewResources:
		actions = "E: explore storage, d: details, ESC: back, /: filter, q: quit"
	case navigation.ViewResourceType:
		// Get actions from handler
		handler := a.registry.GetHandlerOrDefault(a.navState.SelectedResourceType)
		if handler != nil && handler.CanExplore() {
			actions = "E: explore, d: details, ESC: back, /: filter, q: quit"
		} else {
			actions = "d: details, ESC: back, /: filter, q: quit"
		}
	case navigation.ViewStorageExplorer:
		actions = "Enter: open container, d: details, ESC: back, /: filter, q: quit"
	case navigation.ViewBlobs:
		actions = "Enter: open folder/details, d: details, ESC: back, /: filter, q: quit"
	case navigation.ViewKeyVaultExplorer:
		actions = "Enter: open item type, ESC: back, /: filter, q: quit"
	case navigation.ViewKeyVaultSecrets:
		actions = "v: view value, d: details, ESC: back, /: filter, q: quit"
	case navigation.ViewKeyVaultKeys:
		actions = "d: details, ESC: back, /: filter, q: quit"
	case navigation.ViewKeyVaultCertificates:
		actions = "d: details, ESC: back, /: filter, q: quit"
	case navigation.ViewMenu:
		actions = "Enter: select resource type, ESC: back, /: filter, q: quit"
	default:
		actions = "ESC: back, /: filter, q: quit"
	}

	a.updateFooterWithActions(totalCount, filteredCount, hasFilter, actions)
}

// updateViewTitle updates the view title based on current navigation state
func (a *App) updateViewTitle() {
	var viewName string
	switch a.navState.CurrentView {
	case navigation.ViewSubscriptions:
		viewName = "Subscriptions"
	case navigation.ViewResourceGroups:
		viewName = fmt.Sprintf("Resource Groups - %s", a.navState.SelectedSubscriptionName)
	case navigation.ViewResourceTypes:
		viewName = fmt.Sprintf("Resource List - %s", a.navState.SelectedResourceGroupName)
	case navigation.ViewResources:
		viewName = fmt.Sprintf("Resources - %s", a.navState.SelectedResourceGroupName)
	case navigation.ViewResourceType:
		// Get display name from handler
		handler := a.registry.GetHandlerOrDefault(a.navState.SelectedResourceType)
		if handler != nil {
			viewName = fmt.Sprintf("%s - %s", handler.GetDisplayName(), a.navState.SelectedResourceGroupName)
		} else {
			// Strip provider prefix from resource type for display
			resourceTypeDisplay := a.navState.SelectedResourceType
			if idx := strings.LastIndex(resourceTypeDisplay, "/"); idx >= 0 && idx < len(resourceTypeDisplay)-1 {
				resourceTypeDisplay = resourceTypeDisplay[idx+1:]
			}
			viewName = fmt.Sprintf("Resources - %s (%s)", a.navState.SelectedResourceGroupName, resourceTypeDisplay)
		}
	case navigation.ViewStorageExplorer:
		viewName = fmt.Sprintf("Storage Explorer - %s", a.navState.SelectedStorageAccount)
	case navigation.ViewBlobs:
		pathDisplay := ""
		if a.navState.BlobPathPrefix != "" {
			pathDisplay = fmt.Sprintf(" - %s", a.navState.BlobPathPrefix)
		}
		viewName = fmt.Sprintf("Blobs - %s/%s%s", a.navState.SelectedStorageAccount, a.navState.SelectedContainer, pathDisplay)
	case navigation.ViewKeyVaultExplorer:
		viewName = fmt.Sprintf("Key Vault Explorer - %s", a.navState.SelectedKeyVault)
	case navigation.ViewKeyVaultSecrets:
		viewName = fmt.Sprintf("Secrets - %s", a.navState.SelectedKeyVault)
	case navigation.ViewKeyVaultKeys:
		viewName = fmt.Sprintf("Keys - %s", a.navState.SelectedKeyVault)
	case navigation.ViewKeyVaultCertificates:
		viewName = fmt.Sprintf("Certificates - %s", a.navState.SelectedKeyVault)
	case navigation.ViewMenu:
		viewName = "Resource Types Menu"
	default:
		viewName = "Unknown View"
	}
	a.viewTitleView.SetViewName(viewName)
}

// updateFooterWithActions updates the footer with counts and action keys
func (a *App) updateFooterWithActions(totalCount, filteredCount int, hasFilter bool, actions string) {
	a.footerView.UpdateCountWithActions(totalCount, filteredCount, hasFilter, actions)
}

// Start initializes and runs the application
func (a *App) Start(ctx context.Context) error {
	// Load user info
	userInfo, err := a.azureClient.GetUserInfo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}
	a.userInfo = userInfo
	a.headerView.UpdateUserInfo(userInfo)

	// Load initial subscriptions
	if err := a.loadSubscriptions(ctx); err != nil {
		return fmt.Errorf("failed to load subscriptions: %w", err)
	}

	return a.Run()
}

// loadSubscriptions loads and displays subscriptions
func (a *App) loadSubscriptions(ctx context.Context) error {
	subscriptions, err := a.azureClient.ListSubscriptions(ctx)
	if err != nil {
		return err
	}

	err = a.subscriptionsView.LoadSubscriptions(ctx, subscriptions)
	if err == nil {
		a.updateFooterForTableView(a.subscriptionsView.TableView)
	}
	return err
}

// navigateToSubscriptions navigates back to subscriptions view
func (a *App) navigateToSubscriptions() {
	a.navState.NavigateToSubscriptions()
	a.headerView.UpdateSelectedSubscription("", "")
	a.updateLayout()
	a.SetFocus(a.subscriptionsView)
}

// navigateToResourceGroups navigates to resource groups view for a subscription
func (a *App) navigateToResourceGroups(subscriptionID, subscriptionName string) {
	a.navState.NavigateToResourceGroups(subscriptionID, subscriptionName)
	a.headerView.UpdateSelectedSubscription(subscriptionName, subscriptionID)

	// Load resource groups
	ctx := context.Background()
	resourceGroups, err := a.azureClient.ListResourceGroups(ctx, subscriptionID)
	if err != nil {
		// TODO: Show error in UI
		return
	}

	err = a.resourceGroupsView.LoadResourceGroups(ctx, resourceGroups, subscriptionID, subscriptionName)
	if err == nil {
		a.updateFooterForTableView(a.resourceGroupsView.TableView)
	}
	a.updateLayout()
	a.SetFocus(a.resourceGroupsView)
}

// showSubscriptionDetails shows the details view for a subscription
func (a *App) showSubscriptionDetails(sub *models.Subscription) {
	a.navState.NavigateToDetails()
	a.detailsView.ShowSubscriptionDetails(sub)
	a.updateLayout()
	a.SetFocus(a.detailsView)
}

// showResourceGroupDetails shows the details view for a resource group
func (a *App) showResourceGroupDetails(rg *models.ResourceGroup) {
	a.navState.NavigateToDetails()
	subscriptionID := a.resourceGroupsView.GetSubscriptionID()
	a.detailsView.ShowResourceGroupDetails(rg, subscriptionID)
	a.updateLayout()
	a.SetFocus(a.detailsView)
}

// navigateBackFromDetails returns from details view to previous view
func (a *App) navigateBackFromDetails() {
	a.navState.NavigateBackFromDetails()
	a.updateLayout()
	switch a.navState.CurrentView {
	case navigation.ViewSubscriptions:
		a.SetFocus(a.subscriptionsView)
	case navigation.ViewResourceGroups:
		a.SetFocus(a.resourceGroupsView)
	case navigation.ViewResourceTypes:
		a.SetFocus(a.resourceTypesView)
	case navigation.ViewResources, navigation.ViewResourceType:
		a.SetFocus(a.resourcesView)
	case navigation.ViewStorageExplorer:
		a.SetFocus(a.storageExplorerView)
	case navigation.ViewBlobs:
		a.SetFocus(a.blobsView)
	case navigation.ViewKeyVaultSecrets:
		a.SetFocus(a.keyVaultSecretsView)
	case navigation.ViewKeyVaultKeys:
		a.SetFocus(a.keyVaultKeysView)
	case navigation.ViewKeyVaultCertificates:
		a.SetFocus(a.keyVaultCertificatesView)
	case navigation.ViewMenu:
		a.SetFocus(a.menuView)
	}
}

// navigateToResourceTypes navigates to the resource types summary view for a resource group
func (a *App) navigateToResourceTypes(resourceGroupName string) {
	a.navState.NavigateToResourceTypes(resourceGroupName)

	// Load resource type counts
	ctx := context.Background()
	subscriptionID := a.resourceGroupsView.GetSubscriptionID()
	subscriptionName := a.navState.SelectedSubscriptionName
	resourceTypes, err := a.azureClient.GetResourceTypeCounts(ctx, subscriptionID, resourceGroupName)
	if err != nil {
		// TODO: Show error in UI
		return
	}

	err = a.resourceTypesView.LoadResourceTypes(ctx, resourceTypes, subscriptionID, subscriptionName, resourceGroupName)
	if err == nil {
		a.updateFooterForTableView(a.resourceTypesView.TableView)
	}
	a.updateLayout()
	a.SetFocus(a.resourceTypesView)
}

// navigateToResourceType navigates to a resource type filtered view
func (a *App) navigateToResourceType(resourceType string) {
	a.navState.NavigateToResourceType(resourceType)

	// Load resources filtered by type
	ctx := context.Background()
	subscriptionID := a.navState.SelectedSubscriptionID
	subscriptionName := a.navState.SelectedSubscriptionName
	resourceGroupName := a.navState.SelectedResourceGroupName

	var resources []*models.Resource
	var err error
	if resourceGroupName != "" {
		// Filter by resource type within the resource group
		resources, err = a.azureClient.ListResourcesByResourceGroup(ctx, subscriptionID, resourceGroupName, resourceType)
	} else {
		// Filter by resource type across the subscription
		resources, err = a.azureClient.ListResources(ctx, subscriptionID, resourceType)
	}

	if err != nil {
		// TODO: Show error in UI
		return
	}

	// Update title (empty since breadcrumb shows navigation path)
	a.resourcesView.SetTitle("")

	err = a.resourcesView.LoadResources(ctx, resources, subscriptionID, subscriptionName, resourceGroupName)
	if err == nil {
		a.updateFooterForTableView(a.resourcesView.TableView)
	}
	a.updateLayout()
	a.SetFocus(a.resourcesView)
}

// navigateBackToResourceTypes returns from resource type view to resource types view
func (a *App) navigateBackToResourceTypes() {
	// Restore the resource types view for the resource group
	resourceGroupName := a.navState.SelectedResourceGroupName
	a.navigateToResourceTypes(resourceGroupName)
}

// navigateBackToResourceType returns from storage explorer to resource type view
func (a *App) navigateBackToResourceType() {
	// Go back to the resource type view (filtered resources)
	resourceType := a.navState.SelectedResourceType
	a.navigateToResourceType(resourceType)
}

// navigateBackToKeyVaultExplorer returns from Key Vault item views to Key Vault explorer
func (a *App) navigateBackToKeyVaultExplorer() {
	keyVaultName := a.navState.SelectedKeyVault
	vaultURL := a.navState.SelectedKeyVaultURL
	
	// Navigate back to Key Vault explorer view
	a.navState.CurrentView = navigation.ViewKeyVaultExplorer
	
	// Reload Key Vault explorer
	ctx := context.Background()
	err := a.keyVaultExplorerView.LoadKeyVault(ctx, keyVaultName, vaultURL)
	if err != nil {
		// TODO: Show error in UI
		return
	}
	
	a.updateLayout()
	a.SetFocus(a.keyVaultExplorerView)
}

// navigateBackToResourceGroups returns from resources view to resource groups view
func (a *App) navigateBackToResourceGroups() {
	subscriptionID := a.navState.SelectedSubscriptionID
	subscriptionName := a.navState.SelectedSubscriptionName
	a.navigateToResourceGroups(subscriptionID, subscriptionName)
}

// showResourceDetails shows the details view for a resource
func (a *App) showResourceDetails(resource *models.Resource) {
	a.navState.NavigateToDetails()
	subscriptionID := a.resourcesView.GetSubscriptionID()
	a.detailsView.ShowResourceDetails(resource, subscriptionID)
	a.updateLayout()
	a.SetFocus(a.detailsView)
}

// navigateToStorageExplorer navigates to the storage explorer view for a storage account
func (a *App) navigateToStorageExplorer(resource *models.Resource) {
	storageAccountName := resource.Name
	a.navState.NavigateToStorageExplorer(storageAccountName)

	// Load containers
	ctx := context.Background()
	subscriptionID := a.navState.SelectedSubscriptionID
	resourceGroupName := resource.ResourceGroup
	containers, err := a.azureClient.ListContainers(ctx, subscriptionID, resourceGroupName, storageAccountName)
	if err != nil {
		// TODO: Show error in UI
		return
	}

	err = a.storageExplorerView.LoadContainers(ctx, containers, storageAccountName)
	if err == nil {
		a.updateFooterForTableView(a.storageExplorerView.TableView)
	}
	a.updateLayout()
	a.SetFocus(a.storageExplorerView)
}

// navigateToBlobs navigates to the blobs view for a container
func (a *App) navigateToBlobs(containerName string) {
	a.navState.NavigateToBlobs(containerName)

	// Load blobs at root level
	a.loadBlobsForCurrentPath()
}

// loadBlobsForCurrentPath loads blobs for the current path prefix
func (a *App) loadBlobsForCurrentPath() {
	ctx := context.Background()
	subscriptionID := a.navState.SelectedSubscriptionID
	resourceGroupName := a.navState.SelectedResourceGroupName
	storageAccountName := a.navState.SelectedStorageAccount
	containerName := a.navState.SelectedContainer
	pathPrefix := a.navState.BlobPathPrefix

	blobs, err := a.azureClient.ListBlobs(ctx, subscriptionID, resourceGroupName, storageAccountName, containerName, pathPrefix)
	if err != nil {
		// TODO: Show error in UI
		return
	}

	err = a.blobsView.LoadBlobs(ctx, blobs, containerName, storageAccountName, pathPrefix)
	if err == nil {
		a.updateFooterForTableView(a.blobsView.TableView)
	}
	a.updateLayout()
	a.SetFocus(a.blobsView)
}

// navigateIntoBlobFolder navigates into a blob folder
func (a *App) navigateIntoBlobFolder(folderPath string) {
	a.navState.NavigateIntoBlobFolder(folderPath)
	a.loadBlobsForCurrentPath()
}

// navigateBackFromBlobs returns from blobs view to storage explorer or parent folder
func (a *App) navigateBackFromBlobs() {
	// Check if we're in a subfolder
	if a.navState.BlobPathPrefix != "" {
		// Go back to parent folder
		a.navState.NavigateBackFromBlobFolder()
		a.loadBlobsForCurrentPath()
		return
	}

	// Go back to storage explorer
	storageAccountName := a.navState.SelectedStorageAccount
	a.navState.NavigateBackFromBlobs()

	// Reload containers
	ctx := context.Background()
	subscriptionID := a.navState.SelectedSubscriptionID
	resourceGroupName := a.navState.SelectedResourceGroupName
	containers, err := a.azureClient.ListContainers(ctx, subscriptionID, resourceGroupName, storageAccountName)
	if err != nil {
		// TODO: Show error in UI
		return
	}

	err = a.storageExplorerView.LoadContainers(ctx, containers, storageAccountName)
	if err == nil {
		a.updateFooterForTableView(a.storageExplorerView.TableView)
	}
	a.updateLayout()
	a.SetFocus(a.storageExplorerView)
}

// showContainerDetails shows the details view for a container
func (a *App) showContainerDetails(container *models.Container) {
	a.navState.NavigateToDetails()
	storageAccountName := a.navState.SelectedStorageAccount
	a.detailsView.ShowContainerDetails(container, storageAccountName)
	a.updateLayout()
	a.SetFocus(a.detailsView)
}

// showBlobDetails shows the details view for a blob
func (a *App) showBlobDetails(blob *models.Blob) {
	a.navState.NavigateToDetails()
	subscriptionID := a.navState.SelectedSubscriptionID
	resourceGroupName := a.navState.SelectedResourceGroupName
	storageAccountName := a.navState.SelectedStorageAccount
	containerName := a.navState.SelectedContainer

	// Get full blob details
	ctx := context.Background()
	fullBlob, err := a.azureClient.GetBlobDetails(ctx, subscriptionID, resourceGroupName, storageAccountName, containerName, blob.Name)
	if err != nil {
		// TODO: Show error in UI
		return
	}

	a.detailsView.ShowBlobDetails(fullBlob, storageAccountName, containerName)
	a.updateLayout()
	a.SetFocus(a.detailsView)
}

// navigateToKeyVaultExplorer navigates to the Key Vault explorer view for a Key Vault
func (a *App) navigateToKeyVaultExplorer(resource *models.Resource) {
	keyVaultName := resource.Name
	
	// Extract vault URL from resource properties
	vaultURL := ""
	if resource.Properties != nil {
		if vaultURIVal, ok := resource.Properties["vaultUri"]; ok {
			if vaultURIStr, ok := vaultURIVal.(string); ok {
				vaultURL = vaultURIStr
			}
		}
	}
	
	// If not in properties, construct it
	if vaultURL == "" {
		vaultURL = fmt.Sprintf("https://%s.vault.azure.net/", keyVaultName)
	}
	
	a.navState.NavigateToKeyVaultExplorer(keyVaultName, vaultURL)
	
	// Load Key Vault explorer
	ctx := context.Background()
	err := a.keyVaultExplorerView.LoadKeyVault(ctx, keyVaultName, vaultURL)
	if err != nil {
		// TODO: Show error in UI
		return
	}
	
	a.updateLayout()
	a.SetFocus(a.keyVaultExplorerView)
}

// navigateToKeyVaultItemType navigates to the selected Key Vault item type (secrets, keys, or certificates)
func (a *App) navigateToKeyVaultItemType(itemType string) {
	ctx := context.Background()
	vaultURL := a.navState.SelectedKeyVaultURL
	keyVaultName := a.navState.SelectedKeyVault
	
	switch itemType {
	case "secrets":
		a.navState.NavigateToKeyVaultSecrets()
		
		// Load secrets
		secrets, err := a.azureClient.ListSecrets(ctx, vaultURL)
		if err != nil {
			// TODO: Show error in UI
			return
		}
		
		err = a.keyVaultSecretsView.LoadSecrets(ctx, secrets, keyVaultName, vaultURL)
		if err == nil {
			a.updateFooterForTableView(a.keyVaultSecretsView.TableView)
		}
		a.updateLayout()
		a.SetFocus(a.keyVaultSecretsView)
		
	case "keys":
		a.navState.NavigateToKeyVaultKeys()
		
		// Load keys
		keys, err := a.azureClient.ListKeys(ctx, vaultURL)
		if err != nil {
			// TODO: Show error in UI
			return
		}
		
		err = a.keyVaultKeysView.LoadKeys(ctx, keys, keyVaultName, vaultURL)
		if err == nil {
			a.updateFooterForTableView(a.keyVaultKeysView.TableView)
		}
		a.updateLayout()
		a.SetFocus(a.keyVaultKeysView)
		
	case "certificates":
		a.navState.NavigateToKeyVaultCertificates()
		
		// Load certificates
		certificates, err := a.azureClient.ListCertificates(ctx, vaultURL)
		if err != nil {
			// TODO: Show error in UI
			return
		}
		
		err = a.keyVaultCertificatesView.LoadCertificates(ctx, certificates, keyVaultName, vaultURL)
		if err == nil {
			a.updateFooterForTableView(a.keyVaultCertificatesView.TableView)
		}
		a.updateLayout()
		a.SetFocus(a.keyVaultCertificatesView)
	}
}

// showSecretDetails shows the details view for a secret
func (a *App) showSecretDetails(secret *models.Secret) {
	a.navState.NavigateToDetails()
	keyVaultName := a.navState.SelectedKeyVault
	a.detailsView.ShowSecretDetails(secret, keyVaultName)
	a.updateLayout()
	a.SetFocus(a.detailsView)
}

// viewSecretValue shows the secret value with a confirmation dialog
func (a *App) viewSecretValue(secret *models.Secret) {
	// Create a modal for confirmation
	modal := tview.NewModal().
		SetText(fmt.Sprintf("Are you sure you want to view the value of secret '%s'?\n\n⚠️ This will display sensitive information on screen.", secret.Name)).
		AddButtons([]string{"View", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "View" {
				// Fetch and show the secret value
				ctx := context.Background()
				vaultURL := a.navState.SelectedKeyVaultURL
				value, err := a.azureClient.GetSecretValue(ctx, vaultURL, secret.Name)
				if err != nil {
					// TODO: Show error in UI
					a.SetRoot(a.mainFlex, true)
					return
				}
				
				// Create a modal to display the value
				valueModal := tview.NewModal().
					SetText(fmt.Sprintf("Secret: %s\n\nValue:\n%s\n\nPress any key to close.", secret.Name, value)).
					AddButtons([]string{"Close"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						a.SetRoot(a.mainFlex, true)
						a.SetFocus(a.keyVaultSecretsView)
					})
				a.SetRoot(valueModal, true)
			} else {
				a.SetRoot(a.mainFlex, true)
				a.SetFocus(a.keyVaultSecretsView)
			}
		})
	
	a.SetRoot(modal, true)
}

// showKeyDetails shows the details view for a key
func (a *App) showKeyDetails(key *models.Key) {
	a.navState.NavigateToDetails()
	keyVaultName := a.navState.SelectedKeyVault
	
	// Get full key details
	ctx := context.Background()
	vaultURL := a.navState.SelectedKeyVaultURL
	fullKey, err := a.azureClient.GetKeyDetails(ctx, vaultURL, key.Name)
	if err != nil {
		// TODO: Show error in UI
		return
	}
	
	a.detailsView.ShowKeyDetails(fullKey, keyVaultName)
	a.updateLayout()
	a.SetFocus(a.detailsView)
}

// showCertificateDetails shows the details view for a certificate
func (a *App) showCertificateDetails(cert *models.Certificate) {
	a.navState.NavigateToDetails()
	keyVaultName := a.navState.SelectedKeyVault
	
	// Get full certificate details
	ctx := context.Background()
	vaultURL := a.navState.SelectedKeyVaultURL
	fullCert, err := a.azureClient.GetCertificateDetails(ctx, vaultURL, cert.Name)
	if err != nil {
		// TODO: Show error in UI
		return
	}
	
	a.detailsView.ShowCertificateDetails(fullCert, keyVaultName)
	a.updateLayout()
	a.SetFocus(a.detailsView)
}

// navigateToMenu navigates to the menu view
func (a *App) navigateToMenu() {
	a.navState.NavigateToMenu()

	// Load menu with current context
	ctx := context.Background()
	subscriptionID := a.navState.SelectedSubscriptionID
	subscriptionName := a.navState.SelectedSubscriptionName
	resourceGroupName := a.navState.SelectedResourceGroupName

	err := a.menuView.LoadResourceTypes(ctx, subscriptionID, subscriptionName, resourceGroupName)
	if err == nil {
		a.updateFooterForTableView(a.menuView.TableView)
	}
	a.updateLayout()
	a.SetFocus(a.menuView)
}

// navigateBackFromMenu returns from menu view
func (a *App) navigateBackFromMenu() {
	a.navState.NavigateBackFromMenu()
	a.updateLayout()
	switch a.navState.CurrentView {
	case navigation.ViewSubscriptions:
		a.SetFocus(a.subscriptionsView)
	case navigation.ViewResourceGroups:
		a.SetFocus(a.resourceGroupsView)
	case navigation.ViewResourceTypes:
		a.SetFocus(a.resourceTypesView)
	case navigation.ViewResources, navigation.ViewResourceType:
		a.SetFocus(a.resourcesView)
	}
}

// navigateToResourceTypeFromMenu navigates to a resource type list from the menu
func (a *App) navigateToResourceTypeFromMenu(resourceType string) {
	// Navigate to resource type view
	a.navState.NavigateToResourceType(resourceType)

	// Load resources filtered by type
	ctx := context.Background()
	subscriptionID := a.navState.SelectedSubscriptionID
	subscriptionName := a.navState.SelectedSubscriptionName
	resourceGroupName := a.navState.SelectedResourceGroupName

	var resources []*models.Resource
	var err error
	if resourceGroupName != "" {
		// Filter by resource type within the resource group
		resources, err = a.azureClient.ListResourcesByResourceGroup(ctx, subscriptionID, resourceGroupName, resourceType)
	} else {
		// Filter by resource type across the subscription
		resources, err = a.azureClient.ListResources(ctx, subscriptionID, resourceType)
	}

	if err != nil {
		// TODO: Show error in UI
		return
	}

	// Update title (empty since breadcrumb shows navigation path)
	a.resourcesView.SetTitle("")

	err = a.resourcesView.LoadResources(ctx, resources, subscriptionID, subscriptionName, resourceGroupName)
	if err == nil {
		a.updateFooterForTableView(a.resourcesView.TableView)
	}
	a.updateLayout()
	a.SetFocus(a.resourcesView)
}

// applyFilter applies a filter to the current table view
func (a *App) applyFilter(filterText string) {
	if a.navState.InDetailsView {
		return
	}

	switch a.navState.CurrentView {
	case navigation.ViewSubscriptions:
		a.subscriptionsView.SetFilter(filterText)
		a.updateFooterForTableView(a.subscriptionsView.TableView)
	case navigation.ViewResourceGroups:
		a.resourceGroupsView.SetFilter(filterText)
		a.updateFooterForTableView(a.resourceGroupsView.TableView)
	case navigation.ViewResourceTypes:
		a.resourceTypesView.SetFilter(filterText)
		a.updateFooterForTableView(a.resourceTypesView.TableView)
	case navigation.ViewResources, navigation.ViewResourceType:
		a.resourcesView.SetFilter(filterText)
		a.updateFooterForTableView(a.resourcesView.TableView)
	case navigation.ViewStorageExplorer:
		a.storageExplorerView.SetFilter(filterText)
		a.updateFooterForTableView(a.storageExplorerView.TableView)
	case navigation.ViewBlobs:
		a.blobsView.SetFilter(filterText)
		a.updateFooterForTableView(a.blobsView.TableView)
	case navigation.ViewMenu:
		a.menuView.SetFilter(filterText)
		a.updateFooterForTableView(a.menuView.TableView)
	}

	a.updateLayout()
	a.SetFocus(a.currentView)
}

// clearFilter clears the filter from the current table view
func (a *App) clearFilter() {
	if a.navState.InDetailsView {
		return
	}

	switch a.navState.CurrentView {
	case navigation.ViewSubscriptions:
		a.subscriptionsView.ClearFilter()
		a.updateFooterForTableView(a.subscriptionsView.TableView)
	case navigation.ViewResourceGroups:
		a.resourceGroupsView.ClearFilter()
		a.updateFooterForTableView(a.resourceGroupsView.TableView)
	case navigation.ViewResourceTypes:
		a.resourceTypesView.ClearFilter()
		a.updateFooterForTableView(a.resourceTypesView.TableView)
	case navigation.ViewResources, navigation.ViewResourceType:
		a.resourcesView.ClearFilter()
		a.updateFooterForTableView(a.resourcesView.TableView)
	case navigation.ViewStorageExplorer:
		a.storageExplorerView.ClearFilter()
		a.updateFooterForTableView(a.storageExplorerView.TableView)
	case navigation.ViewBlobs:
		a.blobsView.ClearFilter()
		a.updateFooterForTableView(a.blobsView.TableView)
	case navigation.ViewMenu:
		a.menuView.ClearFilter()
		a.updateFooterForTableView(a.menuView.TableView)
	}
}
