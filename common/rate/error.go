package rate

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrLimitExceeded = status.Error(codes.ResourceExhausted, "limit exceeded")
)
