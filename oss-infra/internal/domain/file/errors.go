package file

import "errors"

// 领域错误定义
var (
	// 文件相关错误
	ErrFileNotFound     = errors.New("文件不存在")
	ErrFileAlreadyExists = errors.New("文件已存在")
	ErrFileDeleted      = errors.New("文件已被删除")
	ErrFileAccessDenied = errors.New("文件访问被拒绝")
	ErrFileUploadFailed = errors.New("文件上传失败")
	ErrFileDownloadFailed = errors.New("文件下载失败")
	ErrFileStorageFailed = errors.New("文件存储失败")
	
	// 验证相关错误
	ErrInvalidFilename   = errors.New("无效的文件名")
	ErrInvalidFileKey    = errors.New("无效的文件键")
	ErrInvalidFileSize   = errors.New("无效的文件大小")
	ErrInvalidMimeType   = errors.New("无效的MIME类型")
	ErrInvalidFileHash   = errors.New("无效的文件哈希")
	ErrInvalidOwnerID    = errors.New("无效的所有者ID")
	ErrInvalidVisibility = errors.New("无效的可见性设置")
	ErrInvalidCategory   = errors.New("无效的文件分类")
	ErrInvalidUserID     = errors.New("无效的用户ID")
	
	// 权限相关错误
	ErrPermissionDenied   = errors.New("权限不足")
	ErrNotFileOwner      = errors.New("非文件所有者")
	ErrPrivateFileAccess = errors.New("私有文件访问被拒绝")
	
	// 业务逻辑错误
	ErrFileSizeTooLarge   = errors.New("文件大小超出限制")
	ErrUnsupportedFileType = errors.New("不支持的文件类型")
	ErrQuotaExceeded      = errors.New("存储配额已超出")
	ErrDuplicateFileHash  = errors.New("重复的文件哈希")
	
	// 存储相关错误
	ErrStorageUnavailable = errors.New("存储服务不可用")
	ErrStorageQuotaFull   = errors.New("存储空间已满")
	ErrStorageCorrupted   = errors.New("存储文件损坏")
	ErrStorageTimeout     = errors.New("存储操作超时")
	
	// URL相关错误
	ErrURLGenerationFailed = errors.New("URL生成失败")
	ErrURLExpired          = errors.New("URL已过期")
	ErrInvalidURL          = errors.New("无效的URL")
)
