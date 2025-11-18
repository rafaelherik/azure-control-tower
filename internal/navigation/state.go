package navigation

// ViewType represents the current view type
type ViewType int

const (
	ViewSubscriptions ViewType = iota
	ViewResourceGroups
	ViewResourceTypes
	ViewResources
	ViewResourceType
	ViewDetails
	ViewStorageExplorer
	ViewBlobs
	ViewKeyVaultExplorer
	ViewKeyVaultSecrets
	ViewKeyVaultKeys
	ViewKeyVaultCertificates
	ViewMenu
)

// State manages navigation state
type State struct {
	CurrentView               ViewType
	SelectedSubscriptionID    string
	SelectedSubscriptionName  string
	SelectedResourceGroupName string
	SelectedResourceType      string
	InDetailsView             bool
	SelectedStorageAccount    string
	SelectedContainer         string
	SelectedBlob              string
	BlobPathPrefix            string // Current folder path prefix in blob view
	SelectedKeyVault          string
	SelectedKeyVaultURL       string
}

// NewState creates a new navigation state
func NewState() *State {
	return &State{
		CurrentView: ViewSubscriptions,
	}
}

// NavigateToSubscriptions navigates to the subscriptions view
func (s *State) NavigateToSubscriptions() {
	s.CurrentView = ViewSubscriptions
	s.SelectedSubscriptionID = ""
	s.SelectedSubscriptionName = ""
	s.SelectedResourceGroupName = ""
	s.SelectedResourceType = ""
	s.InDetailsView = false
}

// NavigateToResourceGroups navigates to the resource groups view for a subscription
func (s *State) NavigateToResourceGroups(subscriptionID, subscriptionName string) {
	s.CurrentView = ViewResourceGroups
	s.SelectedSubscriptionID = subscriptionID
	s.SelectedSubscriptionName = subscriptionName
	s.SelectedResourceGroupName = ""
	s.SelectedResourceType = ""
	s.InDetailsView = false
}

// NavigateToResourceTypes navigates to the resource types summary view for a resource group
func (s *State) NavigateToResourceTypes(resourceGroupName string) {
	s.CurrentView = ViewResourceTypes
	s.SelectedResourceGroupName = resourceGroupName
	s.SelectedResourceType = ""
	s.InDetailsView = false
}

// NavigateToResources navigates to the resources view for a resource group
func (s *State) NavigateToResources(resourceGroupName string) {
	s.CurrentView = ViewResources
	s.SelectedResourceGroupName = resourceGroupName
	s.SelectedResourceType = ""
	s.InDetailsView = false
}

// NavigateToResourceType navigates to the resource type filtered view
func (s *State) NavigateToResourceType(resourceType string) {
	s.CurrentView = ViewResourceType
	s.SelectedResourceType = resourceType
	s.InDetailsView = false
}

// NavigateToDetails navigates to the details view
func (s *State) NavigateToDetails() {
	s.InDetailsView = true
}

// NavigateBackFromDetails returns from details view to previous view
func (s *State) NavigateBackFromDetails() {
	s.InDetailsView = false
}

// NavigateToStorageExplorer navigates to the storage explorer view
func (s *State) NavigateToStorageExplorer(storageAccountName string) {
	s.CurrentView = ViewStorageExplorer
	s.SelectedStorageAccount = storageAccountName
	s.SelectedContainer = ""
	s.SelectedBlob = ""
	s.InDetailsView = false
}

// NavigateToBlobs navigates to the blobs view for a container
func (s *State) NavigateToBlobs(containerName string) {
	s.CurrentView = ViewBlobs
	s.SelectedContainer = containerName
	s.SelectedBlob = ""
	s.BlobPathPrefix = ""
	s.InDetailsView = false
}

// NavigateBackFromBlobs returns from blobs view to storage explorer
func (s *State) NavigateBackFromBlobs() {
	s.CurrentView = ViewStorageExplorer
	s.SelectedContainer = ""
	s.SelectedBlob = ""
	s.BlobPathPrefix = ""
}

// NavigateIntoBlobFolder navigates into a blob folder
func (s *State) NavigateIntoBlobFolder(folderPath string) {
	s.BlobPathPrefix = folderPath
	s.SelectedBlob = ""
}

// NavigateBackFromBlobFolder returns from a blob folder to parent folder
func (s *State) NavigateBackFromBlobFolder() {
	if s.BlobPathPrefix == "" {
		// Already at root, go back to storage explorer
		s.NavigateBackFromBlobs()
		return
	}
	
	// Remove last path segment
	lastSlash := -1
	for i := len(s.BlobPathPrefix) - 2; i >= 0; i-- {
		if s.BlobPathPrefix[i] == '/' {
			lastSlash = i
			break
		}
	}
	
	if lastSlash >= 0 {
		s.BlobPathPrefix = s.BlobPathPrefix[:lastSlash+1]
	} else {
		s.BlobPathPrefix = ""
	}
	s.SelectedBlob = ""
}

// NavigateToMenu navigates to the menu view
func (s *State) NavigateToMenu() {
	s.CurrentView = ViewMenu
	s.InDetailsView = false
}

// NavigateBackFromMenu returns from menu view to previous view
func (s *State) NavigateBackFromMenu() {
	// Restore previous view - for now, go back to subscriptions
	// In a more sophisticated implementation, we could track the previous view
	s.NavigateToSubscriptions()
}

// NavigateToKeyVaultExplorer navigates to the Key Vault explorer view
func (s *State) NavigateToKeyVaultExplorer(keyVaultName, vaultURL string) {
	s.CurrentView = ViewKeyVaultExplorer
	s.SelectedKeyVault = keyVaultName
	s.SelectedKeyVaultURL = vaultURL
	s.InDetailsView = false
}

// NavigateToKeyVaultSecrets navigates to the secrets view for a Key Vault
func (s *State) NavigateToKeyVaultSecrets() {
	s.CurrentView = ViewKeyVaultSecrets
	s.InDetailsView = false
}

// NavigateToKeyVaultKeys navigates to the keys view for a Key Vault
func (s *State) NavigateToKeyVaultKeys() {
	s.CurrentView = ViewKeyVaultKeys
	s.InDetailsView = false
}

// NavigateToKeyVaultCertificates navigates to the certificates view for a Key Vault
func (s *State) NavigateToKeyVaultCertificates() {
	s.CurrentView = ViewKeyVaultCertificates
	s.InDetailsView = false
}

// NavigateBackFromKeyVaultSecrets returns from secrets view to Key Vault explorer
func (s *State) NavigateBackFromKeyVaultSecrets() {
	s.CurrentView = ViewKeyVaultExplorer
}

// NavigateBackFromKeyVaultKeys returns from keys view to Key Vault explorer
func (s *State) NavigateBackFromKeyVaultKeys() {
	s.CurrentView = ViewKeyVaultExplorer
}

// NavigateBackFromKeyVaultCertificates returns from certificates view to Key Vault explorer
func (s *State) NavigateBackFromKeyVaultCertificates() {
	s.CurrentView = ViewKeyVaultExplorer
}
