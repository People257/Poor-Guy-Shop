package server

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var InternalStatus = status.Error(codes.Internal, "Internal Server Error")
