package file

import (
	"fmt"
	"time"
)

// File 文件领域实体
type File struct {
	ID         string
	Filename   string
	FileKey    string
	FilePath   string
	FileSize   int64
	MimeType   string
	FileHash   string
	Category   string
	OwnerID    string
	Visibility string
	Status     int16
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// FileAccessLog 文件访问日志实体
type FileAccessLog struct {
	ID         string
	FileID     string
	UserID     *string
	Action     string
	IPAddress  string
	StatusCode *int16
	AccessedAt time.Time
}

// NewFile 创建新的文件实体
func NewFile(
	filename, fileKey, filePath string,
	fileSize int64,
	mimeType, fileHash, category, ownerID, visibility string,
) *File {
	now := time.Now()
	return &File{
		Filename:   filename,
		FileKey:    fileKey,
		FilePath:   filePath,
		FileSize:   fileSize,
		MimeType:   mimeType,
		FileHash:   fileHash,
		Category:   category,
		OwnerID:    ownerID,
		Visibility: visibility,
		Status:     1, // 正常状态
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// IsPublic 判断文件是否为公开的
func (f *File) IsPublic() bool {
	return f.Visibility == "public"
}

// IsOwner 判断是否为文件所有者
func (f *File) IsOwner(userID string) bool {
	return f.OwnerID == userID
}

// IsActive 判断文件是否为活跃状态
func (f *File) IsActive() bool {
	return f.Status == 1
}

// CanAccess 判断用户是否可以访问文件
func (f *File) CanAccess(userID string) bool {
	// 文件必须是活跃状态
	if !f.IsActive() {
		return false
	}
	
	// 公开文件或者是文件所有者
	return f.IsPublic() || f.IsOwner(userID)
}

// SoftDelete 软删除文件
func (f *File) SoftDelete() {
	f.Status = 2
	f.UpdatedAt = time.Now()
}

// UpdateVisibility 更新文件可见性
func (f *File) UpdateVisibility(visibility string) error {
	if visibility != "public" && visibility != "private" {
		return ErrInvalidVisibility
	}
	f.Visibility = visibility
	f.UpdatedAt = time.Now()
	return nil
}

// GetDisplayName 获取显示名称
func (f *File) GetDisplayName() string {
	if f.Filename != "" {
		return f.Filename
	}
	return f.FileKey
}

// IsImage 判断是否为图片文件
func (f *File) IsImage() bool {
	switch f.MimeType {
	case "image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp", "image/svg+xml":
		return true
	default:
		return false
	}
}

// IsVideo 判断是否为视频文件
func (f *File) IsVideo() bool {
	switch f.MimeType {
	case "video/mp4", "video/avi", "video/mov", "video/wmv", "video/flv", "video/webm":
		return true
	default:
		return false
	}
}

// IsDocument 判断是否为文档文件
func (f *File) IsDocument() bool {
	switch f.MimeType {
	case "application/pdf", "application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		 "application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		 "application/vnd.ms-powerpoint", "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		 "text/plain", "text/csv":
		return true
	default:
		return false
	}
}

// GetFileType 获取文件类型分类
func (f *File) GetFileType() string {
	if f.IsImage() {
		return "image"
	}
	if f.IsVideo() {
		return "video"
	}
	if f.IsDocument() {
		return "document"
	}
	return "other"
}

// GetSizeDisplay 获取文件大小的友好显示
func (f *File) GetSizeDisplay() string {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
	)

	size := float64(f.FileSize)
	switch {
	case f.FileSize >= GB:
		return fmt.Sprintf("%.2f GB", size/GB)
	case f.FileSize >= MB:
		return fmt.Sprintf("%.2f MB", size/MB)
	case f.FileSize >= KB:
		return fmt.Sprintf("%.2f KB", size/KB)
	default:
		return fmt.Sprintf("%d B", f.FileSize)
	}
}

// Validate 验证文件实体
func (f *File) Validate() error {
	if f.Filename == "" {
		return ErrInvalidFilename
	}
	if f.FileKey == "" {
		return ErrInvalidFileKey
	}
	if f.FileSize <= 0 {
		return ErrInvalidFileSize
	}
	if f.MimeType == "" {
		return ErrInvalidMimeType
	}
	if f.FileHash == "" {
		return ErrInvalidFileHash
	}
	if f.OwnerID == "" {
		return ErrInvalidOwnerID
	}
	if f.Visibility != "public" && f.Visibility != "private" {
		return ErrInvalidVisibility
	}
	if f.Category != "" && !isValidCategory(f.Category) {
		return ErrInvalidCategory
	}
	return nil
}

// isValidCategory 验证文件分类是否有效
func isValidCategory(category string) bool {
	validCategories := map[string]bool{
		"avatar":   true,
		"product":  true,
		"document": true,
		"temp":     true,
	}
	return validCategories[category]
}
