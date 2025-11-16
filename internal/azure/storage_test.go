package azure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsImmediateChild(t *testing.T) {
	tests := []struct {
		name     string
		blobPath string
		prefix   string
		expected bool
	}{
		{
			name:     "Root level file",
			blobPath: "file.txt",
			prefix:   "",
			expected: true,
		},
		{
			name:     "Root level folder",
			blobPath: "folder/",
			prefix:   "",
			expected: false, // Has a slash, so not an immediate child
		},
		{
			name:     "Nested file from root",
			blobPath: "folder/file.txt",
			prefix:   "",
			expected: false,
		},
		{
			name:     "Immediate child file in folder",
			blobPath: "folder/file.txt",
			prefix:   "folder/",
			expected: true,
		},
		{
			name:     "Immediate child subfolder",
			blobPath: "folder/subfolder/",
			prefix:   "folder/",
			expected: false, // Has a slash after prefix removal, so not immediate
		},
		{
			name:     "Nested file in subfolder",
			blobPath: "folder/subfolder/file.txt",
			prefix:   "folder/",
			expected: false,
		},
		{
			name:     "File two levels deep",
			blobPath: "folder/subfolder/file.txt",
			prefix:   "folder/subfolder/",
			expected: true,
		},
		{
			name:     "File three levels deep from root",
			blobPath: "a/b/c/file.txt",
			prefix:   "",
			expected: false,
		},
		{
			name:     "File in different folder",
			blobPath: "otherfolder/file.txt",
			prefix:   "folder/",
			expected: false,
		},
		{
			name:     "Empty blob path",
			blobPath: "",
			prefix:   "",
			expected: true,
		},
		{
			name:     "Prefix longer than blob path",
			blobPath: "file.txt",
			prefix:   "folder/subfolder/",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isImmediateChild(tt.blobPath, tt.prefix)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		blobPath string
		prefix   string
		expected string
	}{
		{
			name:     "Root level file",
			blobPath: "file.txt",
			prefix:   "",
			expected: "file.txt",
		},
		{
			name:     "Root level folder",
			blobPath: "folder/",
			prefix:   "",
			expected: "folder/",
		},
		{
			name:     "File with prefix removed",
			blobPath: "folder/file.txt",
			prefix:   "folder/",
			expected: "file.txt",
		},
		{
			name:     "Subfolder with prefix removed",
			blobPath: "folder/subfolder/",
			prefix:   "folder/",
			expected: "subfolder/",
		},
		{
			name:     "Deep path with prefix",
			blobPath: "a/b/c/file.txt",
			prefix:   "a/b/",
			expected: "c/file.txt",
		},
		{
			name:     "No prefix match",
			blobPath: "otherfolder/file.txt",
			prefix:   "folder/",
			expected: "otherfolder/file.txt",
		},
		{
			name:     "Empty prefix",
			blobPath: "folder/subfolder/file.txt",
			prefix:   "",
			expected: "folder/subfolder/file.txt",
		},
		{
			name:     "Exact match (directory marker)",
			blobPath: "folder/",
			prefix:   "folder/",
			expected: "",
		},
		{
			name:     "Path with special characters",
			blobPath: "folder/my-file_2024.txt",
			prefix:   "folder/",
			expected: "my-file_2024.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getDisplayName(tt.blobPath, tt.prefix)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetParentFolderPath(t *testing.T) {
	tests := []struct {
		name          string
		blobPath      string
		currentPrefix string
		expected      string
	}{
		{
			name:          "Root level nested file",
			blobPath:      "folder/file.txt",
			currentPrefix: "",
			expected:      "folder/",
		},
		{
			name:          "Root level deeply nested file",
			blobPath:      "folder/subfolder/file.txt",
			currentPrefix: "",
			expected:      "folder/",
		},
		{
			name:          "File in subfolder",
			blobPath:      "folder/subfolder/file.txt",
			currentPrefix: "folder/",
			expected:      "folder/subfolder/",
		},
		{
			name:          "Deeply nested from subfolder",
			blobPath:      "folder/subfolder/deep/file.txt",
			currentPrefix: "folder/",
			expected:      "folder/subfolder/",
		},
		{
			name:          "Multiple levels",
			blobPath:      "a/b/c/d/file.txt",
			currentPrefix: "a/b/",
			expected:      "a/b/c/",
		},
		{
			name:          "No parent (immediate child)",
			blobPath:      "folder/file.txt",
			currentPrefix: "folder/",
			expected:      "",
		},
		{
			name:          "Empty prefix with nested path",
			blobPath:      "a/b/file.txt",
			currentPrefix: "",
			expected:      "a/",
		},
		{
			name:          "No slash in remaining path",
			blobPath:      "file.txt",
			currentPrefix: "",
			expected:      "",
		},
		{
			name:          "Folder marker",
			blobPath:      "folder/subfolder/",
			currentPrefix: "folder/",
			expected:      "folder/subfolder/", // Returns the folder itself since it's a directory marker
		},
		{
			name:          "Deep folder marker from root",
			blobPath:      "folder/subfolder/deep/",
			currentPrefix: "",
			expected:      "folder/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getParentFolderPath(tt.blobPath, tt.currentPrefix)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetParentFolderPath_EdgeCases(t *testing.T) {
	t.Run("Very deep nesting", func(t *testing.T) {
		blobPath := "a/b/c/d/e/f/g/h/file.txt"
		prefix := "a/b/c/"
		expected := "a/b/c/d/"

		result := getParentFolderPath(blobPath, prefix)
		assert.Equal(t, expected, result)
	})

	t.Run("Path with special characters", func(t *testing.T) {
		blobPath := "folder-1/sub_folder_2/my-file.txt"
		prefix := "folder-1/"
		expected := "folder-1/sub_folder_2/"

		result := getParentFolderPath(blobPath, prefix)
		assert.Equal(t, expected, result)
	})

	t.Run("Unicode characters in path", func(t *testing.T) {
		blobPath := "フォルダ/サブフォルダ/ファイル.txt"
		prefix := "フォルダ/"
		expected := "フォルダ/サブフォルダ/"

		result := getParentFolderPath(blobPath, prefix)
		assert.Equal(t, expected, result)
	})
}

func TestStorageHelpers_Integration(t *testing.T) {
	// Test the helpers work together correctly for hierarchical navigation
	t.Run("Navigate folder hierarchy", func(t *testing.T) {
		// Simulate blobs in a container
		blobs := []string{
			"file1.txt",
			"file2.txt",
			"folder1/file3.txt",
			"folder1/subfolder1/file4.txt",
			"folder1/subfolder2/file5.txt",
			"folder2/file6.txt",
		}

		// At root level, we should see:
		// - file1.txt (immediate child)
		// - file2.txt (immediate child)
		// - folder1/ (folder)
		// - folder2/ (folder)
		rootChildren := make([]string, 0)
		for _, blob := range blobs {
			if isImmediateChild(blob, "") {
				rootChildren = append(rootChildren, getDisplayName(blob, ""))
			}
		}
		assert.Contains(t, rootChildren, "file1.txt")
		assert.Contains(t, rootChildren, "file2.txt")

		// At folder1/ level, we should see:
		// - file3.txt (immediate child)
		// - subfolder1/ (folder)
		// - subfolder2/ (folder)
		folder1Children := make([]string, 0)
		for _, blob := range blobs {
			if isImmediateChild(blob, "folder1/") {
				folder1Children = append(folder1Children, getDisplayName(blob, "folder1/"))
			}
		}
		assert.Contains(t, folder1Children, "file3.txt")

		// Files in subfolders are not immediate children
		assert.NotContains(t, folder1Children, "subfolder1/file4.txt")
	})
}

func BenchmarkIsImmediateChild(b *testing.B) {
	blobPath := "folder/subfolder/deep/file.txt"
	prefix := "folder/subfolder/"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = isImmediateChild(blobPath, prefix)
	}
}

func BenchmarkGetDisplayName(b *testing.B) {
	blobPath := "folder/subfolder/deep/file.txt"
	prefix := "folder/subfolder/"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getDisplayName(blobPath, prefix)
	}
}

func BenchmarkGetParentFolderPath(b *testing.B) {
	blobPath := "folder/subfolder/deep/file.txt"
	prefix := "folder/"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getParentFolderPath(blobPath, prefix)
	}
}
