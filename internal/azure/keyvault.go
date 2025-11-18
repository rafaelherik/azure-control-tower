package azure

import (
	"context"
	"fmt"
	"strings"

	"azure-control-tower/internal/models"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

// ListKeyVaults lists all Key Vaults in a resource group
func (c *Client) ListKeyVaults(ctx context.Context, subscriptionID, resourceGroupName string) ([]*models.KeyVault, error) {
	client, err := armkeyvault.NewVaultsClient(subscriptionID, c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Key Vault client: %w", err)
	}

	pager := client.NewListByResourceGroupPager(resourceGroupName, nil)
	var keyVaults []*models.KeyVault

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page: %w", err)
		}

		for _, vault := range page.Value {
			if vault == nil {
				continue
			}

			kv := &models.KeyVault{
				Tags:       make(map[string]*string),
				Properties: make(map[string]interface{}),
			}

			if vault.ID != nil {
				kv.ID = *vault.ID
			}
			if vault.Name != nil {
				kv.Name = *vault.Name
			}
			if vault.Location != nil {
				kv.Location = *vault.Location
			}

			// Extract resource group from ID
			if kv.ID != "" {
				parts := strings.Split(kv.ID, "/")
				for i, part := range parts {
					if strings.EqualFold(part, "resourceGroups") && i+1 < len(parts) {
						kv.ResourceGroup = parts[i+1]
						break
					}
				}
			}

			if vault.Properties != nil {
				if vault.Properties.VaultURI != nil {
					kv.VaultURI = *vault.Properties.VaultURI
				}
				if vault.Properties.TenantID != nil {
					kv.TenantID = *vault.Properties.TenantID
				}
				if vault.Properties.EnabledForDeployment != nil {
					kv.EnabledForDeploy = *vault.Properties.EnabledForDeployment
				}
				if vault.Properties.EnabledForDiskEncryption != nil {
					kv.EnabledForDisk = *vault.Properties.EnabledForDiskEncryption
				}
				if vault.Properties.EnabledForTemplateDeployment != nil {
					kv.EnabledForTemplate = *vault.Properties.EnabledForTemplateDeployment
				}

				// Add SKU info
				if vault.Properties.SKU != nil && vault.Properties.SKU.Name != nil {
					kv.SKU = string(*vault.Properties.SKU.Name)
					kv.Properties["sku"] = kv.SKU
				}

				// Add network rules info
				if vault.Properties.NetworkACLs != nil {
					if vault.Properties.NetworkACLs.DefaultAction != nil {
						kv.Properties["networkDefaultAction"] = string(*vault.Properties.NetworkACLs.DefaultAction)
					}
				}
			}

			if vault.Tags != nil {
				for k, v := range vault.Tags {
					kv.Tags[k] = v
				}
			}

			keyVaults = append(keyVaults, kv)
		}
	}

	return keyVaults, nil
}

// ListSecrets lists all secrets in a Key Vault
func (c *Client) ListSecrets(ctx context.Context, vaultURL string) ([]*models.Secret, error) {
	client, err := azsecrets.NewClient(vaultURL, c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create secrets client: %w", err)
	}

	pager := client.NewListSecretPropertiesPager(nil)
	var secrets []*models.Secret

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page: %w", err)
		}

		for _, secretProps := range page.Value {
			if secretProps == nil || secretProps.ID == nil {
				continue
			}

			secret := &models.Secret{
				Tags: make(map[string]string),
			}

			// Extract name from ID (format: https://vault.vault.azure.net/secrets/name/version)
			if secretProps.ID != nil {
				parts := strings.Split(secretProps.ID.Name(), "/")
				if len(parts) > 0 {
					secret.Name = parts[len(parts)-1]
				}
			}

			if secretProps.Attributes != nil {
				if secretProps.Attributes.Enabled != nil {
					secret.Enabled = *secretProps.Attributes.Enabled
				}
				secret.Created = secretProps.Attributes.Created
				secret.Updated = secretProps.Attributes.Updated
				secret.Expires = secretProps.Attributes.Expires
				secret.NotBefore = secretProps.Attributes.NotBefore
			}

			if secretProps.ContentType != nil {
				secret.ContentType = *secretProps.ContentType
			}

			if secretProps.Tags != nil {
				for k, v := range secretProps.Tags {
					if v != nil {
						secret.Tags[k] = *v
					}
				}
			}

			secrets = append(secrets, secret)
		}
	}

	return secrets, nil
}

// GetSecretValue retrieves the actual value of a secret
func (c *Client) GetSecretValue(ctx context.Context, vaultURL, secretName string) (string, error) {
	client, err := azsecrets.NewClient(vaultURL, c.credential, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create secrets client: %w", err)
	}

	resp, err := client.GetSecret(ctx, secretName, "", nil)
	if err != nil {
		return "", fmt.Errorf("failed to get secret: %w", err)
	}

	if resp.Value == nil {
		return "", fmt.Errorf("secret value is nil")
	}

	return *resp.Value, nil
}

// ListKeys lists all keys in a Key Vault
func (c *Client) ListKeys(ctx context.Context, vaultURL string) ([]*models.Key, error) {
	client, err := azkeys.NewClient(vaultURL, c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create keys client: %w", err)
	}

	pager := client.NewListKeyPropertiesPager(nil)
	var keys []*models.Key

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page: %w", err)
		}

		for _, keyProps := range page.Value {
			if keyProps == nil || keyProps.KID == nil {
				continue
			}

			key := &models.Key{
				Tags: make(map[string]string),
			}

			// Extract name from ID
			if keyProps.KID != nil {
				parts := strings.Split(keyProps.KID.Name(), "/")
				if len(parts) > 0 {
					key.Name = parts[len(parts)-1]
				}
			}

			if keyProps.Attributes != nil {
				if keyProps.Attributes.Enabled != nil {
					key.Enabled = *keyProps.Attributes.Enabled
				}
				key.Created = keyProps.Attributes.Created
				key.Updated = keyProps.Attributes.Updated
				key.Expires = keyProps.Attributes.Expires
				key.NotBefore = keyProps.Attributes.NotBefore
			}

			if keyProps.Tags != nil {
				for k, v := range keyProps.Tags {
					if v != nil {
						key.Tags[k] = *v
					}
				}
			}

			keys = append(keys, key)
		}
	}

	return keys, nil
}

// GetKeyDetails retrieves detailed information about a key
func (c *Client) GetKeyDetails(ctx context.Context, vaultURL, keyName string) (*models.Key, error) {
	client, err := azkeys.NewClient(vaultURL, c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create keys client: %w", err)
	}

	resp, err := client.GetKey(ctx, keyName, "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	key := &models.Key{
		Tags: make(map[string]string),
	}

	if resp.Key.KID != nil {
		parts := strings.Split(resp.Key.KID.Name(), "/")
		if len(parts) > 0 {
			key.Name = parts[len(parts)-1]
		}
		if resp.Key.KID.Version() != "" {
			key.Version = resp.Key.KID.Version()
		}
	}

	if resp.Key.Kty != nil {
		key.KeyType = string(*resp.Key.Kty)
	}

	if resp.Attributes != nil {
		if resp.Attributes.Enabled != nil {
			key.Enabled = *resp.Attributes.Enabled
		}
		key.Created = resp.Attributes.Created
		key.Updated = resp.Attributes.Updated
		key.Expires = resp.Attributes.Expires
		key.NotBefore = resp.Attributes.NotBefore
	}

	if resp.Tags != nil {
		for k, v := range resp.Tags {
			if v != nil {
				key.Tags[k] = *v
			}
		}
	}

	return key, nil
}

// ListCertificates lists all certificates in a Key Vault
func (c *Client) ListCertificates(ctx context.Context, vaultURL string) ([]*models.Certificate, error) {
	client, err := azcertificates.NewClient(vaultURL, c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificates client: %w", err)
	}

	pager := client.NewListCertificatePropertiesPager(nil)
	var certificates []*models.Certificate

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page: %w", err)
		}

		for _, certProps := range page.Value {
			if certProps == nil || certProps.ID == nil {
				continue
			}

			cert := &models.Certificate{
				Tags: make(map[string]string),
			}

			// Extract name from ID
			if certProps.ID != nil {
				parts := strings.Split(certProps.ID.Name(), "/")
				if len(parts) > 0 {
					cert.Name = parts[len(parts)-1]
				}
			}

			if certProps.Attributes != nil {
				if certProps.Attributes.Enabled != nil {
					cert.Enabled = *certProps.Attributes.Enabled
				}
				cert.Created = certProps.Attributes.Created
				cert.Updated = certProps.Attributes.Updated
				cert.Expires = certProps.Attributes.Expires
				cert.NotBefore = certProps.Attributes.NotBefore
			}

			if certProps.X509Thumbprint != nil {
				cert.Thumbprint = fmt.Sprintf("%x", certProps.X509Thumbprint)
			}

			if certProps.Tags != nil {
				for k, v := range certProps.Tags {
					if v != nil {
						cert.Tags[k] = *v
					}
				}
			}

			certificates = append(certificates, cert)
		}
	}

	return certificates, nil
}

// GetCertificateDetails retrieves detailed information about a certificate
func (c *Client) GetCertificateDetails(ctx context.Context, vaultURL, certName string) (*models.Certificate, error) {
	client, err := azcertificates.NewClient(vaultURL, c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificates client: %w", err)
	}

	resp, err := client.GetCertificate(ctx, certName, "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get certificate: %w", err)
	}

	cert := &models.Certificate{
		Tags: make(map[string]string),
	}

	if resp.ID != nil {
		parts := strings.Split(resp.ID.Name(), "/")
		if len(parts) > 0 {
			cert.Name = parts[len(parts)-1]
		}
		if resp.ID.Version() != "" {
			cert.Version = resp.ID.Version()
		}
	}

	if resp.Attributes != nil {
		if resp.Attributes.Enabled != nil {
			cert.Enabled = *resp.Attributes.Enabled
		}
		cert.Created = resp.Attributes.Created
		cert.Updated = resp.Attributes.Updated
		cert.Expires = resp.Attributes.Expires
		cert.NotBefore = resp.Attributes.NotBefore
	}

	// Certificate properties are in the Policy field
	if resp.Policy != nil {
		if resp.Policy.X509CertificateProperties != nil {
			if resp.Policy.X509CertificateProperties.Subject != nil {
				cert.Subject = *resp.Policy.X509CertificateProperties.Subject
			}
		}
		if resp.Policy.IssuerParameters != nil && resp.Policy.IssuerParameters.Name != nil {
			cert.Issuer = *resp.Policy.IssuerParameters.Name
		}
	}

	// Get thumbprint from CER
	if resp.CER != nil {
		// Calculate thumbprint from certificate data
		cert.Thumbprint = fmt.Sprintf("%x", resp.CER)[:40] // First 40 chars (SHA1)
	}

	if resp.Tags != nil {
		for k, v := range resp.Tags {
			if v != nil {
				cert.Tags[k] = *v
			}
		}
	}

	return cert, nil
}
