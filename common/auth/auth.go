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

func UserIDFromContext(ctx context.Context) string {
	values := metadata.ValueFromIncomingContext(ctx, GrpcUserIDMetadataKey)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
