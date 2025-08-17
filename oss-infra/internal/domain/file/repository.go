package file

import (
	"context"
)

// Repository 文件仓储接口
type Repository interface {
	// 基本CRUD操作
	Create(ctx context.Context, file *File) error
	GetByID(ctx context.Context, id string) (*File, error)
	GetByFileKey(ctx context.Context, fileKey string) (*File, error)
	GetByHash(ctx context.Context, hash string) (*File, error)
	Update(ctx context.Context, file *File) error
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) error
	
	// 查询操作
	List(ctx context.Context, req *ListFilesReq) ([]*File, int64, error)
	ListByOwner(ctx context.Context, ownerID string, req *ListFilesReq) ([]*File, int64, error)
	ListByCategory(ctx context.Context, category string, req *ListFilesReq) ([]*File, int64, error)
	
	// 统计操作
	CountByOwner(ctx context.Context, ownerID string) (int64, error)
	CountByCategory(ctx context.Context, category string) (int64, error)
	GetTotalSizeByOwner(ctx context.Context, ownerID string) (int64, error)
	
	// 批量操作
	BatchDelete(ctx context.Context, ids []string) error
	BatchUpdateVisibility(ctx context.Context, ids []string, visibility string) error
	
	// 清理操作
	CleanupDeletedFiles(ctx context.Context, beforeDate string) (int64, error)
	
	// 访问日志操作
	CreateAccessLog(ctx context.Context, log *FileAccessLog) error
	GetAccessLogs(ctx context.Context, req *GetAccessLogsReq) ([]*FileAccessLog, int64, error)
	CleanupAccessLogs(ctx context.Context, beforeDate string) (int64, error)
}

// StorageRepository 文件存储仓储接口
type StorageRepository interface {
	// 文件存储操作
	SaveFile(ctx context.Context, req *SaveFileReq) error
	GetFile(ctx context.Context, fileKey string) ([]byte, error)
	DeleteFile(ctx context.Context, fileKey string) error
	FileExists(ctx context.Context, fileKey string) (bool, error)
	
	// URL生成
	GenerateUploadURL(ctx context.Context, fileKey string, expiresIn int32) (string, error)
	GenerateDownloadURL(ctx context.Context, fileKey string, expiresIn int32) (string, error)
	
	// 文件操作
	CopyFile(ctx context.Context, srcKey, dstKey string) error
	MoveFile(ctx context.Context, srcKey, dstKey string) error
	
	// 元数据操作
	GetFileMetadata(ctx context.Context, fileKey string) (*FileMetadata, error)
	UpdateFileMetadata(ctx context.Context, fileKey string, metadata *FileMetadata) error
	
	// 批量操作
	BatchDeleteFiles(ctx context.Context, fileKeys []string) error
	
	// 存储统计
	GetStorageUsage(ctx context.Context, prefix string) (*StorageUsage, error)
}

// ListFilesReq 文件列表查询请求
type ListFilesReq struct {
	OwnerID    string
	Category   string
	Visibility string
	Status     *int16
	Page       int32
	PageSize   int32
	OrderBy    string // created_at, file_size, filename
	OrderDesc  bool
	Search     string // 搜索文件名
}

// GetAccessLogsReq 访问日志查询请求
type GetAccessLogsReq struct {
	FileID     string
	UserID     string
	Action     string
	StartTime  string
	EndTime    string
	Page       int32
	PageSize   int32
}

// SaveFileReq 保存文件请求
type SaveFileReq struct {
	FileKey  string
	FileData []byte
	MimeType string
	Metadata map[string]string
}

// FileMetadata 文件元数据
type FileMetadata struct {
	ContentType     string
	ContentLength   int64
	LastModified    string
	ETag            string
	CustomMetadata  map[string]string
}

// StorageUsage 存储使用情况
type StorageUsage struct {
	TotalSize   int64
	FileCount   int64
	LastUpdated string
}
