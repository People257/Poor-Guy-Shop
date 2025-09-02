package file

import (
	"context"
	"fmt"
	"mime"
	"path/filepath"
	"time"

	"github.com/google/uuid"
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
	// 参数验证
	if len(req.FileData) == 0 {
		return nil, fmt.Errorf("文件数据不能为空")
	}
	if req.Filename == "" {
		return nil, fmt.Errorf("文件名不能为空")
	}

	// 生成文件ID和文件键
	fileID := uuid.New().String()
	fileKey := s.generateFileKey(req.UserID, fileID, req.Filename)

	// 检测MIME类型
	mimeType := mime.TypeByExtension(filepath.Ext(req.Filename))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// 设置默认值
	if req.Category == "" {
		req.Category = "general"
	}
	if req.Visibility == "" {
		req.Visibility = "private"
	}

	// 文件大小验证
	err := s.fileDomain.ValidateFileSize(int64(len(req.FileData)))
	if err != nil {
		return nil, fmt.Errorf("文件大小验证失败: %w", err)
	}

	// 文件类型验证
	err = s.fileDomain.ValidateFileType(mimeType)
	if err != nil {
		return nil, fmt.Errorf("文件类型验证失败: %w", err)
	}

	// 上传到存储系统
	err = s.storageRepo.SaveFile(ctx, &file.SaveFileReq{
		FileKey:  fileKey,
		FileData: req.FileData,
		MimeType: mimeType,
		Metadata: map[string]string{
			"category":   req.Category,
			"visibility": req.Visibility,
			"owner_id":   req.UserID,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("文件上传失败: %w", err)
	}

	// 创建文件实体
	fileEntity := &file.File{
		ID:         fileID,
		Filename:   req.Filename,
		FileKey:    fileKey,
		FileSize:   int64(len(req.FileData)),
		MimeType:   mimeType,
		Category:   req.Category,
		OwnerID:    req.UserID,
		Visibility: req.Visibility,
		Status:     1, // 正常状态
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 保存文件记录
	err = s.fileRepo.Create(ctx, fileEntity)
	if err != nil {
		// 如果数据库保存失败，尝试清理已上传的文件
		_ = s.storageRepo.DeleteFile(ctx, fileKey)
		return nil, fmt.Errorf("保存文件记录失败: %w", err)
	}

	// 生成临时下载URL
	downloadURL, _ := s.storageRepo.GenerateDownloadURL(ctx, fileKey, 3600)

	// 记录访问日志
	_ = s.fileRepo.CreateAccessLog(ctx, &file.FileAccessLog{
		FileID:     fileEntity.ID,
		UserID:     &req.UserID,
		Action:     "upload",
		IPAddress:  "", // 从context获取
		StatusCode: func() *int16 { code := int16(200); return &code }(),
		AccessedAt: time.Now(),
	})

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

	// 设置默认过期时间
	expiresIn := req.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = 3600 // 默认1小时
	}

	// 生成下载URL
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

// generateFileKey 生成文件存储键
func (s *Service) generateFileKey(userID, fileID, filename string) string {
	ext := filepath.Ext(filename)
	timestamp := time.Now().Format("20060102")
	return fmt.Sprintf("uploads/%s/%s/%s%s", userID, timestamp, fileID, ext)
}
