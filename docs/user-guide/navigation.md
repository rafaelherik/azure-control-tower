# Navigation

Azure Command Tower provides a hierarchical navigation system for exploring your Azure resources.

## Navigation Hierarchy

```
Subscriptions
  └── Resource Groups
      └── Resource Types
          └── Resources
              └── Storage Explorer (for storage accounts)
                  └── Containers
                      └── Blobs
```

## Views

### Subscriptions View

The initial view when Azure Command Tower starts. Displays all Azure subscriptions you have access to.

**Actions:**
- `Enter`: Navigate to resource groups for the selected subscription
- `d`: View subscription details
- `/`: Filter subscriptions

### Resource Groups View

Shows all resource groups within the selected subscription.

**Actions:**
- `Enter`: Navigate to resource types for the selected resource group
- `d`: View resource group details
- `ESC`: Go back to subscriptions
- `/`: Filter resource groups

### Resource Types View

Displays a summary of resource types in the selected resource group, showing the count of each type.

**Actions:**
- `Enter`: Navigate to resources of the selected type
- `ESC`: Go back to resource groups
- `/`: Filter resource types
- `m`: Open resource type menu

### Resources View

Shows all resources of a specific type (or all resources in a resource group).

**Actions:**
- `Enter`: View resource details
- `d`: View resource details
- `e`: Explore storage (for storage accounts)
- `ESC`: Go back to resource types
- `/`: Filter resources

### Storage Explorer View

Available for storage account resources. Shows all containers in the storage account.

**Actions:**
- `Enter`: Navigate to blobs in the selected container
- `d`: View container details
- `ESC`: Go back to resources
- `/`: Filter containers

### Blobs View

Displays blobs and folders within a container.

**Actions:**
- `Enter`: Navigate into folder or view blob details
- `d`: View blob details
- `ESC`: Go back to storage explorer or parent folder
- `/`: Filter blobs

### Details View

Shows detailed information about the selected resource, subscription, resource group, container, or blob.

**Actions:**
- `ESC`: Go back to previous view

## Breadcrumb Navigation

The breadcrumb at the top of the screen shows your current navigation path, making it easy to understand where you are in the hierarchy.

## Menu

Press `m` to open the resource type menu, which provides quick access to all available resource types in the current context.

