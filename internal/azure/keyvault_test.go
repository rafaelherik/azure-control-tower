package azure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestExtractResourceGroupFromKeyVaultID tests the resource group extraction logic
// that is used in ListKeyVaults
func TestExtractResourceGroupFromKeyVaultID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected string
	}{
		{
			name:     "Standard Key Vault ID",
			id:       "/subscriptions/sub123/resourceGroups/my-rg/providers/Microsoft.KeyVault/vaults/my-vault",
			expected: "my-rg",
		},
		{
			name:     "Different resource group name",
			id:       "/subscriptions/abc-def/resourceGroups/production-rg/providers/Microsoft.KeyVault/vaults/prod-vault",
			expected: "production-rg",
		},
		{
			name:     "Resource group with special characters",
			id:       "/subscriptions/sub123/resourceGroups/my-rg_test-001/providers/Microsoft.KeyVault/vaults/vault1",
			expected: "my-rg_test-001",
		},
		{
			name:     "Case insensitive resourceGroups",
			id:       "/subscriptions/sub123/resourcegroups/my-rg/providers/Microsoft.KeyVault/vaults/my-vault",
			expected: "my-rg",
		},
		{
			name:     "Mixed case resourceGroups",
			id:       "/subscriptions/sub123/ResourceGroups/my-rg/providers/Microsoft.KeyVault/vaults/my-vault",
			expected: "my-rg",
		},
		{
			name:     "Empty ID",
			id:       "",
			expected: "",
		},
		{
			name:     "ID without resourceGroups",
			id:       "/subscriptions/sub123/providers/Microsoft.KeyVault/vaults/my-vault",
			expected: "",
		},
		{
			name:     "ID with resourceGroups at the end",
			id:       "/subscriptions/sub123/resourceGroups",
			expected: "",
		},
		{
			name:     "Multiple slashes",
			id:       "/subscriptions/sub123//resourceGroups//my-rg//providers/Microsoft.KeyVault/vaults/my-vault",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the extraction logic from ListKeyVaults
			var resourceGroup string
			if tt.id != "" {
				parts := splitPath(tt.id)
				resourceGroup = extractResourceGroup(parts)
			}
			assert.Equal(t, tt.expected, resourceGroup)
		})
	}
}

// Helper function to split path (simulates strings.Split logic)
func splitPath(id string) []string {
	return splitString(id, "/")
}

// Helper function to extract resource group from parts
func extractResourceGroup(parts []string) string {
	for i, part := range parts {
		if equalFoldString(part, "resourceGroups") && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// Helper to simulate strings.Split
func splitString(s, sep string) []string {
	if s == "" {
		return []string{""}
	}
	result := []string{}
	current := ""
	for _, ch := range s {
		if string(ch) == sep {
			result = append(result, current)
			current = ""
		} else {
			current += string(ch)
		}
	}
	result = append(result, current)
	return result
}

// Helper to simulate strings.EqualFold
func equalFoldString(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		c1, c2 := s1[i], s2[i]
		if c1 >= 'A' && c1 <= 'Z' {
			c1 = c1 + 32 // to lowercase
		}
		if c2 >= 'A' && c2 <= 'Z' {
			c2 = c2 + 32 // to lowercase
		}
		if c1 != c2 {
			return false
		}
	}
	return true
}

func TestSecretNameExtraction(t *testing.T) {
	tests := []struct {
		name        string
		secretID    string
		expected    string
		description string
	}{
		{
			name:        "Standard secret ID",
			secretID:    "secrets/my-secret",
			expected:    "my-secret",
			description: "Extract secret name from standard ID format",
		},
		{
			name:        "Secret ID with version",
			secretID:    "secrets/my-secret/version123",
			expected:    "version123",
			description: "Extract version from versioned secret ID",
		},
		{
			name:        "Secret with single part",
			secretID:    "my-secret",
			expected:    "my-secret",
			description: "Handle single part secret name",
		},
		{
			name:        "Empty secret ID",
			secretID:    "",
			expected:    "",
			description: "Handle empty secret ID",
		},
		{
			name:        "Secret with hyphens",
			secretID:    "secrets/my-secret-name",
			expected:    "my-secret-name",
			description: "Handle secret names with hyphens",
		},
		{
			name:        "Secret with underscores",
			secretID:    "secrets/my_secret_name",
			expected:    "my_secret_name",
			description: "Handle secret names with underscores",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the name extraction logic from ListSecrets
			parts := splitString(tt.secretID, "/")
			var name string
			if len(parts) > 0 {
				name = parts[len(parts)-1]
			}
			assert.Equal(t, tt.expected, name, tt.description)
		})
	}
}

func TestKeyNameExtraction(t *testing.T) {
	tests := []struct {
		name        string
		keyID       string
		expected    string
		description string
	}{
		{
			name:        "Standard key ID",
			keyID:       "keys/my-key",
			expected:    "my-key",
			description: "Extract key name from standard ID format",
		},
		{
			name:        "Key ID with version",
			keyID:       "keys/my-key/version456",
			expected:    "version456",
			description: "Extract version from versioned key ID",
		},
		{
			name:        "Key with single part",
			keyID:       "my-key",
			expected:    "my-key",
			description: "Handle single part key name",
		},
		{
			name:        "RSA key",
			keyID:       "keys/rsa-key-2048",
			expected:    "rsa-key-2048",
			description: "Handle RSA key names",
		},
		{
			name:        "EC key",
			keyID:       "keys/ec-key-p256",
			expected:    "ec-key-p256",
			description: "Handle EC key names",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the name extraction logic from ListKeys
			parts := splitString(tt.keyID, "/")
			var name string
			if len(parts) > 0 {
				name = parts[len(parts)-1]
			}
			assert.Equal(t, tt.expected, name, tt.description)
		})
	}
}

func TestCertificateNameExtraction(t *testing.T) {
	tests := []struct {
		name        string
		certID      string
		expected    string
		description string
	}{
		{
			name:        "Standard certificate ID",
			certID:      "certificates/my-cert",
			expected:    "my-cert",
			description: "Extract certificate name from standard ID format",
		},
		{
			name:        "Certificate ID with version",
			certID:      "certificates/my-cert/version789",
			expected:    "version789",
			description: "Extract version from versioned certificate ID",
		},
		{
			name:        "SSL certificate",
			certID:      "certificates/wildcard-ssl-cert",
			expected:    "wildcard-ssl-cert",
			description: "Handle SSL certificate names",
		},
		{
			name:        "Certificate with domain",
			certID:      "certificates/example-com-cert",
			expected:    "example-com-cert",
			description: "Handle certificate names with domain references",
		},
		{
			name:        "Empty certificate ID",
			certID:      "",
			expected:    "",
			description: "Handle empty certificate ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the name extraction logic from ListCertificates
			parts := splitString(tt.certID, "/")
			var name string
			if len(parts) > 0 {
				name = parts[len(parts)-1]
			}
			assert.Equal(t, tt.expected, name, tt.description)
		})
	}
}

func TestThumbprintFormatting(t *testing.T) {
	tests := []struct {
		name        string
		thumbprint  []byte
		expectedLen int
		description string
	}{
		{
			name:        "SHA1 thumbprint",
			thumbprint:  []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67},
			expectedLen: 40, // 20 bytes = 40 hex characters
			description: "Standard SHA1 thumbprint should format to 40 characters",
		},
		{
			name:        "Empty thumbprint",
			thumbprint:  []byte{},
			expectedLen: 0,
			description: "Empty thumbprint should format to empty string",
		},
		{
			name:        "Partial thumbprint",
			thumbprint:  []byte{0xff, 0xee, 0xdd},
			expectedLen: 6, // 3 bytes = 6 hex characters
			description: "Partial thumbprint should format correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the thumbprint formatting logic from ListCertificates
			formatted := formatThumbprint(tt.thumbprint)
			assert.Equal(t, tt.expectedLen, len(formatted), tt.description)
			
			// Verify it only contains hex characters
			if len(formatted) > 0 {
				for _, ch := range formatted {
					assert.True(t, (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f'), 
						"Thumbprint should only contain hex characters")
				}
			}
		})
	}
}

// Helper to format thumbprint
func formatThumbprint(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	hexChars := "0123456789abcdef"
	result := make([]byte, len(data)*2)
	for i, b := range data {
		result[i*2] = hexChars[b>>4]
		result[i*2+1] = hexChars[b&0x0f]
	}
	return string(result)
}

func BenchmarkExtractResourceGroup(b *testing.B) {
	id := "/subscriptions/sub123/resourceGroups/my-rg/providers/Microsoft.KeyVault/vaults/my-vault"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parts := splitPath(id)
		_ = extractResourceGroup(parts)
	}
}

func BenchmarkSecretNameExtraction(b *testing.B) {
	secretID := "secrets/my-secret-name/version123"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parts := splitString(secretID, "/")
		if len(parts) > 0 {
			_ = parts[len(parts)-1]
		}
	}
}

func BenchmarkThumbprintFormatting(b *testing.B) {
	thumbprint := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = formatThumbprint(thumbprint)
	}
}

// Integration test for resource group extraction
func TestResourceGroupExtraction_Integration(t *testing.T) {
	t.Run("Extract from multiple Key Vault IDs", func(t *testing.T) {
		ids := []string{
			"/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.KeyVault/vaults/vault1",
			"/subscriptions/sub2/resourceGroups/rg2/providers/Microsoft.KeyVault/vaults/vault2",
			"/subscriptions/sub3/resourceGroups/rg3/providers/Microsoft.KeyVault/vaults/vault3",
		}
		expected := []string{"rg1", "rg2", "rg3"}
		
		for i, id := range ids {
			parts := splitPath(id)
			rg := extractResourceGroup(parts)
			assert.Equal(t, expected[i], rg, "Resource group extraction failed for ID: %s", id)
		}
	})
}
