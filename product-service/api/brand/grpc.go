package brand

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	brandpb "github.com/people257/poor-guy-shop/product-service/gen/proto/proto/product/brand"
	"github.com/people257/poor-guy-shop/product-service/internal/application/brand"
)

// BrandServer 品牌gRPC服务器
type BrandServer struct {
	brandpb.UnimplementedBrandServiceServer
	brandService *brand.Service
}

// NewBrandServer 创建品牌gRPC服务器
func NewBrandServer(brandService *brand.Service) *BrandServer {
	return &BrandServer{
		brandService: brandService,
	}
}

// CreateBrand 创建品牌
func (s *BrandServer) CreateBrand(ctx context.Context, req *brandpb.CreateBrandReq) (*brandpb.CreateBrandResp, error) {
	dto := &brand.CreateBrandDTO{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		LogoURL:     req.LogoUrl,
		WebsiteURL:  req.WebsiteUrl,
		SortOrder:   int(req.SortOrder),
	}

	result, err := s.brandService.CreateBrand(ctx, dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建品牌失败: %v", err)
	}

	return &brandpb.CreateBrandResp{
		Brand: s.toBrandPB(result),
	}, nil
}

// UpdateBrand 更新品牌
func (s *BrandServer) UpdateBrand(ctx context.Context, req *brandpb.UpdateBrandReq) (*brandpb.UpdateBrandResp, error) {
	dto := &brand.UpdateBrandDTO{
		ID:          req.Id,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		LogoURL:     req.LogoUrl,
		WebsiteURL:  req.WebsiteUrl,
		SortOrder:   int(req.SortOrder),
		IsActive:    req.IsActive,
	}

	result, err := s.brandService.UpdateBrand(ctx, dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "更新品牌失败: %v", err)
	}

	return &brandpb.UpdateBrandResp{
		Brand: s.toBrandPB(result),
	}, nil
}

// DeleteBrand 删除品牌
func (s *BrandServer) DeleteBrand(ctx context.Context, req *brandpb.DeleteBrandReq) (*brandpb.DeleteBrandResp, error) {
	err := s.brandService.DeleteBrand(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "删除品牌失败: %v", err)
	}

	return &brandpb.DeleteBrandResp{
		Success: true,
	}, nil
}

// GetBrand 获取品牌详情
func (s *BrandServer) GetBrand(ctx context.Context, req *brandpb.GetBrandReq) (*brandpb.GetBrandResp, error) {
	result, err := s.brandService.GetBrand(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取品牌失败: %v", err)
	}

	return &brandpb.GetBrandResp{
		Brand: s.toBrandPB(result),
	}, nil
}

// ListBrands 获取品牌列表
func (s *BrandServer) ListBrands(ctx context.Context, req *brandpb.ListBrandsReq) (*brandpb.ListBrandsResp, error) {
	dto := &brand.ListBrandsDTO{
		Page:      int(req.Page),
		PageSize:  int(req.PageSize),
		IsActive:  &req.IsActive,
		Keyword:   req.Keyword,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	result, err := s.brandService.ListBrands(ctx, dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取品牌列表失败: %v", err)
	}

	brands := make([]*brandpb.Brand, len(result.Brands))
	for i, b := range result.Brands {
		brands[i] = s.toBrandPB(b)
	}

	return &brandpb.ListBrandsResp{
		Brands:   brands,
		Total:    result.Total,
		Page:     int32(result.Page),
		PageSize: int32(result.PageSize),
	}, nil
}

// toBrandPB 转换为protobuf品牌对象
func (s *BrandServer) toBrandPB(b *brand.BrandDTO) *brandpb.Brand {
	return &brandpb.Brand{
		Id:          b.ID,
		Name:        b.Name,
		Slug:        b.Slug,
		Description: b.Description,
		LogoUrl:     b.LogoURL,
		WebsiteUrl:  b.WebsiteURL,
		SortOrder:   int32(b.SortOrder),
		IsActive:    b.IsActive,
		CreatedAt:   parseTime(b.CreatedAt),
		UpdatedAt:   parseTime(b.UpdatedAt),
	}
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
