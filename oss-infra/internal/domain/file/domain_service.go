package file

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// DomainService 文件领域服务
type DomainService struct {
	maxFileSize      int64  // 最大文件大小（字节）
	allowedMimeTypes map[string]bool
	storagePrefix    string
}

// NewDomainService 创建文件领域服务
func NewDomainService(maxFileSize int64, allowedMimeTypes []string, storagePrefix string) *DomainService {
	mimeTypeMap := make(map[string]bool)
	for _, mt := range allowedMimeTypes {
		mimeTypeMap[mt] = true
	}
	
	return &DomainService{
		maxFileSize:      maxFileSize,
		allowedMimeTypes: mimeTypeMap,
		storagePrefix:    storagePrefix,
	}
}

// CreateFileReq 创建文件请求
type CreateFileReq struct {
	FileData   []byte
	Filename   string
	Category   string
	Visibility string
	OwnerID    string
}

// CreateFile 创建文件领域对象
func (s *DomainService) CreateFile(req *CreateFileReq) (*File, error) {
	// 验证输入参数
	if err := s.validateCreateFileReq(req); err != nil {
		return nil, err
	}
	
	// 生成文件哈希
	fileHash := s.generateFileHash(req.FileData)
	
	// 检测MIME类型
	mimeType := s.detectMimeType(req.Filename, req.FileData)
	if !s.isAllowedMimeType(mimeType) {
		return nil, ErrUnsupportedFileType
	}
	
	// 生成文件键
	fileKey := s.generateFileKey(req.Filename, req.Category, req.OwnerID)
	
	// 生成文件路径
	filePath := s.generateFilePath(fileKey)
	
	// 设置默认可见性
	visibility := req.Visibility
	if visibility == "" {
		visibility = "private"
	}
	
	// 创建文件实体
	file := NewFile(
		req.Filename,
		fileKey,
		filePath,
		int64(len(req.FileData)),
		mimeType,
		fileHash,
		req.Category,
		req.OwnerID,
		visibility,
	)
	
	// 验证文件实体
	if err := file.Validate(); err != nil {
		return nil, err
	}
	
	return file, nil
}

// CheckFileAccess 检查文件访问权限
func (s *DomainService) CheckFileAccess(file *File, userID string) (bool, error) {
	if file == nil {
		return false, ErrFileNotFound
	}
	
	// 检查文件是否已删除
	if !file.IsActive() {
		return false, ErrFileDeleted
	}
	
	// 公开文件或文件所有者可以访问
	return file.CanAccess(userID), nil
}

// ValidateFileSize 验证文件大小
func (s *DomainService) ValidateFileSize(fileSize int64) error {
	if fileSize <= 0 {
		return ErrInvalidFileSize
	}
	if fileSize > s.maxFileSize {
		return ErrFileSizeTooLarge
	}
	return nil
}

// ValidateFileType 验证文件类型
func (s *DomainService) ValidateFileType(mimeType string) error {
	if !s.isAllowedMimeType(mimeType) {
		return ErrUnsupportedFileType
	}
	return nil
}

// GenerateFileKey 生成文件存储键
func (s *DomainService) GenerateFileKey(filename, category, ownerID string) string {
	return s.generateFileKey(filename, category, ownerID)
}

// CalculateFileHash 计算文件哈希
func (s *DomainService) CalculateFileHash(data []byte) string {
	return s.generateFileHash(data)
}

// 私有方法

// validateCreateFileReq 验证创建文件请求
func (s *DomainService) validateCreateFileReq(req *CreateFileReq) error {
	if req.FileData == nil || len(req.FileData) == 0 {
		return ErrInvalidFileSize
	}
	
	if req.Filename == "" {
		return ErrInvalidFilename
	}
	
	if req.OwnerID == "" {
		return ErrInvalidOwnerID
	}
	
	// 验证文件大小
	if err := s.ValidateFileSize(int64(len(req.FileData))); err != nil {
		return err
	}
	
	// 验证可见性设置
	if req.Visibility != "" && req.Visibility != "public" && req.Visibility != "private" {
		return ErrInvalidVisibility
	}
	
	// 验证文件分类
	if req.Category != "" && !isValidCategory(req.Category) {
		return ErrInvalidCategory
	}
	
	return nil
}

// generateFileHash 生成文件SHA256哈希
func (s *DomainService) generateFileHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// detectMimeType 检测文件MIME类型
func (s *DomainService) detectMimeType(filename string, data []byte) string {
	// 首先根据文件扩展名检测
	ext := strings.ToLower(filepath.Ext(filename))
	if mimeType := mime.TypeByExtension(ext); mimeType != "" {
		return mimeType
	}
	
	// 根据文件头检测常见类型
	if len(data) >= 512 {
		return detectMimeTypeByHeader(data[:512])
	}
	
	// 默认返回二进制类型
	return "application/octet-stream"
}

// detectMimeTypeByHeader 根据文件头检测MIME类型
func detectMimeTypeByHeader(data []byte) string {
	if len(data) < 4 {
		return "application/octet-stream"
	}
	
	// 常见文件头签名检测
	signatures := map[string]string{
		"\xFF\xD8\xFF":                     "image/jpeg",
		"\x89PNG\r\n\x1A\n":               "image/png",
		"GIF87a":                          "image/gif",
		"GIF89a":                          "image/gif",
		"RIFF":                            "image/webp", // 需要进一步检查
		"\x00\x00\x00\x20ftypmp41":       "video/mp4",
		"\x00\x00\x00\x1CftypM4V":        "video/mp4",
		"%PDF":                            "application/pdf",
		"\xD0\xCF\x11\xE0\xA1\xB1\x1A\xE1": "application/msword",
		"PK\x03\x04":                      "application/zip", // 也可能是docx/xlsx等
	}
	
	dataStr := string(data)
	for sig, mimeType := range signatures {
		if strings.HasPrefix(dataStr, sig) {
			// 特殊处理RIFF格式
			if sig == "RIFF" && len(data) >= 12 && string(data[8:12]) == "WEBP" {
				return "image/webp"
			}
			return mimeType
		}
	}
	
	return "application/octet-stream"
}

// isAllowedMimeType 检查是否为允许的MIME类型
func (s *DomainService) isAllowedMimeType(mimeType string) bool {
	// 如果没有限制，则允许所有类型
	if len(s.allowedMimeTypes) == 0 {
		return true
	}
	
	return s.allowedMimeTypes[mimeType]
}

// generateFileKey 生成文件存储键
func (s *DomainService) generateFileKey(filename, category, ownerID string) string {
	// 生成UUID作为唯一标识
	fileID := uuid.New().String()
	
	// 获取文件扩展名
	ext := filepath.Ext(filename)
	
	// 构建路径: prefix/category/year/month/day/uuid.ext
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")
	day := now.Format("02")
	
	var pathParts []string
	
	// 添加存储前缀
	if s.storagePrefix != "" {
		pathParts = append(pathParts, s.storagePrefix)
	}
	
	// 添加分类
	if category != "" {
		pathParts = append(pathParts, category)
	} else {
		pathParts = append(pathParts, "general")
	}
	
	// 添加日期路径
	pathParts = append(pathParts, year, month, day)
	
	// 添加文件名
	fileName := fileID + ext
	pathParts = append(pathParts, fileName)
	
	return strings.Join(pathParts, "/")
}

// generateFilePath 生成文件存储路径
func (s *DomainService) generateFilePath(fileKey string) string {
	// 对于本地存储，返回完整路径
	// 对于云存储，返回相对路径
	return fmt.Sprintf("/data/oss/%s", fileKey)
}

// GetDefaultAllowedMimeTypes 获取默认允许的MIME类型
func GetDefaultAllowedMimeTypes() []string {
	return []string{
		// 图片
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/svg+xml",
		
		// 视频
		"video/mp4",
		"video/avi",
		"video/mov",
		"video/wmv",
		"video/flv",
		"video/webm",
		
		// 音频
		"audio/mpeg",
		"audio/wav",
		"audio/ogg",
		"audio/mp3",
		
		// 文档
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.ms-powerpoint",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"text/plain",
		"text/csv",
		
		// 压缩文件
		"application/zip",
		"application/x-rar-compressed",
		"application/x-7z-compressed",
		"application/gzip",
		
		// 其他
		"application/json",
		"application/xml",
		"text/xml",
	}
}
