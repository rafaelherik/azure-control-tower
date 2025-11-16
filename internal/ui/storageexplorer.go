package ui

import (
	"context"
	"fmt"

	"azure-control-tower/internal/models"

	"github.com/rivo/tview"
)

// ContainerRowData wraps Container with context info for display
type ContainerRowData struct {
	Container *models.Container
}

// StorageExplorerView displays containers in a storage account
type StorageExplorerView struct {
	*TableView
	containers      []*models.Container
	storageAccount  string
	onSelect        func(container *models.Container)
	onShowDetails   func(container *models.Container)
}

// NewStorageExplorerView creates a new storage explorer view
func NewStorageExplorerView() *StorageExplorerView {
	sev := &StorageExplorerView{}

	// Create table configuration
	config := &TableConfig{
		Title: "", // Will be set in LoadContainers
		Columns: []ColumnConfig{
			{Name: "Name", Align: tview.AlignLeft},
			{Name: "Public Access", Align: tview.AlignLeft},
			{Name: "Last Modified", Align: tview.AlignLeft},
		},
		RowActions: []RowAction{
			{
				Rune:  'd',
				Label: "Details",
				Callback: func(rowIndex int, data interface{}) bool {
					if rowData, ok := data.(*ContainerRowData); ok && sev.onShowDetails != nil {
						sev.onShowDetails(rowData.Container)
						return true
					}
					return false
				},
			},
		},
		OnSelect: func(rowIndex int, data interface{}) {
			// Enter key on a container - navigate to blobs
			if rowData, ok := data.(*ContainerRowData); ok && sev.onSelect != nil {
				sev.onSelect(rowData.Container)
			}
		},
		GetCellValue: func(data interface{}, columnIndex int) string {
			rowData, ok := data.(*ContainerRowData)
			if !ok {
				return ""
			}
			switch columnIndex {
			case 0:
				return rowData.Container.Name
			case 1:
				if rowData.Container.PublicAccess == "" {
					return "Private"
				}
				return rowData.Container.PublicAccess
			case 2:
				return rowData.Container.LastModified.Format("2006-01-02 15:04:05")
			default:
				return ""
			}
		},
	}

	sev.TableView = NewTableView(config)
	return sev
}

// LoadContainers loads containers into the view
func (sev *StorageExplorerView) LoadContainers(ctx context.Context, containers []*models.Container, storageAccount string) error {
	sev.containers = containers
	sev.storageAccount = storageAccount

	// Update title
	sev.SetTitle("")

	// Convert to interface{} slice with row data
	data := make([]interface{}, len(containers))
	for i, container := range containers {
		data[i] = &ContainerRowData{
			Container: container,
		}
	}

	sev.LoadData(data)
	return nil
}

// SetOnSelect sets the callback for when a container is selected (Enter key)
func (sev *StorageExplorerView) SetOnSelect(callback func(*models.Container)) {
	sev.onSelect = callback
}

// SetOnShowDetails sets the callback for when details are requested (d key)
func (sev *StorageExplorerView) SetOnShowDetails(callback func(*models.Container)) {
	sev.onShowDetails = callback
}

// GetStorageAccount returns the current storage account name
func (sev *StorageExplorerView) GetStorageAccount() string {
	return sev.storageAccount
}

// BlobRowData wraps Blob with context info for display
type BlobRowData struct {
	Blob *models.Blob
}

// BlobsView displays blobs in a container
type BlobsView struct {
	*TableView
	blobs          []*models.Blob
	containerName  string
	storageAccount  string
	pathPrefix      string // Current folder path prefix
	onShowDetails   func(blob *models.Blob)
	onNavigateFolder func(folderPath string) // Callback for folder navigation
}

// NewBlobsView creates a new blobs view
func NewBlobsView() *BlobsView {
	bv := &BlobsView{}

	// Create table configuration
	config := &TableConfig{
		Title: "", // Will be set in LoadBlobs
		Columns: []ColumnConfig{
			{Name: "Name", Align: tview.AlignLeft},
			{Name: "Size", Align: tview.AlignRight},
			{Name: "Content Type", Align: tview.AlignLeft},
			{Name: "Last Modified", Align: tview.AlignLeft},
		},
		RowActions: []RowAction{
			{
				Rune:  'd',
				Label: "Details",
				Callback: func(rowIndex int, data interface{}) bool {
					if rowData, ok := data.(*BlobRowData); ok && bv.onShowDetails != nil {
						bv.onShowDetails(rowData.Blob)
						return true
					}
					return false
				},
			},
		},
		OnSelect: func(rowIndex int, data interface{}) {
			// Enter key on a blob - navigate into folder or show details
			if rowData, ok := data.(*BlobRowData); ok {
				if rowData.Blob.IsDirectory && bv.onNavigateFolder != nil {
					// Navigate into folder
					bv.onNavigateFolder(rowData.Blob.Name)
				} else if bv.onShowDetails != nil {
					// Show details for file
					bv.onShowDetails(rowData.Blob)
				}
			}
		},
		GetCellValue: func(data interface{}, columnIndex int) string {
			rowData, ok := data.(*BlobRowData)
			if !ok {
				return ""
			}
			switch columnIndex {
			case 0:
				// Show icon and display name
				icon := "📄" // File icon
				if rowData.Blob.IsDirectory {
					icon = "📁" // Folder icon
				}
				displayName := rowData.Blob.DisplayName
				if displayName == "" {
					displayName = rowData.Blob.Name
				}
				return fmt.Sprintf("%s %s", icon, displayName)
			case 1:
				if rowData.Blob.IsDirectory {
					return "-"
				}
				return formatSize(rowData.Blob.Size)
			case 2:
				if rowData.Blob.IsDirectory {
					return "-"
				}
				return rowData.Blob.ContentType
			case 3:
				if rowData.Blob.IsDirectory {
					return "-"
				}
				return rowData.Blob.LastModified.Format("2006-01-02 15:04:05")
			default:
				return ""
			}
		},
	}

	bv.TableView = NewTableView(config)
	return bv
}

// LoadBlobs loads blobs into the view
func (bv *BlobsView) LoadBlobs(ctx context.Context, blobs []*models.Blob, containerName, storageAccount, pathPrefix string) error {
	bv.blobs = blobs
	bv.containerName = containerName
	bv.storageAccount = storageAccount
	bv.pathPrefix = pathPrefix

	// Update title
	bv.SetTitle("")

	// Convert to interface{} slice with row data
	data := make([]interface{}, len(blobs))
	for i, blob := range blobs {
		data[i] = &BlobRowData{
			Blob: blob,
		}
	}

	bv.LoadData(data)
	return nil
}

// SetOnShowDetails sets the callback for when details are requested (d key or Enter)
func (bv *BlobsView) SetOnShowDetails(callback func(*models.Blob)) {
	bv.onShowDetails = callback
}

// SetOnNavigateFolder sets the callback for when a folder is selected (Enter key)
func (bv *BlobsView) SetOnNavigateFolder(callback func(string)) {
	bv.onNavigateFolder = callback
}

// GetContainerName returns the current container name
func (bv *BlobsView) GetContainerName() string {
	return bv.containerName
}

// GetPathPrefix returns the current path prefix
func (bv *BlobsView) GetPathPrefix() string {
	return bv.pathPrefix
}

// formatSize formats a size in bytes to a human-readable string
func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

