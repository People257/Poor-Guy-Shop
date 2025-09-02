package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/people257/poor-guy-shop/oss-infra/cmd/grpc/internal/config"
	"github.com/people257/poor-guy-shop/oss-infra/internal/domain/file"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// QiniuKodoStorage 七牛云Kodo存储实现
type QiniuKodoStorage struct {
	mac           *qbox.Mac
	cfg           *storage.Config
	bucketManager *storage.BucketManager
	config        *config.QiniuStorageConfig
}

// NewQiniuKodoStorage 创建七牛云存储
func NewQiniuKodoStorage(cfg *config.QiniuStorageConfig) (file.StorageRepository, error) {
	mac := qbox.NewMac(cfg.AccessKey, cfg.SecretKey)

	// 根据区域设置配置
	storageCfg := &storage.Config{
		UseHTTPS:      true,
		UseCdnDomains: false,
	}

	// 设置存储区域
	switch cfg.Zone {
	case "z0":
		storageCfg.Zone = &storage.ZoneHuadong
	case "z1":
		storageCfg.Zone = &storage.ZoneHuabei
	case "z2":
		storageCfg.Zone = &storage.ZoneHuanan
	case "na0":
		storageCfg.Zone = &storage.ZoneBeimei
	case "as0":
		storageCfg.Zone = &storage.ZoneXinjiapo
	default:
		storageCfg.Zone = &storage.ZoneHuadong // 默认华东
	}

	bucketManager := storage.NewBucketManager(mac, storageCfg)

	return &QiniuKodoStorage{
		mac:           mac,
		cfg:           storageCfg,
		bucketManager: bucketManager,
		config:        cfg,
	}, nil
}

// SaveFile 保存文件
func (s *QiniuKodoStorage) SaveFile(ctx context.Context, req *file.SaveFileReq) error {
	putPolicy := storage.PutPolicy{
		Scope: s.config.Bucket,
	}

	// 设置自定义元数据
	if len(req.Metadata) > 0 {
		putPolicy.PersistentOps = s.buildMetadataOps(req.Metadata)
	}

	upToken := putPolicy.UploadToken(s.mac)
	formUploader := storage.NewFormUploader(s.cfg)
	ret := storage.PutRet{}

	putExtra := storage.PutExtra{
		Params: map[string]string{},
	}

	// 设置MIME类型
	if req.MimeType != "" {
		putExtra.MimeType = req.MimeType
	}

	err := formUploader.Put(ctx, &ret, upToken, req.FileKey, strings.NewReader(string(req.FileData)), int64(len(req.FileData)), &putExtra)
	if err != nil {
		return fmt.Errorf("上传文件到七牛云失败: %w", err)
	}

	return nil
}

// GetFile 获取文件内容
func (s *QiniuKodoStorage) GetFile(ctx context.Context, fileKey string) ([]byte, error) {
	// 七牛云不提供直接获取文件内容的API，需要通过下载URL获取
	downloadURL, err := s.GenerateDownloadURL(ctx, fileKey, 3600)
	if err != nil {
		return nil, fmt.Errorf("生成下载URL失败: %w", err)
	}

	// 这里需要实现HTTP下载逻辑
	// 为了简化，暂时返回错误
	return nil, fmt.Errorf("七牛云不支持直接获取文件内容，请使用下载URL: %s", downloadURL)
}

// DeleteFile 删除文件
func (s *QiniuKodoStorage) DeleteFile(ctx context.Context, fileKey string) error {
	err := s.bucketManager.Delete(s.config.Bucket, fileKey)
	if err != nil {
		return fmt.Errorf("从七牛云删除文件失败: %w", err)
	}
	return nil
}

// FileExists 检查文件是否存在
func (s *QiniuKodoStorage) FileExists(ctx context.Context, fileKey string) (bool, error) {
	_, err := s.bucketManager.Stat(s.config.Bucket, fileKey)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return false, nil
		}
		return false, fmt.Errorf("检查文件存在性失败: %w", err)
	}
	return true, nil
}

// GenerateUploadURL 生成上传URL
func (s *QiniuKodoStorage) GenerateUploadURL(ctx context.Context, fileKey string, expiresIn int32) (string, error) {
	putPolicy := storage.PutPolicy{
		Scope:   fmt.Sprintf("%s:%s", s.config.Bucket, fileKey),
		Expires: uint64(time.Now().Add(time.Duration(expiresIn) * time.Second).Unix()),
	}
	upToken := putPolicy.UploadToken(s.mac)

	// 返回上传表单信息，实际使用时需要客户端构建上传请求
	return fmt.Sprintf("token:%s", upToken), nil
}

// GenerateDownloadURL 生成下载URL
func (s *QiniuKodoStorage) GenerateDownloadURL(ctx context.Context, fileKey string, expiresIn int32) (string, error) {
	domain := s.config.Domain
	if domain == "" {
		return "", fmt.Errorf("未配置七牛云域名")
	}

	publicURL := storage.MakePublicURL(domain, fileKey)

	// 如果需要私有访问，生成带签名的URL
	deadline := time.Now().Add(time.Duration(expiresIn) * time.Second).Unix()
	privateURL := storage.MakePrivateURL(s.mac, domain, fileKey, deadline)

	// 这里假设使用私有访问，实际应该根据bucket配置决定
	return privateURL, nil
}

// CopyFile 复制文件
func (s *QiniuKodoStorage) CopyFile(ctx context.Context, srcKey, dstKey string) error {
	err := s.bucketManager.Copy(s.config.Bucket, srcKey, s.config.Bucket, dstKey, true)
	if err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}
	return nil
}

// MoveFile 移动文件
func (s *QiniuKodoStorage) MoveFile(ctx context.Context, srcKey, dstKey string) error {
	err := s.bucketManager.Move(s.config.Bucket, srcKey, s.config.Bucket, dstKey, true)
	if err != nil {
		return fmt.Errorf("移动文件失败: %w", err)
	}
	return nil
}

// GetFileMetadata 获取文件元数据
func (s *QiniuKodoStorage) GetFileMetadata(ctx context.Context, fileKey string) (*file.FileMetadata, error) {
	fileInfo, err := s.bucketManager.Stat(s.config.Bucket, fileKey)
	if err != nil {
		return nil, fmt.Errorf("获取文件元数据失败: %w", err)
	}

	metadata := &file.FileMetadata{
		ContentType:    fileInfo.MimeType,
		ContentLength:  fileInfo.Fsize,
		LastModified:   time.Unix(fileInfo.PutTime/10000000, 0).Format(time.RFC3339),
		ETag:           fileInfo.Hash,
		CustomMetadata: make(map[string]string),
	}

	return metadata, nil
}

// UpdateFileMetadata 更新文件元数据
func (s *QiniuKodoStorage) UpdateFileMetadata(ctx context.Context, fileKey string, metadata *file.FileMetadata) error {
	// 七牛云不支持直接更新元数据，需要通过其他方式实现
	return fmt.Errorf("七牛云暂不支持直接更新文件元数据")
}

// BatchDeleteFiles 批量删除文件
func (s *QiniuKodoStorage) BatchDeleteFiles(ctx context.Context, fileKeys []string) error {
	if len(fileKeys) == 0 {
		return nil
	}

	// 构建批量删除操作
	deleteOps := make([]string, len(fileKeys))
	for i, key := range fileKeys {
		deleteOps[i] = storage.URIDelete(s.config.Bucket, key)
	}

	rets, err := s.bucketManager.Batch(deleteOps)
	if err != nil {
		return fmt.Errorf("批量删除文件失败: %w", err)
	}

	// 检查每个操作的结果
	for i, ret := range rets {
		if ret.Code != 200 {
			return fmt.Errorf("删除文件 %s 失败: %s", fileKeys[i], ret.Error)
		}
	}

	return nil
}

// GetStorageUsage 获取存储使用情况
func (s *QiniuKodoStorage) GetStorageUsage(ctx context.Context, prefix string) (*file.StorageUsage, error) {
	// 简化实现，返回基本信息
	// 实际实现需要通过API获取bucket统计信息
	return &file.StorageUsage{
		TotalFiles: 0,
		TotalSize:  0,
		UsedSpace:  0,
	}, nil
}

// buildMetadataOps 构建元数据操作
func (s *QiniuKodoStorage) buildMetadataOps(metadata map[string]string) string {
	// 七牛云通过persistent ops设置元数据
	// 这里简化处理，实际需要根据七牛云API文档构建
	return ""
}
