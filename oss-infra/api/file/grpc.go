package file

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
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

// GetFileInfo 获取文件信息
func (h *Handler) GetFileInfo(ctx context.Context, req *filepb.GetFileInfoReq) (*filepb.GetFileInfoResp, error) {
	result, err := h.fileApp.GetFileInfo(ctx, &file.GetFileInfoReq{
		FileID: req.FileId,
	})
	if err != nil {
		return nil, err
	}

	return &filepb.GetFileInfoResp{
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

// ListFiles 获取文件列表
func (h *Handler) ListFiles(ctx context.Context, req *filepb.ListFilesReq) (*filepb.ListFilesResp, error) {
	result, err := h.fileApp.ListFiles(ctx, &file.ListFilesReq{
		Category: req.Category,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	files := make([]*filepb.FileInfo, 0, len(result.Files))
	for _, f := range result.Files {
		files = append(files, &filepb.FileInfo{
			FileId:      f.FileID,
			Filename:    f.Filename,
			FileKey:     f.FileKey,
			FileSize:    f.FileSize,
			MimeType:    f.MimeType,
			Category:    f.Category,
			OwnerId:     f.OwnerID,
			Visibility:  f.Visibility,
			CreatedAt:   timestamppb.New(f.CreatedAt),
			DownloadUrl: f.DownloadURL,
		})
	}

	return &filepb.ListFilesResp{
		Files:    files,
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
	}, nil
}

// DeleteFile 删除文件
func (h *Handler) DeleteFile(ctx context.Context, req *filepb.DeleteFileReq) (*emptypb.Empty, error) {
	err := h.fileApp.DeleteFile(ctx, &file.DeleteFileReq{
		FileID: req.FileId,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// CheckFileAccess 检查文件访问权限 (内部RPC)
func (h *Handler) CheckFileAccess(ctx context.Context, req *filepb.CheckFileAccessReq) (*filepb.CheckFileAccessResp, error) {
	result, err := h.fileApp.CheckFileAccess(ctx, &file.CheckFileAccessReq{
		FileID: req.FileId,
		UserID: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	return &filepb.CheckFileAccessResp{
		HasAccess: result.HasAccess,
		Reason:    result.Reason,
	}, nil
}
