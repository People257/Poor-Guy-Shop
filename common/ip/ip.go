package ip

import (
	"context"
	"google.golang.org/grpc/metadata"
	"strings"
)

// GetIPFromMetadata 从 Metadata 中获取 IP 地址 (XForwardedFor)
func GetIPFromMetadata(ctx context.Context) string {
	addrs := metadata.ValueFromIncomingContext(ctx, "x-forwarded-for")
	if len(addrs) > 0 {
		return strings.Split(addrs[0], ",")[0]
	}
	return ""
}
