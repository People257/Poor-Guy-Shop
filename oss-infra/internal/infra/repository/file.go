package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/people257/poor-guy-shop/oss-infra/gen/gen/model"
	"github.com/people257/poor-guy-shop/oss-infra/gen/gen/query"
	"github.com/people257/poor-guy-shop/oss-infra/internal/domain/file"
)

// fileRepository 文件仓储实现
type fileRepository struct {
	db *gorm.DB
	q  *query.Query
}

// NewFileRepository 创建文件仓储
func NewFileRepository(db *gorm.DB) file.Repository {
	return &fileRepository{
		db: db,
		q:  query.Use(db),
	}
}

// Create 创建文件
func (r *fileRepository) Create(ctx context.Context, fileEntity *file.File) error {
	modelFile := r.domainToModel(fileEntity)
	err := r.q.File.WithContext(ctx).Create(modelFile)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	fileEntity.ID = modelFile.ID
	return nil
}

// GetByID 根据ID获取文件
func (r *fileRepository) GetByID(ctx context.Context, id string) (*file.File, error) {
	modelFile, err := r.q.File.WithContext(ctx).Where(r.q.File.ID.Eq(id)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, file.ErrFileNotFound
		}
		return nil, fmt.Errorf("获取文件失败: %w", err)
	}
	return r.modelToDomain(modelFile), nil
}

// GetByFileKey 根据文件键获取文件
func (r *fileRepository) GetByFileKey(ctx context.Context, fileKey string) (*file.File, error) {
	modelFile, err := r.q.File.WithContext(ctx).Where(r.q.File.FileKey.Eq(fileKey)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, file.ErrFileNotFound
		}
		return nil, fmt.Errorf("获取文件失败: %w", err)
	}
	return r.modelToDomain(modelFile), nil
}

// GetByHash 根据哈希获取文件
func (r *fileRepository) GetByHash(ctx context.Context, hash string) (*file.File, error) {
	modelFile, err := r.q.File.WithContext(ctx).Where(r.q.File.FileHash.Eq(hash)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, file.ErrFileNotFound
		}
		return nil, fmt.Errorf("获取文件失败: %w", err)
	}
	return r.modelToDomain(modelFile), nil
}

// Update 更新文件
func (r *fileRepository) Update(ctx context.Context, fileEntity *file.File) error {
	modelFile := r.domainToModel(fileEntity)
	_, err := r.q.File.WithContext(ctx).Where(r.q.File.ID.Eq(fileEntity.ID)).Updates(modelFile)
	if err != nil {
		return fmt.Errorf("更新文件失败: %w", err)
	}
	return nil
}

// Delete 物理删除文件
func (r *fileRepository) Delete(ctx context.Context, id string) error {
	_, err := r.q.File.WithContext(ctx).Where(r.q.File.ID.Eq(id)).Delete()
	if err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}
	return nil
}

// SoftDelete 软删除文件
func (r *fileRepository) SoftDelete(ctx context.Context, id string) error {
	now := time.Now()
	_, err := r.q.File.WithContext(ctx).Where(r.q.File.ID.Eq(id)).Updates(map[string]interface{}{
		"status":     2,
		"updated_at": now,
	})
	if err != nil {
		return fmt.Errorf("软删除文件失败: %w", err)
	}
	return nil
}

// List 获取文件列表
func (r *fileRepository) List(ctx context.Context, req *file.ListFilesReq) ([]*file.File, int64, error) {
	query := r.q.File.WithContext(ctx)

	// 构建查询条件
	if req.OwnerID != "" {
		query = query.Where(r.q.File.OwnerID.Eq(req.OwnerID))
	}
	if req.Category != "" {
		query = query.Where(r.q.File.Category.Eq(&req.Category))
	}
	if req.Visibility != "" {
		query = query.Where(r.q.File.Visibility.Eq(req.Visibility))
	}
	if req.Status != nil {
		query = query.Where(r.q.File.Status.Eq(*req.Status))
	} else {
		// 默认只查询正常状态的文件
		query = query.Where(r.q.File.Status.Eq(1))
	}
	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		query = query.Where(r.q.File.Filename.Like(searchPattern))
	}

	// 排序
	orderBy := req.OrderBy
	if orderBy == "" {
		orderBy = "created_at"
	}

	switch orderBy {
	case "created_at":
		if req.OrderDesc {
			query = query.Order(r.q.File.CreatedAt.Desc())
		} else {
			query = query.Order(r.q.File.CreatedAt.Asc())
		}
	case "file_size":
		if req.OrderDesc {
			query = query.Order(r.q.File.FileSize.Desc())
		} else {
			query = query.Order(r.q.File.FileSize.Asc())
		}
	case "filename":
		if req.OrderDesc {
			query = query.Order(r.q.File.Filename.Desc())
		} else {
			query = query.Order(r.q.File.Filename.Asc())
		}
	default:
		query = query.Order(r.q.File.CreatedAt.Desc())
	}

	// 分页
	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(int(offset)).Limit(int(req.PageSize))

	// 查询数据
	modelFiles, err := query.Find()
	if err != nil {
		return nil, 0, fmt.Errorf("查询文件列表失败: %w", err)
	}

	// 查询总数
	countQuery := r.q.File.WithContext(ctx)
	if req.OwnerID != "" {
		countQuery = countQuery.Where(r.q.File.OwnerID.Eq(req.OwnerID))
	}
	if req.Category != "" {
		countQuery = countQuery.Where(r.q.File.Category.Eq(&req.Category))
	}
	if req.Visibility != "" {
		countQuery = countQuery.Where(r.q.File.Visibility.Eq(req.Visibility))
	}
	if req.Status != nil {
		countQuery = countQuery.Where(r.q.File.Status.Eq(*req.Status))
	} else {
		countQuery = countQuery.Where(r.q.File.Status.Eq(1))
	}
	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		countQuery = countQuery.Where(r.q.File.Filename.Like(searchPattern))
	}

	total, err := countQuery.Count()
	if err != nil {
		return nil, 0, fmt.Errorf("查询文件总数失败: %w", err)
	}

	// 转换为领域对象
	files := make([]*file.File, 0, len(modelFiles))
	for _, modelFile := range modelFiles {
		files = append(files, r.modelToDomain(modelFile))
	}

	return files, total, nil
}

// ListByOwner 根据所有者获取文件列表
func (r *fileRepository) ListByOwner(ctx context.Context, ownerID string, req *file.ListFilesReq) ([]*file.File, int64, error) {
	req.OwnerID = ownerID
	return r.List(ctx, req)
}

// ListByCategory 根据分类获取文件列表
func (r *fileRepository) ListByCategory(ctx context.Context, category string, req *file.ListFilesReq) ([]*file.File, int64, error) {
	req.Category = category
	return r.List(ctx, req)
}

// CountByOwner 统计用户文件数量
func (r *fileRepository) CountByOwner(ctx context.Context, ownerID string) (int64, error) {
	count, err := r.q.File.WithContext(ctx).Where(
		r.q.File.OwnerID.Eq(ownerID),
		r.q.File.Status.Eq(1),
	).Count()
	if err != nil {
		return 0, fmt.Errorf("统计用户文件数量失败: %w", err)
	}
	return count, nil
}

// CountByCategory 统计分类文件数量
func (r *fileRepository) CountByCategory(ctx context.Context, category string) (int64, error) {
	count, err := r.q.File.WithContext(ctx).Where(
		r.q.File.Category.Eq(&category),
		r.q.File.Status.Eq(1),
	).Count()
	if err != nil {
		return 0, fmt.Errorf("统计分类文件数量失败: %w", err)
	}
	return count, nil
}

// GetTotalSizeByOwner 获取用户文件总大小
func (r *fileRepository) GetTotalSizeByOwner(ctx context.Context, ownerID string) (int64, error) {
	var totalSize int64
	err := r.db.WithContext(ctx).Model(&model.File{}).
		Select("COALESCE(SUM(file_size), 0)").
		Where("owner_id = ? AND status = ?", ownerID, 1).
		Scan(&totalSize).Error
	if err != nil {
		return 0, fmt.Errorf("查询用户文件总大小失败: %w", err)
	}
	return totalSize, nil
}

// BatchDelete 批量删除文件
func (r *fileRepository) BatchDelete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := r.q.File.WithContext(ctx).Where(r.q.File.ID.In(ids...)).Delete()
	if err != nil {
		return fmt.Errorf("批量删除文件失败: %w", err)
	}
	return nil
}

// BatchUpdateVisibility 批量更新文件可见性
func (r *fileRepository) BatchUpdateVisibility(ctx context.Context, ids []string, visibility string) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := r.q.File.WithContext(ctx).Where(r.q.File.ID.In(ids...)).Updates(map[string]interface{}{
		"visibility": visibility,
		"updated_at": time.Now(),
	})
	if err != nil {
		return fmt.Errorf("批量更新文件可见性失败: %w", err)
	}
	return nil
}

// CleanupDeletedFiles 清理已删除的文件
func (r *fileRepository) CleanupDeletedFiles(ctx context.Context, beforeDate string) (int64, error) {
	result, err := r.q.File.WithContext(ctx).Where(
		r.q.File.Status.Eq(2),
		r.q.File.UpdatedAt.Lt(parseTime(beforeDate)),
	).Delete()
	if err != nil {
		return 0, fmt.Errorf("清理已删除文件失败: %w", err)
	}
	return result.RowsAffected, nil
}

// CreateAccessLog 创建访问日志
func (r *fileRepository) CreateAccessLog(ctx context.Context, log *file.FileAccessLog) error {
	modelLog := &model.FileAccessLog{
		FileID:     log.FileID,
		UserID:     log.UserID,
		Action:     log.Action,
		IPAddress:  log.IPAddress,
		StatusCode: log.StatusCode,
		AccessedAt: log.AccessedAt,
	}

	err := r.q.FileAccessLog.WithContext(ctx).Create(modelLog)
	if err != nil {
		return fmt.Errorf("创建访问日志失败: %w", err)
	}
	log.ID = modelLog.ID
	return nil
}

// GetAccessLogs 获取访问日志
func (r *fileRepository) GetAccessLogs(ctx context.Context, req *file.GetAccessLogsReq) ([]*file.FileAccessLog, int64, error) {
	query := r.q.FileAccessLog.WithContext(ctx)

	// 构建查询条件
	if req.FileID != "" {
		query = query.Where(r.q.FileAccessLog.FileID.Eq(req.FileID))
	}
	if req.UserID != "" {
		query = query.Where(r.q.FileAccessLog.UserID.Eq(&req.UserID))
	}
	if req.Action != "" {
		query = query.Where(r.q.FileAccessLog.Action.Eq(req.Action))
	}
	if req.StartTime != "" {
		query = query.Where(r.q.FileAccessLog.AccessedAt.Gte(parseTime(req.StartTime)))
	}
	if req.EndTime != "" {
		query = query.Where(r.q.FileAccessLog.AccessedAt.Lte(parseTime(req.EndTime)))
	}

	// 排序
	query = query.Order(r.q.FileAccessLog.AccessedAt.Desc())

	// 分页
	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(int(offset)).Limit(int(req.PageSize))

	// 查询数据
	modelLogs, err := query.Find()
	if err != nil {
		return nil, 0, fmt.Errorf("查询访问日志失败: %w", err)
	}

	// 查询总数
	countQuery := r.q.FileAccessLog.WithContext(ctx)
	if req.FileID != "" {
		countQuery = countQuery.Where(r.q.FileAccessLog.FileID.Eq(req.FileID))
	}
	if req.UserID != "" {
		countQuery = countQuery.Where(r.q.FileAccessLog.UserID.Eq(&req.UserID))
	}
	if req.Action != "" {
		countQuery = countQuery.Where(r.q.FileAccessLog.Action.Eq(req.Action))
	}
	if req.StartTime != "" {
		countQuery = countQuery.Where(r.q.FileAccessLog.AccessedAt.Gte(parseTime(req.StartTime)))
	}
	if req.EndTime != "" {
		countQuery = countQuery.Where(r.q.FileAccessLog.AccessedAt.Lte(parseTime(req.EndTime)))
	}

	total, err := countQuery.Count()
	if err != nil {
		return nil, 0, fmt.Errorf("查询访问日志总数失败: %w", err)
	}

	// 转换为领域对象
	logs := make([]*file.FileAccessLog, 0, len(modelLogs))
	for _, modelLog := range modelLogs {
		logs = append(logs, &file.FileAccessLog{
			ID:         modelLog.ID,
			FileID:     modelLog.FileID,
			UserID:     modelLog.UserID,
			Action:     modelLog.Action,
			IPAddress:  modelLog.IPAddress,
			StatusCode: modelLog.StatusCode,
			AccessedAt: modelLog.AccessedAt,
		})
	}

	return logs, total, nil
}

// CleanupAccessLogs 清理访问日志
func (r *fileRepository) CleanupAccessLogs(ctx context.Context, beforeDate string) (int64, error) {
	result, err := r.q.FileAccessLog.WithContext(ctx).Where(
		r.q.FileAccessLog.AccessedAt.Lt(parseTime(beforeDate)),
	).Delete()
	if err != nil {
		return 0, fmt.Errorf("清理访问日志失败: %w", err)
	}
	return result.RowsAffected, nil
}

// 私有方法

// domainToModel 领域对象转换为数据模型
func (r *fileRepository) domainToModel(fileEntity *file.File) *model.File {
	var category *string
	if fileEntity.Category != "" {
		category = &fileEntity.Category
	}

	return &model.File{
		ID:         fileEntity.ID,
		Filename:   fileEntity.Filename,
		FileKey:    fileEntity.FileKey,
		FilePath:   fileEntity.FilePath,
		FileSize:   fileEntity.FileSize,
		MimeType:   fileEntity.MimeType,
		FileHash:   fileEntity.FileHash,
		Category:   category,
		OwnerID:    fileEntity.OwnerID,
		Visibility: fileEntity.Visibility,
		Status:     fileEntity.Status,
		CreatedAt:  fileEntity.CreatedAt,
		UpdatedAt:  fileEntity.UpdatedAt,
	}
}

// modelToDomain 数据模型转换为领域对象
func (r *fileRepository) modelToDomain(modelFile *model.File) *file.File {
	category := ""
	if modelFile.Category != nil {
		category = *modelFile.Category
	}

	return &file.File{
		ID:         modelFile.ID,
		Filename:   modelFile.Filename,
		FileKey:    modelFile.FileKey,
		FilePath:   modelFile.FilePath,
		FileSize:   modelFile.FileSize,
		MimeType:   modelFile.MimeType,
		FileHash:   modelFile.FileHash,
		Category:   category,
		OwnerID:    modelFile.OwnerID,
		Visibility: modelFile.Visibility,
		Status:     modelFile.Status,
		CreatedAt:  modelFile.CreatedAt,
		UpdatedAt:  modelFile.UpdatedAt,
	}
}

// parseTime 解析时间字符串
func parseTime(timeStr string) time.Time {
	formats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t
		}
	}

	return time.Now()
}
