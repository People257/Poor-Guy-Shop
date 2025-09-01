package file

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	filepb "github.com/people257/poor-guy-shop/oss-infra/gen/proto/oss/file"
	"github.com/people257/poor-guy-shop/oss-infra/internal/application/file"
)

// Handler OSS文件服务gRPC处理器
type Handler struct {
	filepb.UnimplementedFileServiceServer
	fileApp *file.Service
}

// NewHandler 创建文件服务处理器
func NewHandler(fileApp *file.Service) *Handler {
	return &Handler{
		fileApp: fileApp,
	}
}

// UploadFile 上传文件
func (h *Handler) UploadFile(ctx context.Context, req *filepb.UploadFileReq) (*filepb.UploadFileResp, error) {
	result, err := h.fileApp.UploadFile(ctx, &file.UploadFileReq{
		FileData:   req.FileData,
		Filename:   req.Filename,
		Category:   req.Category,
		Visibility: req.Visibility,
	})
	if err != nil {
		return nil, err
	}

	return &filepb.UploadFileResp{
		FileInfo: &filepb.FileInfo{
			FileId:      result.FileInfo.FileID,
			Filename:    result.FileInfo.Filename,
			FileKey:     result.FileInfo.FileKey,
			FileSize:    result.FileInfo.FileSize,
			MimeType:    result.FileInfo.MimeType,
			Category:    result.FileInfo.Category,
			OwnerId:     result.FileInfo.OwnerID,
			Visibility:  result.FileInfo.Visibility,
			CreatedAt:   timestamppb.New(result.FileInfo.CreatedAt),
			DownloadUrl: result.FileInfo.DownloadURL,
		},
	}, nil
}

// GetDownloadUrl 获取文件下载URL
func (h *Handler) GetDownloadUrl(ctx context.Context, req *filepb.GetDownloadUrlReq) (*filepb.GetDownloadUrlResp, error) {
	result, err := h.fileApp.GetDownloadURL(ctx, &file.GetDownloadURLReq{
		FileID:    req.FileId,
		ExpiresIn: req.ExpiresIn,
	})
	if err != nil {
		return nil, err
	}

	return &filepb.GetDownloadUrlResp{
		DownloadUrl: result.DownloadURL,
		ExpiresIn:   result.ExpiresIn,
	}, nil
}
