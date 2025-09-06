package category

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	categorypb "github.com/people257/poor-guy-shop/product-service/gen/proto/proto/product/category"
	"github.com/people257/poor-guy-shop/product-service/internal/application/category"
)

// CategoryServer 分类gRPC服务器
type CategoryServer struct {
	categorypb.UnimplementedCategoryServiceServer
	categoryService *category.Service
}

// NewCategoryServer 创建分类gRPC服务器
func NewCategoryServer(categoryService *category.Service) *CategoryServer {
	return &CategoryServer{
		categoryService: categoryService,
	}
}

// CreateCategory 创建分类
func (s *CategoryServer) CreateCategory(ctx context.Context, req *categorypb.CreateCategoryReq) (*categorypb.CreateCategoryResp, error) {
	dto := &category.CreateCategoryDTO{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ParentID:    &req.ParentId,
		SortOrder:   int(req.SortOrder),
		IconURL:     req.IconUrl,
		BannerURL:   req.BannerUrl,
	}

	if req.ParentId == "" {
		dto.ParentID = nil
	}

	result, err := s.categoryService.CreateCategory(ctx, dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建分类失败: %v", err)
	}

	return &categorypb.CreateCategoryResp{
		Category: s.toCategoryPB(result),
	}, nil
}

// UpdateCategory 更新分类
func (s *CategoryServer) UpdateCategory(ctx context.Context, req *categorypb.UpdateCategoryReq) (*categorypb.UpdateCategoryResp, error) {
	dto := &category.UpdateCategoryDTO{
		ID:          req.Id,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ParentID:    &req.ParentId,
		SortOrder:   int(req.SortOrder),
		IconURL:     req.IconUrl,
		BannerURL:   req.BannerUrl,
		IsActive:    req.IsActive,
	}

	if req.ParentId == "" {
		dto.ParentID = nil
	}

	result, err := s.categoryService.UpdateCategory(ctx, dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "更新分类失败: %v", err)
	}

	return &categorypb.UpdateCategoryResp{
		Category: s.toCategoryPB(result),
	}, nil
}

// DeleteCategory 删除分类
func (s *CategoryServer) DeleteCategory(ctx context.Context, req *categorypb.DeleteCategoryReq) (*categorypb.DeleteCategoryResp, error) {
	err := s.categoryService.DeleteCategory(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "删除分类失败: %v", err)
	}

	return &categorypb.DeleteCategoryResp{
		Success: true,
	}, nil
}

// GetCategory 获取分类详情
func (s *CategoryServer) GetCategory(ctx context.Context, req *categorypb.GetCategoryReq) (*categorypb.GetCategoryResp, error) {
	result, err := s.categoryService.GetCategory(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取分类失败: %v", err)
	}

	return &categorypb.GetCategoryResp{
		Category: s.toCategoryPB(result),
	}, nil
}

// ListCategories 获取分类列表
func (s *CategoryServer) ListCategories(ctx context.Context, req *categorypb.ListCategoriesReq) (*categorypb.ListCategoriesResp, error) {
	dto := &category.ListCategoriesDTO{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
		ParentID: &req.ParentId,
		Level: func() *int {
			if req.Level != 0 {
				v := int(req.Level)
				return &v
			}
			return nil
		}(),
		IsActive:  &req.IsActive,
		Keyword:   req.Keyword,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	if req.ParentId == "" {
		dto.ParentID = nil
	}
	if req.Level == 0 {
		dto.Level = nil
	}

	result, err := s.categoryService.ListCategories(ctx, dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取分类列表失败: %v", err)
	}

	categories := make([]*categorypb.Category, len(result.Categories))
	for i, c := range result.Categories {
		categories[i] = s.toCategoryPB(c)
	}

	return &categorypb.ListCategoriesResp{
		Categories: categories,
		Total:      result.Total,
		Page:       int32(result.Page),
		PageSize:   int32(result.PageSize),
	}, nil
}

// GetCategoryTree 获取分类树
func (s *CategoryServer) GetCategoryTree(ctx context.Context, req *categorypb.GetCategoryTreeReq) (*categorypb.GetCategoryTreeResp, error) {
	result, err := s.categoryService.GetCategoryTree(ctx, req.ActiveOnly)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取分类树失败: %v", err)
	}

	categories := make([]*categorypb.Category, len(result))
	for i, c := range result {
		categories[i] = s.toCategoryPB(c)
	}

	return &categorypb.GetCategoryTreeResp{
		Categories: categories,
	}, nil
}

// toCategoryPB 转换为protobuf分类对象
func (s *CategoryServer) toCategoryPB(c *category.CategoryDTO) *categorypb.Category {
	pb := &categorypb.Category{
		Id:          c.ID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		Level:       int32(c.Level),
		SortOrder:   int32(c.SortOrder),
		IconUrl:     c.IconURL,
		BannerUrl:   c.BannerURL,
		IsActive:    c.IsActive,
		CreatedAt:   parseTime(c.CreatedAt),
		UpdatedAt:   parseTime(c.UpdatedAt),
	}

	if c.ParentID != nil {
		pb.ParentId = *c.ParentID
	}

	// 转换子分类
	if len(c.Children) > 0 {
		pb.Children = make([]*categorypb.Category, len(c.Children))
		for i, child := range c.Children {
			pb.Children[i] = s.toCategoryPB(child)
		}
	}

	return pb
}

// parseTime 解析时间字符串为protobuf时间戳
func parseTime(timeStr string) *timestamppb.Timestamp {
	if timeStr == "" {
		return nil
	}

	// 尝试多种时间格式
	formats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z07:00",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return timestamppb.New(t)
		}
	}

	return nil
}
