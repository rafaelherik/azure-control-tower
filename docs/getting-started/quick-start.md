# Quick Start

Get up and running with Azure Command Tower (azct) in minutes.

## Authentication

Before using azct, you need to authenticate with Azure:

```bash
az login
```

This will open your browser to authenticate. Once authenticated, azct will use these credentials automatically.

## Running azct

Simply run:

```bash
azct
```

The application will start and display your Azure subscriptions.

## Basic Navigation

### Viewing Subscriptions

When azct starts, you'll see a list of your Azure subscriptions. Use the arrow keys to navigate and press `Enter` to select a subscription.

### Exploring Resource Groups

After selecting a subscription, you'll see the resource groups in that subscription. Select a resource group to see its resources.

### Viewing Resources

From the resource groups view, you can:
- See a summary of resource types
- Navigate to specific resource types
- View individual resources

### Storage Explorer

For storage accounts, you can explore containers and blobs:
1. Navigate to a storage account resource
2. Press `e` to explore
3. Browse containers and blobs

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑/↓` | Navigate up/down |
| `Enter` | Select/view details |
| `/` | Open filter/search |
| `d` | Show details |
| `e` | Explore storage |
| `m` | Open resource type menu |
| `ESC` | Go back |
| `q` | Quit |

## Next Steps

- Learn more about [Navigation](user-guide/navigation.md)
- See all [Keyboard Shortcuts](user-guide/keyboard-shortcuts.md)
- Explore the [Storage Explorer](user-guide/storage-explorer.md)

