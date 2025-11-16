# API Reference

Azure Command Tower's internal API reference.

## Packages

### `internal/auth`

Azure authentication package.

#### Functions

- `NewAzureAuth() (*azidentity.DefaultAzureCredential, error)`
  - Creates Azure credentials using Azure CLI
  - Returns error if authentication fails

### `internal/azure`

Azure SDK client wrappers.

#### Types

- `Client`: Main Azure client wrapper
  - `SubscriptionsClient`: Azure subscriptions client
  - Methods for listing subscriptions, resource groups, resources, etc.

### `internal/models`

Data models for Azure resources.

#### Types

- `Subscription`: Azure subscription
- `ResourceGroup`: Azure resource group
- `Resource`: Azure resource
- `Container`: Storage container
- `Blob`: Storage blob
- `UserInfo`: Azure user information

### `internal/navigation`

Navigation state management.

#### Types

- `State`: Navigation state
  - Tracks current view
  - Tracks selected resources
  - Manages navigation hierarchy

### `internal/ui`

Terminal UI components.

#### Types

- `App`: Main application coordinator
- `SubscriptionsView`: Subscriptions list view
- `ResourceGroupsView`: Resource groups list view
- `ResourcesView`: Resources list view
- `DetailsView`: Resource details view
- `StorageExplorerView`: Storage explorer view
- `BlobsView`: Blobs list view

### `pkg/resource`

Resource handler system.

#### Types

- `Registry`: Resource handler registry
- `Handler`: Resource handler interface
- `DefaultHandler`: Default resource handler
- `StorageHandler`: Storage account handler

## Extending Azure Command Tower

### Adding a Resource Handler

Implement the `Handler` interface:

```go
type Handler interface {
    GetResourceType() string
    GetDisplayName() string
    CanExplore() bool
    GetActions() []Action
    // ... other methods
}
```

Register in `main.go`:

```go
registry.RegisterHandler(yourHandler)
```

### Adding a View

1. Create view component in `internal/ui`
2. Add view type to `internal/navigation`
3. Add navigation methods to `App`
4. Update key bindings

## Internal APIs

These APIs are internal and may change. For stable APIs, see the public interfaces in `pkg/resource`.

