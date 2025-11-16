# Authentication

Azure Command Tower uses Azure CLI credentials for authentication.

## Prerequisites

Before using Azure Command Tower, you need to have the Azure CLI installed and authenticated.

## Setting Up Authentication

1. **Install Azure CLI** (if not already installed):
   - macOS: `brew install azure-cli`
   - Linux: Follow [Azure CLI installation guide](https://docs.microsoft.com/cli/azure/install-azure-cli)
   - Windows: Download from [Azure CLI website](https://aka.ms/installazurecliwindows)

2. **Login to Azure**:
   ```bash
   az login
   ```
   
   This will open your browser to authenticate with your Azure account.

3. **Verify Authentication**:
   ```bash
   az account show
   ```
   
   This should display your current Azure account information.

## How Azure Command Tower Uses Credentials

Azure Command Tower uses the `DefaultAzureCredential` from the Azure SDK, which automatically:
- Uses credentials from Azure CLI (`az login`)
- Falls back to environment variables if configured
- Uses managed identity when running on Azure resources

No additional configuration is needed - just run `az login` and Azure Command Tower will use those credentials automatically.

## Multiple Subscriptions

If you have access to multiple Azure subscriptions, Azure Command Tower will display all of them when you start the application. You can select which subscription to explore from the subscriptions view.

## Troubleshooting

### "Failed to authenticate with Azure" Error

If you see an authentication error:
1. Make sure you've run `az login`
2. Verify your credentials are still valid: `az account show`
3. Try logging in again: `az login`

### No Subscriptions Visible

If no subscriptions appear:
1. Verify you have access to subscriptions: `az account list`
2. Check if you need to set a default subscription: `az account set --subscription <subscription-id>`

