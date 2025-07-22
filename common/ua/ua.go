package ua

import (
	"context"
	"google.golang.org/grpc/metadata"
)

func GetUserAgentFromMetadata(ctx context.Context) string {
	values := metadata.ValueFromIncomingContext(ctx, "grpcgateway-user-agent")
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
