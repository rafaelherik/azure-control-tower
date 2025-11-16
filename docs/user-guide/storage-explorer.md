# Storage Explorer

Azure Command Tower includes a built-in storage explorer for browsing Azure Storage accounts, containers, and blobs.

## Accessing Storage Explorer

1. Navigate to a storage account resource
2. Press `e` to explore the storage account
3. You'll see the Storage Explorer view with all containers

## Features

### Container View

The Storage Explorer shows:
- Container names
- Container properties (public access level, etc.)
- Last modified dates

**Actions:**
- `Enter`: Open container to view blobs
- `d`: View container details
- `/`: Filter containers

### Blob View

When you open a container, you'll see:
- Blobs and folders
- Blob sizes
- Last modified dates
- Content types

**Navigation:**
- Folders can be navigated like a file system
- Press `Enter` on a folder to navigate into it
- Press `ESC` to go back to parent folder or container list

### Blob Details

View detailed information about blobs:
- Full path
- Size
- Content type
- Last modified
- ETag
- Metadata

## Folder Navigation

The blob view supports hierarchical folder navigation:
- Navigate into folders with `Enter`
- Use `ESC` to go back to parent folder
- The breadcrumb shows your current path

## Filtering

Filter containers and blobs just like other views:
- Press `/` to activate filter
- Type to search by name
- Filter is case-insensitive

## Use Cases

- **Quick File Access**: Browse and find files in storage accounts
- **Container Management**: View all containers in a storage account
- **Blob Inspection**: Check blob properties and metadata
- **Folder Navigation**: Navigate through blob storage like a file system

## Tips

- Use the breadcrumb to see your current path in blob storage
- Filter is useful for finding specific blobs in large containers
- Container details show public access configuration
- Blob details include all metadata and properties

