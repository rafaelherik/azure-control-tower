# Filtering

Azure Command Tower includes a powerful filtering system to quickly find resources.

## Activating Filter Mode

Press `/` in any table view to activate filter mode. A filter input field will appear at the top of the view.

## How Filtering Works

- **Case Insensitive**: Filters are case-insensitive
- **Multi-Column**: Filters match against all visible columns
- **Real-time**: Results update as you type
- **Clear Indication**: The footer shows filtered count vs total count

## Using Filters

1. Press `/` to open the filter input
2. Type your search term
3. The table updates automatically as you type
4. Press `ESC` to cancel the filter
5. Press `Enter` to apply (though filtering happens automatically)

## Filter Examples

### Filter by Name
Type part of a resource name to find matching resources:
```
storage
```
Will match all resources containing "storage" in their name.

### Filter by Type
Type a resource type to filter:
```
Microsoft.Storage/storageAccounts
```

### Filter by Resource Group
In resources view, filter by resource group name:
```
production
```

## Filter Indicators

The footer shows:
- **Total count**: Total number of items
- **Filtered count**: Number of items matching the filter
- **Filter indicator**: Shows when a filter is active

Example footer:
```
Showing 5 of 20 items (filtered) | Enter: select, d: details, ESC: back, /: filter, q: quit
```

## Clearing Filters

- Press `ESC` while in filter mode to cancel
- The filter is automatically cleared when you navigate away from a view

## Tips

- Filters persist while you're in the same view
- Use partial matches for broader searches
- Combine with navigation to quickly find specific resources
- Filter works in all table views: subscriptions, resource groups, resources, containers, and blobs

