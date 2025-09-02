package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/people257/poor-guy-shop/oss-infra/cmd/grpc/internal/config"
	"github.com/people257/poor-guy-shop/oss-infra/internal/domain/file"
)

// AliyunOSSStorage 阿里云OSS存储实现
type AliyunOSSStorage struct {
	client *oss.Client
	bucket *oss.Bucket
	config *config.AliyunStorageConfig
}

// NewAliyunOSSStorage 创建阿里云OSS存储
func NewAliyunOSSStorage(cfg *config.AliyunStorageConfig) (file.StorageRepository, error) {
	client, err := oss.New(cfg.Endpoint, cfg.AccessKeyID, cfg.SecretAccessKey)
	if err != nil {
		return nil, fmt.Errorf("创建OSS客户端失败: %w", err)
	}

	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("获取OSS Bucket失败: %w", err)
	}

	return &AliyunOSSStorage{
		client: client,
		bucket: bucket,
		config: cfg,
	}, nil
}

// SaveFile 保存文件
func (s *AliyunOSSStorage) SaveFile(ctx context.Context, req *file.SaveFileReq) error {
	options := []oss.Option{}

	// 设置内容类型
	if req.MimeType != "" {
		options = append(options, oss.ContentType(req.MimeType))
	}

	// 设置自定义元数据
	for key, value := range req.Metadata {
		options = append(options, oss.Meta(key, value))
	}

	err := s.bucket.PutObject(req.FileKey, strings.NewReader(string(req.FileData)), options...)
	if err != nil {
		return fmt.Errorf("上传文件到OSS失败: %w", err)
	}

	return nil
}

// GetFile 获取文件内容
func (s *AliyunOSSStorage) GetFile(ctx context.Context, fileKey string) ([]byte, error) {
	body, err := s.bucket.GetObject(fileKey)
	if err != nil {
		return nil, fmt.Errorf("从OSS获取文件失败: %w", err)
	}
	defer body.Close()

	data := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := body.Read(buf)
		if n > 0 {
			data = append(data, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	return data, nil
}

// DeleteFile 删除文件
func (s *AliyunOSSStorage) DeleteFile(ctx context.Context, fileKey string) error {
	err := s.bucket.DeleteObject(fileKey)
	if err != nil {
		return fmt.Errorf("从OSS删除文件失败: %w", err)
	}
	return nil
}

// FileExists 检查文件是否存在
func (s *AliyunOSSStorage) FileExists(ctx context.Context, fileKey string) (bool, error) {
	_, err := s.bucket.GetObjectMeta(fileKey)
	if err != nil {
		if ossErr, ok := err.(oss.ServiceError); ok && ossErr.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("检查文件存在性失败: %w", err)
	}
	return true, nil
}

// GenerateUploadURL 生成上传URL
func (s *AliyunOSSStorage) GenerateUploadURL(ctx context.Context, fileKey string, expiresIn int32) (string, error) {
	expiration := time.Duration(expiresIn) * time.Second
	url, err := s.bucket.SignURL(fileKey, oss.HTTPPut, int64(expiration.Seconds()))
	if err != nil {
		return "", fmt.Errorf("生成上传URL失败: %w", err)
	}
	return url, nil
}

// GenerateDownloadURL 生成下载URL
func (s *AliyunOSSStorage) GenerateDownloadURL(ctx context.Context, fileKey string, expiresIn int32) (string, error) {
	expiration := time.Duration(expiresIn) * time.Second
	url, err := s.bucket.SignURL(fileKey, oss.HTTPGet, int64(expiration.Seconds()))
	if err != nil {
		return "", fmt.Errorf("生成下载URL失败: %w", err)
	}
	return url, nil
}

// CopyFile 复制文件
func (s *AliyunOSSStorage) CopyFile(ctx context.Context, srcKey, dstKey string) error {
	_, err := s.bucket.CopyObject(srcKey, dstKey)
	if err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}
	return nil
}

// MoveFile 移动文件
func (s *AliyunOSSStorage) MoveFile(ctx context.Context, srcKey, dstKey string) error {
	// OSS不支持直接移动，需要先复制再删除
	if err := s.CopyFile(ctx, srcKey, dstKey); err != nil {
		return fmt.Errorf("移动文件失败(复制阶段): %w", err)
	}

	if err := s.DeleteFile(ctx, srcKey); err != nil {
		return fmt.Errorf("移动文件失败(删除阶段): %w", err)
	}

	return nil
}

// GetFileMetadata 获取文件元数据
func (s *AliyunOSSStorage) GetFileMetadata(ctx context.Context, fileKey string) (*file.FileMetadata, error) {
	meta, err := s.bucket.GetObjectDetailedMeta(fileKey)
	if err != nil {
		return nil, fmt.Errorf("获取文件元数据失败: %w", err)
	}

	metadata := &file.FileMetadata{
		ContentType:    meta.Get("Content-Type"),
		ContentLength:  0, // 需要从Content-Length解析
		LastModified:   meta.Get("Last-Modified"),
		ETag:           meta.Get("Etag"),
		CustomMetadata: make(map[string]string),
	}

	// 提取自定义元数据
	for key, values := range meta {
		if strings.HasPrefix(key, "X-Oss-Meta-") {
			customKey := strings.TrimPrefix(key, "X-Oss-Meta-")
			if len(values) > 0 {
				metadata.CustomMetadata[customKey] = values[0]
			}
		}
	}

	return metadata, nil
}

// UpdateFileMetadata 更新文件元数据
func (s *AliyunOSSStorage) UpdateFileMetadata(ctx context.Context, fileKey string, metadata *file.FileMetadata) error {
	options := []oss.Option{}

	if metadata.ContentType != "" {
		options = append(options, oss.ContentType(metadata.ContentType))
	}

	for key, value := range metadata.CustomMetadata {
		options = append(options, oss.Meta(key, value))
	}

	// OSS需要通过ModifyObjectMeta来更新元数据
	err := s.bucket.SetObjectMeta(fileKey, options...)
	if err != nil {
		return fmt.Errorf("更新文件元数据失败: %w", err)
	}

	return nil
}

// BatchDeleteFiles 批量删除文件
func (s *AliyunOSSStorage) BatchDeleteFiles(ctx context.Context, fileKeys []string) error {
	if len(fileKeys) == 0 {
		return nil
	}

	// OSS支持批量删除，最多1000个
	const batchSize = 1000
	for i := 0; i < len(fileKeys); i += batchSize {
		end := i + batchSize
		if end > len(fileKeys) {
			end = len(fileKeys)
		}

		batch := fileKeys[i:end]
		_, err := s.bucket.DeleteObjects(batch)
		if err != nil {
			return fmt.Errorf("批量删除文件失败: %w", err)
		}
	}

	return nil
}

// GetStorageUsage 获取存储使用情况
func (s *AliyunOSSStorage) GetStorageUsage(ctx context.Context, prefix string) (*file.StorageUsage, error) {
	// 简化实现，返回基本信息
	// 实际实现需要遍历对象并统计
	return &file.StorageUsage{
		TotalFiles: 0,
		TotalSize:  0,
		UsedSpace:  0,
	}, nil
}
