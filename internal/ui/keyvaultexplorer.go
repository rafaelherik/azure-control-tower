package ui

import (
	"context"
	"fmt"
	"time"

	"azure-control-tower/internal/models"

	"github.com/rivo/tview"
)

// KeyVaultExplorerView displays the main navigation for a Key Vault
type KeyVaultExplorerView struct {
	*TableView
	keyVaultName string
	vaultURL     string
	onSelect     func(itemType string)
}

// NewKeyVaultExplorerView creates a new Key Vault explorer view
func NewKeyVaultExplorerView() *KeyVaultExplorerView {
	kve := &KeyVaultExplorerView{}

	// Create table configuration
	config := &TableConfig{
		Title: "",
		Columns: []ColumnConfig{
			{Name: "Item Type", Align: tview.AlignLeft},
			{Name: "Description", Align: tview.AlignLeft},
		},
		OnSelect: func(rowIndex int, data interface{}) {
			// Enter key on an item - navigate to that item type
			if itemData, ok := data.(string); ok && kve.onSelect != nil {
				kve.onSelect(itemData)
			}
		},
		GetCellValue: func(data interface{}, columnIndex int) string {
			itemType, ok := data.(string)
			if !ok {
				return ""
			}
			switch columnIndex {
			case 0:
				switch itemType {
				case "secrets":
					return "🔐 Secrets"
				case "keys":
					return "🔑 Keys"
				case "certificates":
					return "📜 Certificates"
				default:
					return itemType
				}
			case 1:
				switch itemType {
				case "secrets":
					return "Manage secret values and configurations"
				case "keys":
					return "Manage cryptographic keys"
				case "certificates":
					return "Manage SSL/TLS certificates"
				default:
					return ""
				}
			default:
				return ""
			}
		},
	}

	kve.TableView = NewTableView(config)
	return kve
}

// LoadKeyVault loads the Key Vault explorer view
func (kve *KeyVaultExplorerView) LoadKeyVault(ctx context.Context, keyVaultName, vaultURL string) error {
	kve.keyVaultName = keyVaultName
	kve.vaultURL = vaultURL

	// Load the three main item types
	data := []interface{}{
		"secrets",
		"keys",
		"certificates",
	}

	kve.LoadData(data)
	return nil
}

// SetOnSelect sets the callback for when an item type is selected (Enter key)
func (kve *KeyVaultExplorerView) SetOnSelect(callback func(string)) {
	kve.onSelect = callback
}

// GetKeyVaultName returns the current Key Vault name
func (kve *KeyVaultExplorerView) GetKeyVaultName() string {
	return kve.keyVaultName
}

// GetVaultURL returns the current vault URL
func (kve *KeyVaultExplorerView) GetVaultURL() string {
	return kve.vaultURL
}

// SecretRowData wraps Secret with context info for display
type SecretRowData struct {
	Secret *models.Secret
}

// KeyVaultSecretsView displays secrets in a Key Vault
type KeyVaultSecretsView struct {
	*TableView
	secrets      []*models.Secret
	keyVaultName string
	vaultURL     string
	onShowDetails func(secret *models.Secret)
	onViewValue   func(secret *models.Secret)
}

// NewKeyVaultSecretsView creates a new Key Vault secrets view
func NewKeyVaultSecretsView() *KeyVaultSecretsView {
	ksv := &KeyVaultSecretsView{}

	// Create table configuration
	config := &TableConfig{
		Title: "",
		Columns: []ColumnConfig{
			{Name: "Name", Align: tview.AlignLeft},
			{Name: "Enabled", Align: tview.AlignCenter},
			{Name: "Content Type", Align: tview.AlignLeft},
			{Name: "Updated", Align: tview.AlignLeft},
		},
		RowActions: []RowAction{
			{
				Rune:  'v',
				Label: "View Value",
				Callback: func(rowIndex int, data interface{}) bool {
					if rowData, ok := data.(*SecretRowData); ok && ksv.onViewValue != nil {
						ksv.onViewValue(rowData.Secret)
						return true
					}
					return false
				},
			},
			{
				Rune:  'd',
				Label: "Details",
				Callback: func(rowIndex int, data interface{}) bool {
					if rowData, ok := data.(*SecretRowData); ok && ksv.onShowDetails != nil {
						ksv.onShowDetails(rowData.Secret)
						return true
					}
					return false
				},
			},
		},
		OnSelect: func(rowIndex int, data interface{}) {
			// Enter key on a secret - show details
			if rowData, ok := data.(*SecretRowData); ok && ksv.onShowDetails != nil {
				ksv.onShowDetails(rowData.Secret)
			}
		},
		GetCellValue: func(data interface{}, columnIndex int) string {
			rowData, ok := data.(*SecretRowData)
			if !ok {
				return ""
			}
			switch columnIndex {
			case 0:
				return rowData.Secret.Name
			case 1:
				if rowData.Secret.Enabled {
					return "✓"
				}
				return "✗"
			case 2:
				return rowData.Secret.ContentType
			case 3:
				if rowData.Secret.Updated != nil {
					return rowData.Secret.Updated.Format("2006-01-02 15:04:05")
				}
				return "-"
			default:
				return ""
			}
		},
	}

	ksv.TableView = NewTableView(config)
	return ksv
}

// LoadSecrets loads secrets into the view
func (ksv *KeyVaultSecretsView) LoadSecrets(ctx context.Context, secrets []*models.Secret, keyVaultName, vaultURL string) error {
	ksv.secrets = secrets
	ksv.keyVaultName = keyVaultName
	ksv.vaultURL = vaultURL

	// Convert to interface{} slice with row data
	data := make([]interface{}, len(secrets))
	for i, secret := range secrets {
		data[i] = &SecretRowData{
			Secret: secret,
		}
	}

	ksv.LoadData(data)
	return nil
}

// SetOnShowDetails sets the callback for when details are requested
func (ksv *KeyVaultSecretsView) SetOnShowDetails(callback func(*models.Secret)) {
	ksv.onShowDetails = callback
}

// SetOnViewValue sets the callback for when viewing a secret value
func (ksv *KeyVaultSecretsView) SetOnViewValue(callback func(*models.Secret)) {
	ksv.onViewValue = callback
}

// GetKeyVaultName returns the current Key Vault name
func (ksv *KeyVaultSecretsView) GetKeyVaultName() string {
	return ksv.keyVaultName
}

// GetVaultURL returns the current vault URL
func (ksv *KeyVaultSecretsView) GetVaultURL() string {
	return ksv.vaultURL
}

// KeyRowData wraps Key with context info for display
type KeyRowData struct {
	Key *models.Key
}

// KeyVaultKeysView displays keys in a Key Vault
type KeyVaultKeysView struct {
	*TableView
	keys         []*models.Key
	keyVaultName string
	vaultURL     string
	onShowDetails func(key *models.Key)
}

// NewKeyVaultKeysView creates a new Key Vault keys view
func NewKeyVaultKeysView() *KeyVaultKeysView {
	kkv := &KeyVaultKeysView{}

	// Create table configuration
	config := &TableConfig{
		Title: "",
		Columns: []ColumnConfig{
			{Name: "Name", Align: tview.AlignLeft},
			{Name: "Type", Align: tview.AlignLeft},
			{Name: "Enabled", Align: tview.AlignCenter},
			{Name: "Updated", Align: tview.AlignLeft},
		},
		RowActions: []RowAction{
			{
				Rune:  'd',
				Label: "Details",
				Callback: func(rowIndex int, data interface{}) bool {
					if rowData, ok := data.(*KeyRowData); ok && kkv.onShowDetails != nil {
						kkv.onShowDetails(rowData.Key)
						return true
					}
					return false
				},
			},
		},
		OnSelect: func(rowIndex int, data interface{}) {
			// Enter key on a key - show details
			if rowData, ok := data.(*KeyRowData); ok && kkv.onShowDetails != nil {
				kkv.onShowDetails(rowData.Key)
			}
		},
		GetCellValue: func(data interface{}, columnIndex int) string {
			rowData, ok := data.(*KeyRowData)
			if !ok {
				return ""
			}
			switch columnIndex {
			case 0:
				return rowData.Key.Name
			case 1:
				return rowData.Key.KeyType
			case 2:
				if rowData.Key.Enabled {
					return "✓"
				}
				return "✗"
			case 3:
				if rowData.Key.Updated != nil {
					return rowData.Key.Updated.Format("2006-01-02 15:04:05")
				}
				return "-"
			default:
				return ""
			}
		},
	}

	kkv.TableView = NewTableView(config)
	return kkv
}

// LoadKeys loads keys into the view
func (kkv *KeyVaultKeysView) LoadKeys(ctx context.Context, keys []*models.Key, keyVaultName, vaultURL string) error {
	kkv.keys = keys
	kkv.keyVaultName = keyVaultName
	kkv.vaultURL = vaultURL

	// Convert to interface{} slice with row data
	data := make([]interface{}, len(keys))
	for i, key := range keys {
		data[i] = &KeyRowData{
			Key: key,
		}
	}

	kkv.LoadData(data)
	return nil
}

// SetOnShowDetails sets the callback for when details are requested
func (kkv *KeyVaultKeysView) SetOnShowDetails(callback func(*models.Key)) {
	kkv.onShowDetails = callback
}

// GetKeyVaultName returns the current Key Vault name
func (kkv *KeyVaultKeysView) GetKeyVaultName() string {
	return kkv.keyVaultName
}

// GetVaultURL returns the current vault URL
func (kkv *KeyVaultKeysView) GetVaultURL() string {
	return kkv.vaultURL
}

// CertificateRowData wraps Certificate with context info for display
type CertificateRowData struct {
	Certificate *models.Certificate
}

// KeyVaultCertificatesView displays certificates in a Key Vault
type KeyVaultCertificatesView struct {
	*TableView
	certificates []*models.Certificate
	keyVaultName string
	vaultURL     string
	onShowDetails func(cert *models.Certificate)
}

// NewKeyVaultCertificatesView creates a new Key Vault certificates view
func NewKeyVaultCertificatesView() *KeyVaultCertificatesView {
	kcv := &KeyVaultCertificatesView{}

	// Create table configuration
	config := &TableConfig{
		Title: "",
		Columns: []ColumnConfig{
			{Name: "Name", Align: tview.AlignLeft},
			{Name: "Enabled", Align: tview.AlignCenter},
			{Name: "Expires", Align: tview.AlignLeft},
			{Name: "Updated", Align: tview.AlignLeft},
		},
		RowActions: []RowAction{
			{
				Rune:  'd',
				Label: "Details",
				Callback: func(rowIndex int, data interface{}) bool {
					if rowData, ok := data.(*CertificateRowData); ok && kcv.onShowDetails != nil {
						kcv.onShowDetails(rowData.Certificate)
						return true
					}
					return false
				},
			},
		},
		OnSelect: func(rowIndex int, data interface{}) {
			// Enter key on a certificate - show details
			if rowData, ok := data.(*CertificateRowData); ok && kcv.onShowDetails != nil {
				kcv.onShowDetails(rowData.Certificate)
			}
		},
		GetCellValue: func(data interface{}, columnIndex int) string {
			rowData, ok := data.(*CertificateRowData)
			if !ok {
				return ""
			}
			switch columnIndex {
			case 0:
				return rowData.Certificate.Name
			case 1:
				if rowData.Certificate.Enabled {
					return "✓"
				}
				return "✗"
			case 2:
				if rowData.Certificate.Expires != nil {
					expires := rowData.Certificate.Expires.Format("2006-01-02")
					// Add warning if expired or expiring soon
					if rowData.Certificate.Expires.Before(time.Now()) {
						return fmt.Sprintf("⚠️ %s (EXPIRED)", expires)
					}
					return expires
				}
				return "-"
			case 3:
				if rowData.Certificate.Updated != nil {
					return rowData.Certificate.Updated.Format("2006-01-02 15:04:05")
				}
				return "-"
			default:
				return ""
			}
		},
	}

	kcv.TableView = NewTableView(config)
	return kcv
}

// LoadCertificates loads certificates into the view
func (kcv *KeyVaultCertificatesView) LoadCertificates(ctx context.Context, certificates []*models.Certificate, keyVaultName, vaultURL string) error {
	kcv.certificates = certificates
	kcv.keyVaultName = keyVaultName
	kcv.vaultURL = vaultURL

	// Convert to interface{} slice with row data
	data := make([]interface{}, len(certificates))
	for i, cert := range certificates {
		data[i] = &CertificateRowData{
			Certificate: cert,
		}
	}

	kcv.LoadData(data)
	return nil
}

// SetOnShowDetails sets the callback for when details are requested
func (kcv *KeyVaultCertificatesView) SetOnShowDetails(callback func(*models.Certificate)) {
	kcv.onShowDetails = callback
}

// GetKeyVaultName returns the current Key Vault name
func (kcv *KeyVaultCertificatesView) GetKeyVaultName() string {
	return kcv.keyVaultName
}

// GetVaultURL returns the current vault URL
func (kcv *KeyVaultCertificatesView) GetVaultURL() string {
	return kcv.vaultURL
}
