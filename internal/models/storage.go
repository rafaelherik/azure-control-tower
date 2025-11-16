package models

import "time"

// Container represents a storage container
type Container struct {
	Name         string
	LastModified time.Time
	ETag         string
	PublicAccess string
	Metadata     map[string]string
}

// Blob represents a blob in a storage container
type Blob struct {
	Name         string
	DisplayName  string // Display name (without prefix path)
	Size         int64
	ContentType  string
	LastModified time.Time
	ETag         string
	Metadata     map[string]string
	IsDirectory  bool
}
