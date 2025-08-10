package auth

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrUnAuthStatus = status.Error(codes.Unauthenticated, "unauthenticate")
