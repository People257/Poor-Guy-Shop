package auth

import (
	"context"
	"google.golang.org/grpc/metadata"
)

const (
	GrpcUserIDMetadataKey = "user-id"
	GrpcTokenMetadataKey  = "authorization"
	HttpUserIDHeaderKey   = "Grpc-Metadata-User-Id"
)

type userIDCtxKeyType struct{}

var userIDCtxKey userIDCtxKeyType = struct{}{}

func UserIDFromContext(ctx context.Context) string {
	values := metadata.ValueFromIncomingContext(ctx, GrpcUserIDMetadataKey)
	if len(values) == 1 {
		return values[0]
	}

	return ""
}

func GatewayUserIDFromContext(ctx context.Context) string {
	contextUserID, ok := ctx.Value(userIDCtxKey).(string)
	if ok {
		return contextUserID
	}
	return ""
}

// WithUserIDGRPCOutgoingContext 将 userID 写入 gRPC metadata 并返回新的 context
func WithUserIDGRPCOutgoingContext(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, GrpcUserIDMetadataKey, GatewayUserIDFromContext(ctx))
}
