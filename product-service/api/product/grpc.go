package product

import (
	"context"
	"time"

	"github.com/people257/poor-guy-shop/common/auth"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	productpb "github.com/people257/poor-guy-shop/product-service/gen/proto/proto/product/product"
	"github.com/people257/poor-guy-shop/product-service/internal/application/product"
	productdomain "github.com/people257/poor-guy-shop/product-service/internal/domain/product"
)

// ProductServer 商品gRPC服务器
type ProductServer struct {
	productpb.UnimplementedProductServiceServer
	productService *product.Service
}

// NewProductServer 创建商品gRPC服务器
func NewProductServer(productService *product.Service) *ProductServer {
	return &ProductServer{
		productService: productService,
	}
}

// CreateProduct 创建商品
func (s *ProductServer) CreateProduct(ctx context.Context, req *productpb.CreateProductReq) (*productpb.CreateProductResp, error) {
	// 从认证上下文获取用户ID
	userID := auth.UserIDFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "用户未认证")
	}

	// 解析价格
	marketPrice, err := decimal.NewFromString(req.MarketPrice)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "市场价格格式无效: %v", err)
	}

	salePrice, err := decimal.NewFromString(req.SalePrice)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "销售价格格式无效: %v", err)
	}

	costPrice, err := decimal.NewFromString(req.CostPrice)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "成本价格格式无效: %v", err)
	}

	// 转换标签
	tags := make([]string, len(req.Tags))
	for i, tag := range req.Tags {
		tags[i] = tag
	}

	// 转换规格参数
	specifications := make(map[string]any)
	for k, v := range req.Specifications {
		specifications[k] = v
	}

	dto := &product.CreateProductDTO{
		Name:             req.Name,
		Slug:             req.Slug,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		CategoryID:       req.CategoryId,
		BrandID:          &req.BrandId,
		MarketPrice:      marketPrice,
		SalePrice:        salePrice,
		CostPrice:        costPrice,
		MainImageURL:     req.MainImageUrl,
		ImageURLs:        req.ImageUrls,
		VideoURL:         req.VideoUrl,
		Tags:             tags,
		Specifications:   specifications,
		IsFeatured:       req.IsFeatured,
		IsVirtual:        req.IsVirtual,
		SEOTitle:         req.SeoTitle,
		SEODescription:   req.SeoDescription,
		SEOKeywords:      req.SeoKeywords,
		SortOrder:        int(req.SortOrder),
	}

	if req.BrandId == "" {
		dto.BrandID = nil
	}

	result, err := s.productService.CreateProduct(ctx, dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建商品失败: %v", err)
	}

	return &productpb.CreateProductResp{
		Product: s.toProductPB(result),
	}, nil
}

// UpdateProduct 更新商品
func (s *ProductServer) UpdateProduct(ctx context.Context, req *productpb.UpdateProductReq) (*productpb.UpdateProductResp, error) {
	// 从认证上下文获取用户ID
	userID := auth.UserIDFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "用户未认证")
	}

	// 解析价格
	marketPrice, err := decimal.NewFromString(req.MarketPrice)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "市场价格格式无效: %v", err)
	}

	salePrice, err := decimal.NewFromString(req.SalePrice)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "销售价格格式无效: %v", err)
	}

	costPrice, err := decimal.NewFromString(req.CostPrice)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "成本价格格式无效: %v", err)
	}

	// 转换标签
	tags := make([]string, len(req.Tags))
	for i, tag := range req.Tags {
		tags[i] = tag
	}

	// 转换规格参数
	specifications := make(map[string]any)
	for k, v := range req.Specifications {
		specifications[k] = v
	}

	dto := &product.UpdateProductDTO{
		ID:               req.Id,
		Name:             req.Name,
		Slug:             req.Slug,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		CategoryID:       req.CategoryId,
		BrandID:          &req.BrandId,
		MarketPrice:      marketPrice,
		SalePrice:        salePrice,
		CostPrice:        costPrice,
		MainImageURL:     req.MainImageUrl,
		ImageURLs:        req.ImageUrls,
		VideoURL:         req.VideoUrl,
		Tags:             tags,
		Specifications:   specifications,
		Status:           productdomain.ProductStatus(req.Status),
		IsFeatured:       req.IsFeatured,
		IsVirtual:        req.IsVirtual,
		SEOTitle:         req.SeoTitle,
		SEODescription:   req.SeoDescription,
		SEOKeywords:      req.SeoKeywords,
		SortOrder:        int(req.SortOrder),
	}

	if req.BrandId == "" {
		dto.BrandID = nil
	}

	result, err := s.productService.UpdateProduct(ctx, dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "更新商品失败: %v", err)
	}

	return &productpb.UpdateProductResp{
		Product: s.toProductPB(result),
	}, nil
}

// DeleteProduct 删除商品
func (s *ProductServer) DeleteProduct(ctx context.Context, req *productpb.DeleteProductReq) (*productpb.DeleteProductResp, error) {
	// 从认证上下文获取用户ID
	userID := auth.UserIDFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "用户未认证")
	}

	err := s.productService.DeleteProduct(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "删除商品失败: %v", err)
	}

	return &productpb.DeleteProductResp{
		Success: true,
	}, nil
}

// GetProduct 获取商品详情
func (s *ProductServer) GetProduct(ctx context.Context, req *productpb.GetProductReq) (*productpb.GetProductResp, error) {
	result, err := s.productService.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取商品失败: %v", err)
	}

	return &productpb.GetProductResp{
		Product: s.toProductPB(result),
	}, nil
}

// ListProducts 获取商品列表
func (s *ProductServer) ListProducts(ctx context.Context, req *productpb.ListProductsReq) (*productpb.ListProductsResp, error) {
	var priceMin, priceMax *decimal.Decimal
	if req.PriceMin != "" {
		min, err := decimal.NewFromString(req.PriceMin)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "最低价格格式无效: %v", err)
		}
		priceMin = &min
	}
	if req.PriceMax != "" {
		max, err := decimal.NewFromString(req.PriceMax)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "最高价格格式无效: %v", err)
		}
		priceMax = &max
	}

	dto := &product.ListProductsDTO{
		Page:       int(req.Page),
		PageSize:   int(req.PageSize),
		CategoryID: &req.CategoryId,
		BrandID:    &req.BrandId,
		Status:     func() *productdomain.ProductStatus { s := productdomain.ProductStatus(req.Status); return &s }(),
		IsFeatured: &req.IsFeatured,
		Keyword:    req.Keyword,
		PriceMin:   priceMin,
		PriceMax:   priceMax,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
	}

	if req.CategoryId == "" {
		dto.CategoryID = nil
	}
	if req.BrandId == "" {
		dto.BrandID = nil
	}

	result, err := s.productService.ListProducts(ctx, dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取商品列表失败: %v", err)
	}

	products := make([]*productpb.Product, len(result.Products))
	for i, p := range result.Products {
		products[i] = s.toProductPB(p)
	}

	return &productpb.ListProductsResp{
		Products: products,
		Total:    result.Total,
		Page:     int32(result.Page),
		PageSize: int32(result.PageSize),
	}, nil
}

// SearchProducts 搜索商品
func (s *ProductServer) SearchProducts(ctx context.Context, req *productpb.SearchProductsReq) (*productpb.SearchProductsResp, error) {
	var priceMin, priceMax *decimal.Decimal
	if req.PriceMin != "" {
		min, err := decimal.NewFromString(req.PriceMin)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "最低价格格式无效: %v", err)
		}
		priceMin = &min
	}
	if req.PriceMax != "" {
		max, err := decimal.NewFromString(req.PriceMax)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "最高价格格式无效: %v", err)
		}
		priceMax = &max
	}

	dto := &product.SearchProductsDTO{
		Keyword:    req.Keyword,
		Page:       int(req.Page),
		PageSize:   int(req.PageSize),
		CategoryID: &req.CategoryId,
		BrandID:    &req.BrandId,
		PriceMin:   priceMin,
		PriceMax:   priceMax,
	}

	if req.CategoryId == "" {
		dto.CategoryID = nil
	}
	if req.BrandId == "" {
		dto.BrandID = nil
	}

	result, err := s.productService.SearchProducts(ctx, dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "搜索商品失败: %v", err)
	}

	products := make([]*productpb.Product, len(result.Products))
	for i, p := range result.Products {
		products[i] = s.toProductPB(p)
	}

	return &productpb.SearchProductsResp{
		Products: products,
		Total:    result.Total,
		Page:     int32(result.Page),
		PageSize: int32(result.PageSize),
	}, nil
}

// toProductPB 转换为protobuf商品对象
func (s *ProductServer) toProductPB(p *product.ProductDTO) *productpb.Product {
	pb := &productpb.Product{
		Id:               p.ID,
		Name:             p.Name,
		Slug:             p.Slug,
		Description:      p.Description,
		ShortDescription: p.ShortDescription,
		CategoryId:       p.CategoryID,
		MarketPrice:      p.MarketPrice,
		SalePrice:        p.SalePrice,
		CostPrice:        p.CostPrice,
		MainImageUrl:     p.MainImageURL,
		ImageUrls:        p.ImageURLs,
		VideoUrl:         p.VideoURL,
		Tags:             p.Tags,
		Status:           productpb.ProductStatus(p.Status),
		IsFeatured:       p.IsFeatured,
		IsVirtual:        p.IsVirtual,
		SeoTitle:         p.SEOTitle,
		SeoDescription:   p.SEODescription,
		SeoKeywords:      p.SEOKeywords,
		SortOrder:        int32(p.SortOrder),
		CreatedAt:        parseTime(p.CreatedAt),
		UpdatedAt:        parseTime(p.UpdatedAt),
		CategoryName:     p.CategoryName,
		BrandName:        p.BrandName,
	}

	if p.BrandID != nil {
		pb.BrandId = *p.BrandID
	}

	// PublishAt 字段在proto中不存在，暂时注释掉
	// if p.PublishAt != nil {
	//	pb.PublishAt = *p.PublishAt
	// }

	// 转换规格参数
	if p.Specifications != nil {
		pb.Specifications = make(map[string]string)
		for k, v := range p.Specifications {
			if str, ok := v.(string); ok {
				pb.Specifications[k] = str
			}
		}
	}

	// 转换SKUs
	if len(p.SKUs) > 0 {
		pb.Skus = make([]*productpb.ProductSKU, len(p.SKUs))
		for i, sku := range p.SKUs {
			pb.Skus[i] = s.toProductSKUPB(sku)
		}
	}

	return pb
}

// toProductSKUPB 转换为protobuf SKU对象
func (s *ProductServer) toProductSKUPB(sku *product.ProductSKUDTO) *productpb.ProductSKU {
	pb := &productpb.ProductSKU{
		Id:            sku.ID,
		ProductId:     sku.ProductID,
		SkuCode:       sku.SKUCode,
		Name:          sku.Name,
		MarketPrice:   sku.MarketPrice,
		SalePrice:     sku.SalePrice,
		CostPrice:     sku.CostPrice,
		StockQuantity: int32(sku.StockQuantity),
		// ReservedQuantity和SoldQuantity在proto中不存在
		// ReservedQuantity: int32(sku.ReservedQuantity),
		// SoldQuantity:     int32(sku.SoldQuantity),
		Weight:     int32(parseWeight(sku.Weight)),
		ImageUrl:   sku.ImageURL,
		Attributes: sku.Attributes,
		// Status字段在proto中不存在
		// Status:           int32(sku.Status),
		CreatedAt: parseTime(sku.CreatedAt),
		UpdatedAt: parseTime(sku.UpdatedAt),
	}

	// 转换尺寸 - 根据proto定义，Dimensions是string类型
	if sku.Dimensions != nil {
		// 将map转换为字符串，实际使用时可能需要JSON序列化
		pb.Dimensions = "dimensions_placeholder"
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

// parseWeight 解析重量字符串为整数
func parseWeight(weightStr string) int {
	if weightStr == "" {
		return 0
	}

	// 这里可以添加更复杂的重量解析逻辑
	// 暂时返回0，实际使用时需要根据业务需求实现
	return 0
}
