package file

import (
	"context"
	"fmt"
	"time"

	"github.com/people257/poor-guy-shop/oss-infra/internal/domain/file"
)

// Service 文件应用服务
type Service struct {
	fileRepo    file.Repository
	fileDomain  *file.DomainService
	storageRepo file.StorageRepository
}

// NewService 创建文件应用服务
func NewService(
	fileRepo file.Repository,
	fileDomain *file.DomainService,
	storageRepo file.StorageRepository,
) *Service {
	return &Service{
		fileRepo:    fileRepo,
		fileDomain:  fileDomain,
		storageRepo: storageRepo,
	}
}

// UploadFileReq 上传文件请求
type UploadFileReq struct {
	FileData   []byte
	Filename   string
	Category   string
	Visibility string
	UserID     string // 从认证中间件获取
}

// UploadFileResp 上传文件响应
type UploadFileResp struct {
	FileInfo *FileInfoDTO
}

// GetDownloadURLReq 获取下载URL请求
type GetDownloadURLReq struct {
	FileID    string
	ExpiresIn int32
	UserID    string // 从认证中间件获取
}

// GetDownloadURLResp 获取下载URL响应
type GetDownloadURLResp struct {
	DownloadURL string
	ExpiresIn   int32
}

// GetFileInfoReq 获取文件信息请求
type GetFileInfoReq struct {
	FileID string
	UserID string // 从认证中间件获取
}

// GetFileInfoResp 获取文件信息响应
type GetFileInfoResp struct {
	FileInfo *FileInfoDTO
}

// ListFilesReq 获取文件列表请求
type ListFilesReq struct {
	Category string
	Page     int32
	PageSize int32
	UserID   string // 从认证中间件获取
}

// ListFilesResp 获取文件列表响应
type ListFilesResp struct {
	Files    []*FileInfoDTO
	Total    int64
	Page     int32
	PageSize int32
}

// DeleteFileReq 删除文件请求
type DeleteFileReq struct {
	FileID string
	UserID string // 从认证中间件获取
}

// CheckFileAccessReq 检查文件访问权限请求
type CheckFileAccessReq struct {
	FileID string
	UserID string
}

// CheckFileAccessResp 检查文件访问权限响应
type CheckFileAccessResp struct {
	HasAccess bool
	Reason    string
}

// FileInfoDTO 文件信息DTO
type FileInfoDTO struct {
	FileID      string
	Filename    string
	FileKey     string
	FileSize    int64
	MimeType    string
	Category    string
	OwnerID     string
	Visibility  string
	CreatedAt   time.Time
	DownloadURL string
}

// UploadFile 上传文件
func (s *Service) UploadFile(ctx context.Context, req *UploadFileReq) (*UploadFileResp, error) {
	// 创建文件领域对象
	fileEntity, err := s.fileDomain.CreateFile(&file.CreateFileReq{
		FileData:   req.FileData,
		Filename:   req.Filename,
		Category:   req.Category,
		Visibility: req.Visibility,
		OwnerID:    req.UserID,
	})
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}

	// 存储文件到存储系统
	err = s.storageRepo.SaveFile(ctx, &file.SaveFileReq{
		FileKey:  fileEntity.FileKey,
		FileData: req.FileData,
		MimeType: fileEntity.MimeType,
	})
	if err != nil {
		return nil, fmt.Errorf("存储文件失败: %w", err)
	}

	// 保存文件信息到数据库
	err = s.fileRepo.Create(ctx, fileEntity)
	if err != nil {
		// 如果数据库保存失败，删除已上传的文件
		_ = s.storageRepo.DeleteFile(ctx, fileEntity.FileKey)
		return nil, fmt.Errorf("保存文件信息失败: %w", err)
	}

	// 记录访问日志
	_ = s.fileRepo.CreateAccessLog(ctx, &file.FileAccessLog{
		FileID:     fileEntity.ID,
		UserID:     &req.UserID,
		Action:     "upload",
		IPAddress:  "", // 从context获取
		StatusCode: func() *int16 { code := int16(200); return &code }(),
		AccessedAt: time.Now(),
	})

	// 生成下载URL
	downloadURL, err := s.storageRepo.GenerateDownloadURL(ctx, fileEntity.FileKey, 3600) // 1小时
	if err != nil {
		downloadURL = "" // 忽略错误，不影响上传
	}

	return &UploadFileResp{
		FileInfo: &FileInfoDTO{
			FileID:      fileEntity.ID,
			Filename:    fileEntity.Filename,
			FileKey:     fileEntity.FileKey,
			FileSize:    fileEntity.FileSize,
			MimeType:    fileEntity.MimeType,
			Category:    fileEntity.Category,
			OwnerID:     fileEntity.OwnerID,
			Visibility:  fileEntity.Visibility,
			CreatedAt:   fileEntity.CreatedAt,
			DownloadURL: downloadURL,
		},
	}, nil
}

// GetDownloadURL 获取文件下载URL
func (s *Service) GetDownloadURL(ctx context.Context, req *GetDownloadURLReq) (*GetDownloadURLResp, error) {
	// 检查文件是否存在和权限
	fileEntity, err := s.fileRepo.GetByID(ctx, req.FileID)
	if err != nil {
		return nil, fmt.Errorf("文件不存在: %w", err)
	}

	// 检查访问权限
	hasAccess, err := s.fileDomain.CheckFileAccess(fileEntity, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("权限检查失败: %w", err)
	}
	if !hasAccess {
		return nil, fmt.Errorf("无权限访问该文件")
	}

	// 生成下载URL
	expiresIn := req.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = 3600 // 默认1小时
	}

	downloadURL, err := s.storageRepo.GenerateDownloadURL(ctx, fileEntity.FileKey, expiresIn)
	if err != nil {
		return nil, fmt.Errorf("生成下载URL失败: %w", err)
	}

	// 记录访问日志
	_ = s.fileRepo.CreateAccessLog(ctx, &file.FileAccessLog{
		FileID:     fileEntity.ID,
		UserID:     &req.UserID,
		Action:     "download",
		IPAddress:  "", // 从context获取
		StatusCode: func() *int16 { code := int16(200); return &code }(),
		AccessedAt: time.Now(),
	})

	return &GetDownloadURLResp{
		DownloadURL: downloadURL,
		ExpiresIn:   expiresIn,
	}, nil
}

// GetFileInfo 获取文件信息
func (s *Service) GetFileInfo(ctx context.Context, req *GetFileInfoReq) (*GetFileInfoResp, error) {
	// 获取文件信息
	fileEntity, err := s.fileRepo.GetByID(ctx, req.FileID)
	if err != nil {
		return nil, fmt.Errorf("文件不存在: %w", err)
	}

	// 检查访问权限
	hasAccess, err := s.fileDomain.CheckFileAccess(fileEntity, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("权限检查失败: %w", err)
	}
	if !hasAccess {
		return nil, fmt.Errorf("无权限访问该文件")
	}

	// 生成临时下载URL（如果有权限）
	var downloadURL string
	if hasAccess {
		downloadURL, _ = s.storageRepo.GenerateDownloadURL(ctx, fileEntity.FileKey, 3600)
	}

	return &GetFileInfoResp{
		FileInfo: &FileInfoDTO{
			FileID:      fileEntity.ID,
			Filename:    fileEntity.Filename,
			FileKey:     fileEntity.FileKey,
			FileSize:    fileEntity.FileSize,
			MimeType:    fileEntity.MimeType,
			Category:    fileEntity.Category,
			OwnerID:     fileEntity.OwnerID,
			Visibility:  fileEntity.Visibility,
			CreatedAt:   fileEntity.CreatedAt,
			DownloadURL: downloadURL,
		},
	}, nil
}

// ListFiles 获取文件列表
func (s *Service) ListFiles(ctx context.Context, req *ListFilesReq) (*ListFilesResp, error) {
	// 获取文件列表
	files, total, err := s.fileRepo.List(ctx, &file.ListFilesReq{
		OwnerID:  req.UserID,
		Category: req.Category,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("获取文件列表失败: %w", err)
	}

	// 转换为DTO
	fileDTOs := make([]*FileInfoDTO, 0, len(files))
	for _, f := range files {
		// 为每个文件生成临时下载URL
		downloadURL, _ := s.storageRepo.GenerateDownloadURL(ctx, f.FileKey, 3600)

		fileDTOs = append(fileDTOs, &FileInfoDTO{
			FileID:      f.ID,
			Filename:    f.Filename,
			FileKey:     f.FileKey,
			FileSize:    f.FileSize,
			MimeType:    f.MimeType,
			Category:    f.Category,
			OwnerID:     f.OwnerID,
			Visibility:  f.Visibility,
			CreatedAt:   f.CreatedAt,
			DownloadURL: downloadURL,
		})
	}

	return &ListFilesResp{
		Files:    fileDTOs,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// DeleteFile 删除文件
func (s *Service) DeleteFile(ctx context.Context, req *DeleteFileReq) error {
	// 获取文件信息
	fileEntity, err := s.fileRepo.GetByID(ctx, req.FileID)
	if err != nil {
		return fmt.Errorf("文件不存在: %w", err)
	}

	// 检查删除权限（只有所有者可以删除）
	if fileEntity.OwnerID != req.UserID {
		return fmt.Errorf("无权限删除该文件")
	}

	// 软删除文件记录
	err = s.fileRepo.SoftDelete(ctx, req.FileID)
	if err != nil {
		return fmt.Errorf("删除文件记录失败: %w", err)
	}

	// 从存储系统删除文件
	err = s.storageRepo.DeleteFile(ctx, fileEntity.FileKey)
	if err != nil {
		// 存储删除失败，记录日志但不回滚数据库操作
		// 可以通过定时任务清理孤立文件
	}

	// 记录访问日志
	_ = s.fileRepo.CreateAccessLog(ctx, &file.FileAccessLog{
		FileID:     fileEntity.ID,
		UserID:     &req.UserID,
		Action:     "delete",
		IPAddress:  "", // 从context获取
		StatusCode: func() *int16 { code := int16(200); return &code }(),
		AccessedAt: time.Now(),
	})

	return nil
}

// CheckFileAccess 检查文件访问权限
func (s *Service) CheckFileAccess(ctx context.Context, req *CheckFileAccessReq) (*CheckFileAccessResp, error) {
	// 获取文件信息
	fileEntity, err := s.fileRepo.GetByID(ctx, req.FileID)
	if err != nil {
		return &CheckFileAccessResp{
			HasAccess: false,
			Reason:    "文件不存在",
		}, nil
	}

	// 检查访问权限
	hasAccess, err := s.fileDomain.CheckFileAccess(fileEntity, req.UserID)
	if err != nil {
		return &CheckFileAccessResp{
			HasAccess: false,
			Reason:    "权限检查失败",
		}, nil
	}

	reason := ""
	if !hasAccess {
		if fileEntity.Visibility == "private" && fileEntity.OwnerID != req.UserID {
			reason = "文件为私有，且非文件所有者"
		} else if fileEntity.Status != 1 {
			reason = "文件已被删除"
		}
	}

	return &CheckFileAccessResp{
		HasAccess: hasAccess,
		Reason:    reason,
	}, nil
}
