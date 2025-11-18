# Key Vault Explorer

Azure Command Tower includes comprehensive support for browsing and managing Azure Key Vault resources, including secrets, keys, and certificates.

## Accessing Key Vault Explorer

1. Navigate to a Key Vault resource in your resource group
2. Press `e` to explore the Key Vault
3. You'll see the Key Vault Explorer view with three main categories

## Features

### Key Vault Explorer View

The Key Vault Explorer provides access to three main item types:

- **🔐 Secrets**: Manage secret values and configurations
- **🔑 Keys**: Manage cryptographic keys
- **📜 Certificates**: Manage SSL/TLS certificates

**Actions:**
- `Enter`: Open the selected item type to view its contents
- `ESC`: Go back to resource list
- `/`: Filter items

### Secrets Management

#### Listing Secrets

The secrets view displays:
- Secret names
- Enabled status (✓ or ✗)
- Content type
- Last updated date

**Actions:**
- `v`: View secret value (requires confirmation)
- `d`: View secret details
- `Enter`: View secret details
- `ESC`: Go back to Key Vault Explorer
- `/`: Filter secrets

#### Viewing Secret Values

When you press `v` to view a secret value:

1. A security confirmation dialog appears:
   ```
   Are you sure you want to view the value of secret 'my-secret'?
   
   ⚠️ This will display sensitive information on screen.
   
   [View]  [Cancel]
   ```

2. Select "View" to fetch and display the secret value
3. The value is displayed in a modal window
4. Press any key to close and return to the secrets list

**Security Note**: Secret values are only fetched when explicitly requested and are not cached.

#### Secret Details

Secret details include:
- Key Vault name
- Secret name
- Enabled status
- Content type
- Creation date
- Last updated date
- Expiration date (if set)
- Not before date (if set)
- Tags (if any)

### Keys Management

#### Listing Keys

The keys view displays:
- Key names
- Key types (RSA, EC, etc.)
- Enabled status (✓ or ✗)
- Last updated date

**Actions:**
- `d`: View key details
- `Enter`: View key details
- `ESC`: Go back to Key Vault Explorer
- `/`: Filter keys

#### Key Details

Key details include:
- Key Vault name
- Key name
- Key type (RSA, RSA-HSM, EC, EC-HSM, etc.)
- Version
- Enabled status
- Creation date
- Last updated date
- Expiration date (if set)
- Not before date (if set)
- Tags (if any)

### Certificates Management

#### Listing Certificates

The certificates view displays:
- Certificate names
- Enabled status (✓ or ✗)
- Expiration dates with visual warnings
- Last updated date

**Visual Indicators:**
- `⚠️ 2024-01-15 (EXPIRED)`: Shows expired certificates with warning icon

**Actions:**
- `d`: View certificate details
- `Enter`: View certificate details
- `ESC`: Go back to Key Vault Explorer
- `/`: Filter certificates

#### Certificate Details

Certificate details include:
- Key Vault name
- Certificate name
- Subject (DN)
- Issuer
- Thumbprint
- Version
- Enabled status
- Creation date
- Last updated date
- Expiration date with expiration warning
- Not before date (if set)
- Tags (if any)

## Navigation Flow

```
Resource Groups
  └─> Key Vaults (filtered list)
        └─> Key Vault Explorer
              ├─> Secrets
              │     ├─> Secret Details
              │     └─> View Secret Value (with confirmation)
              ├─> Keys
              │     └─> Key Details
              └─> Certificates
                    └─> Certificate Details
```

## Filtering

Filter works across all Key Vault views:
- Press `/` to activate filter
- Type to search by name
- Filter is case-insensitive
- Works in Key Vault list, secrets, keys, and certificates views

## Keyboard Shortcuts Summary

### Key Vault Explorer
| Key | Action |
|-----|--------|
| `Enter` | Open item type (secrets/keys/certificates) |
| `ESC` | Go back to resource list |
| `/` | Filter items |
| `q` | Quit application |

### Secrets View
| Key | Action |
|-----|--------|
| `v` | View secret value (with confirmation) |
| `d` | View secret details |
| `Enter` | View secret details |
| `ESC` | Go back to Key Vault Explorer |
| `/` | Filter secrets |
| `q` | Quit application |

### Keys View
| Key | Action |
|-----|--------|
| `d` | View key details |
| `Enter` | View key details |
| `ESC` | Go back to Key Vault Explorer |
| `/` | Filter keys |
| `q` | Quit application |

### Certificates View
| Key | Action |
|-----|--------|
| `d` | View certificate details |
| `Enter` | View certificate details |
| `ESC` | Go back to Key Vault Explorer |
| `/` | Filter certificates |
| `q` | Quit application |

## Use Cases

- **Secret Management**: Browse and view secrets stored in Key Vaults
- **Security Auditing**: Check which secrets, keys, and certificates are enabled/disabled
- **Expiration Monitoring**: Identify expired or soon-to-expire certificates
- **Quick Access**: Quickly view secret values when needed during troubleshooting
- **Inventory Review**: List all cryptographic assets in a Key Vault

## Security Considerations

### Secret Value Protection

- Secret values are never displayed without explicit user confirmation
- A security warning is shown before displaying sensitive information
- Secret values are fetched on-demand and not cached
- Use the `v` key only when you need to view the actual secret value

### Access Requirements

To use Key Vault Explorer, you need:
- **List permissions** on the Key Vault to see items
- **Get permissions** on secrets/keys/certificates to view their values
- Proper Azure RBAC roles (e.g., "Key Vault Secrets User", "Key Vault Reader")

If you lack permissions, operations will fail with an error message.

## Tips

- **Use filtering** to quickly find specific secrets, keys, or certificates in large vaults
- **Check expiration dates** regularly to avoid expired certificates causing outages
- **Review enabled status** to ensure only necessary items are active
- **Confirm secret values carefully** before viewing to avoid accidental exposure
- **Use breadcrumb navigation** to track your current location in the Key Vault hierarchy
- **Look for visual warnings** (⚠️) on expired certificates

## Common Operations

### Viewing a Secret Value

1. Navigate to Key Vaults in your resource group
2. Select a Key Vault and press `e`
3. Select "Secrets" and press `Enter`
4. Navigate to the desired secret
5. Press `v` to view the value
6. Confirm by selecting "View"
7. The secret value is displayed
8. Press any key to close

### Checking Certificate Expiration

1. Navigate to Key Vaults in your resource group
2. Select a Key Vault and press `e`
3. Select "Certificates" and press `Enter`
4. Review the expiration dates column
5. Look for the ⚠️ warning icon for expired certificates
6. Press `d` or `Enter` on any certificate to see full details

### Finding a Specific Key

1. Navigate to Key Vaults in your resource group
2. Select a Key Vault and press `e`
3. Select "Keys" and press `Enter`
4. Press `/` to activate filter
5. Type the key name (or part of it)
6. Press `Enter` on the found key to see details
