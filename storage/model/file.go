// Package model provides the data models for the application.
package model

import "strings"

const (
	// MetaAppend is the append for the metadata file
	MetaAppend = ".meta"
	// FileAppend is the append for the file
	FileAppend = ".file"
	// ImageAppend is the append for the image
	ImageAppend = ".image_"
	// FileMetaCachePrefix is the prefix for the file meta cache
	FileMetaCachePrefix = "file_meta_"
	// ImageCachePrefixWithWildcard is the prefix for the image cache with wildcard
	ImageCachePrefixWithWildcard = ".image_*"
	// FilePerm is the file permission, only the owner can read and write
	FilePerm = 0600
)

// FileHash contains the hash of the file
type FileHash struct {
	HashSha1   string `json:"sha1"`
	HashSha256 string `json:"sha256"`
	HashMd5    string `json:"md5"`
}

// FileMeta contains the metadata of the file
type FileMeta struct {
	RelativePath string   `json:"relativePath"` // The relative path to the base file storage path
	FileName     string   `json:"fileName"`     // The name of the file
	MimeType     string   `json:"mimeType"`     // The mime type of the file
	UploadedAt   int64    `json:"uploadedAt"`   // The time the file was uploaded
	Hash         FileHash `json:"hash"`         // The hash of the file
}

// FileInfo contains the information of the file
type FileInfo struct {
	FilePath     string `json:"filePath"`     // Relative path to the base file storage path
	FileSize     int64  `json:"fileSize"`     // The size of the file
	LastModified int64  `json:"lastModified"` // The last modified time of the file

	FileMeta FileMeta `json:"fileMeta"` // The metadata of the file
}

// IsImage checks if the file is an image
func (f *FileInfo) IsImage() bool {
	return strings.HasPrefix(f.FileMeta.MimeType, "image/")
}
