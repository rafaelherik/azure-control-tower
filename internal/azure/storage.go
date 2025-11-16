package azure

import (
	"context"
	"fmt"
	"path"
	"strings"

	"azure-control-tower/internal/models"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// ListContainers lists all containers in a storage account
func (c *Client) ListContainers(ctx context.Context, subscriptionID, resourceGroupName, storageAccountName string) ([]*models.Container, error) {
	// Get storage account keys
	keys, err := c.getStorageAccountKeys(ctx, subscriptionID, resourceGroupName, storageAccountName)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage account keys: %w", err)
	}

	if len(keys) == 0 {
		return nil, fmt.Errorf("no storage account keys found")
	}

	// Use the first key to create blob service client
	credential, err := azblob.NewSharedKeyCredential(storageAccountName, keys[0])
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", storageAccountName)
	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create blob client: %w", err)
	}

	// List containers
	pager := client.NewListContainersPager(nil)
	var containers []*models.Container

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page: %w", err)
		}

		for _, containerItem := range page.ContainerItems {
			if containerItem.Name == nil {
				continue
			}

			container := &models.Container{
				Name:     *containerItem.Name,
				Metadata: make(map[string]string),
			}

			if containerItem.Properties != nil {
				if containerItem.Properties.LastModified != nil {
					container.LastModified = *containerItem.Properties.LastModified
				}
				if containerItem.Properties.ETag != nil {
					container.ETag = string(*containerItem.Properties.ETag)
				}
				if containerItem.Properties.PublicAccess != nil {
					container.PublicAccess = string(*containerItem.Properties.PublicAccess)
				}
			}

			if containerItem.Metadata != nil {
				for k, v := range containerItem.Metadata {
					if v != nil {
						container.Metadata[k] = *v
					}
				}
			}

			containers = append(containers, container)
		}
	}

	return containers, nil
}

// ListBlobs lists blobs in a container hierarchically
// prefix: the current folder path prefix (e.g., "folder1/subfolder/")
// Returns only immediate children (folders and files) of the prefix
func (c *Client) ListBlobs(ctx context.Context, subscriptionID, resourceGroupName, storageAccountName, containerName, prefix string) ([]*models.Blob, error) {
	// Get storage account keys
	keys, err := c.getStorageAccountKeys(ctx, subscriptionID, resourceGroupName, storageAccountName)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage account keys: %w", err)
	}

	if len(keys) == 0 {
		return nil, fmt.Errorf("no storage account keys found")
	}

	// Use the first key to create blob service client
	credential, err := azblob.NewSharedKeyCredential(storageAccountName, keys[0])
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", storageAccountName)
	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create blob client: %w", err)
	}

	// List all blobs with the prefix (flat listing)
	options := &azblob.ListBlobsFlatOptions{}
	if prefix != "" {
		options.Prefix = &prefix
	}

	pager := client.NewListBlobsFlatPager(containerName, options)
	var allBlobs []*models.Blob
	seenItems := make(map[string]bool) // Track items we've already added

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page: %w", err)
		}

		for _, blobItem := range page.Segment.BlobItems {
			if blobItem.Name == nil {
				continue
			}

			blobName := *blobItem.Name

			// Skip if not under the current prefix
			if prefix != "" && !strings.HasPrefix(blobName, prefix) {
				continue
			}

			// Check if it's an immediate child
			if !isImmediateChild(blobName, prefix) {
				// It's a nested item - create a folder entry for its parent directory
				folderPath := getParentFolderPath(blobName, prefix)
				if folderPath != "" && !seenItems[folderPath] {
					displayName := getDisplayName(folderPath, prefix)
					folderBlob := &models.Blob{
						Name:        folderPath,
						DisplayName: displayName,
						Metadata:    make(map[string]string),
						IsDirectory: true,
						Size:        0,
					}
					allBlobs = append(allBlobs, folderBlob)
					seenItems[folderPath] = true
				}
				continue
			}

			// It's an immediate child - add it
			displayName := getDisplayName(blobName, prefix)

			// Check if it's a directory marker (ends with /)
			isDir := strings.HasSuffix(blobName, "/")

			blob := &models.Blob{
				Name:        blobName,
				DisplayName: displayName,
				Metadata:    make(map[string]string),
				IsDirectory: isDir,
			}

			if !isDir {
				// Only set properties for files
				if blobItem.Properties != nil {
					if blobItem.Properties.ContentLength != nil {
						blob.Size = *blobItem.Properties.ContentLength
					}
					if blobItem.Properties.ContentType != nil {
						blob.ContentType = *blobItem.Properties.ContentType
					}
					if blobItem.Properties.LastModified != nil {
						blob.LastModified = *blobItem.Properties.LastModified
					}
					if blobItem.Properties.ETag != nil {
						blob.ETag = string(*blobItem.Properties.ETag)
					}
				}

				if blobItem.Metadata != nil {
					for k, v := range blobItem.Metadata {
						if v != nil {
							blob.Metadata[k] = *v
						}
					}
				}
			} else {
				// Directory marker - no size
				blob.Size = 0
			}

			allBlobs = append(allBlobs, blob)
			seenItems[blobName] = true
		}
	}

	blobs := allBlobs

	return blobs, nil
}

// isImmediateChild checks if a blob path is an immediate child of the prefix
func isImmediateChild(blobPath, prefix string) bool {
	// Remove prefix from path
	if prefix != "" {
		if !strings.HasPrefix(blobPath, prefix) {
			return false
		}
		blobPath = blobPath[len(prefix):]
	}

	// Check if there's a path separator in the remaining path
	// If there is, it's not an immediate child
	return !strings.Contains(blobPath, "/")
}

// getDisplayName extracts the display name from a blob path, removing the prefix
func getDisplayName(blobPath, prefix string) string {
	if prefix != "" {
		if strings.HasPrefix(blobPath, prefix) {
			return blobPath[len(prefix):]
		}
	}
	return blobPath
}

// getParentFolderPath extracts the parent folder path for a nested blob
func getParentFolderPath(blobPath, currentPrefix string) string {
	// Remove current prefix
	relativePath := blobPath
	if currentPrefix != "" && strings.HasPrefix(blobPath, currentPrefix) {
		relativePath = blobPath[len(currentPrefix):]
	}

	// Find first slash to get immediate parent folder
	firstSlash := strings.Index(relativePath, "/")
	if firstSlash == -1 {
		return "" // No parent folder
	}

	// Return the folder path including the prefix
	folderName := relativePath[:firstSlash+1]
	if currentPrefix != "" {
		return currentPrefix + folderName
	}
	return folderName
}

// GetBlobDetails gets detailed information about a blob
func (c *Client) GetBlobDetails(ctx context.Context, subscriptionID, resourceGroupName, storageAccountName, containerName, blobName string) (*models.Blob, error) {
	// Get storage account keys
	keys, err := c.getStorageAccountKeys(ctx, subscriptionID, resourceGroupName, storageAccountName)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage account keys: %w", err)
	}

	if len(keys) == 0 {
		return nil, fmt.Errorf("no storage account keys found")
	}

	// Use the first key to create blob service client
	credential, err := azblob.NewSharedKeyCredential(storageAccountName, keys[0])
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", storageAccountName)
	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create blob client: %w", err)
	}

	// Get blob properties
	blobClient := client.ServiceClient().NewContainerClient(containerName).NewBlobClient(blobName)
	props, err := blobClient.GetProperties(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get blob properties: %w", err)
	}

	blob := &models.Blob{
		Name:        blobName,
		DisplayName: path.Base(blobName),
		Metadata:    make(map[string]string),
		IsDirectory: false,
	}

	if props.ContentLength != nil {
		blob.Size = *props.ContentLength
	}
	if props.ContentType != nil {
		blob.ContentType = *props.ContentType
	}
	if props.LastModified != nil {
		blob.LastModified = *props.LastModified
	}
	if props.ETag != nil {
		blob.ETag = string(*props.ETag)
	}

	if props.Metadata != nil {
		for k, v := range props.Metadata {
			if v != nil {
				blob.Metadata[k] = *v
			}
		}
	}

	return blob, nil
}

// getStorageAccountKeys retrieves the storage account keys
func (c *Client) getStorageAccountKeys(ctx context.Context, subscriptionID, resourceGroupName, storageAccountName string) ([]string, error) {
	client, err := armstorage.NewAccountsClient(subscriptionID, c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage accounts client: %w", err)
	}

	resp, err := client.ListKeys(ctx, resourceGroupName, storageAccountName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list storage account keys: %w", err)
	}

	var keys []string
	if resp.Keys != nil {
		for _, key := range resp.Keys {
			if key.Value != nil {
				keys = append(keys, *key.Value)
			}
		}
	}

	return keys, nil
}
