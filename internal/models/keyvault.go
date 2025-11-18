package models

import "time"

// KeyVault represents an Azure Key Vault instance
type KeyVault struct {
	ID                string
	Name              string
	Location          string
	ResourceGroup     string
	VaultURI          string
	TenantID          string
	SKU               string
	EnabledForDeploy  bool
	EnabledForDisk    bool
	EnabledForTemplate bool
	Tags              map[string]*string
	Properties        map[string]interface{}
}

// Secret represents a Key Vault secret
type Secret struct {
	Name        string
	Value       string // Only populated when explicitly retrieved
	Enabled     bool
	Created     *time.Time
	Updated     *time.Time
	Expires     *time.Time
	NotBefore   *time.Time
	Version     string
	ContentType string
	Tags        map[string]string
}

// Key represents a Key Vault key
type Key struct {
	Name      string
	KeyType   string // RSA, EC, etc.
	Enabled   bool
	Created   *time.Time
	Updated   *time.Time
	Expires   *time.Time
	NotBefore *time.Time
	Version   string
	Tags      map[string]string
}

// Certificate represents a Key Vault certificate
type Certificate struct {
	Name        string
	Enabled     bool
	Created     *time.Time
	Updated     *time.Time
	Expires     *time.Time
	NotBefore   *time.Time
	Version     string
	Subject     string
	Issuer      string
	Thumbprint  string
	ContentType string
	Tags        map[string]string
}
