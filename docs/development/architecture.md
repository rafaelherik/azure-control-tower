# Architecture

Overview of Azure Command Tower's architecture and design.

## Project Structure

```
azct/
├── cmd/
│   └── azct/           # Main application entry point
├── internal/
│   ├── auth/           # Azure authentication
│   ├── azure/          # Azure SDK client wrappers
│   ├── models/         # Data models
│   ├── navigation/     # Navigation state management
│   └── ui/             # Terminal UI components
└── pkg/
    └── resource/       # Resource handlers and registry
```

## Components

### Authentication (`internal/auth`)

Handles Azure authentication using Azure CLI credentials:
- Uses `DefaultAzureCredential` from Azure SDK
- Verifies credentials before starting
- Provides clear error messages

### Azure Client (`internal/azure`)

Wraps Azure SDK clients:
- Subscriptions client
- Resources client
- Storage client
- Provides unified interface for Azure operations

### Models (`internal/models`)

Data structures representing Azure resources:
- Subscriptions
- Resource Groups
- Resources
- Storage containers and blobs
- User information

### Navigation (`internal/navigation`)

Manages application navigation state:
- Current view tracking
- Selected resources
- Navigation history
- Breadcrumb information

### UI Components (`internal/ui`)

Terminal UI components built with tview:
- App: Main application coordinator
- Views: Subscriptions, Resource Groups, Resources, etc.
- Header: User info and selected subscription
- Footer: Status and shortcuts
- Details: Resource detail views
- Filter: Search/filter functionality

### Resource Handlers (`pkg/resource`)

Extensible system for handling different resource types:
- Registry: Manages resource handlers
- Handlers: Custom logic for resource types
- Default handler: Generic resource handling
- Storage handler: Special handling for storage accounts

## Design Patterns

### Registry Pattern

Resource handlers use a registry pattern for extensibility:
- Handlers register themselves
- App queries registry for appropriate handler
- Easy to add new resource type handlers

### View Pattern

UI uses a view-based architecture:
- Each view is a separate component
- Views handle their own state
- App coordinates view transitions

### State Management

Navigation state is centralized:
- Single source of truth for navigation
- Views update based on state
- Easy to add new views

## Data Flow

1. User interacts with UI
2. View handles input
3. View calls app navigation methods
4. App updates navigation state
5. App loads data from Azure client
6. App updates view with data
7. View renders to screen

## Extensibility

### Adding Resource Handlers

1. Implement `Handler` interface
2. Register in main.go
3. Handler provides custom actions and display

### Adding Views

1. Create new view component
2. Add to navigation state
3. Add navigation methods to app
4. Update key bindings

## Dependencies

- **tview/tcell**: Terminal UI framework
- **Azure SDK for Go**: Azure API access
- **Go standard library**: Core functionality

## Future Enhancements

- Plugin system for resource handlers
- Customizable themes
- Export functionality
- Resource actions (create, delete, etc.)

