package interceptor

import (
	"context"
	"errors"
	"fmt"

	"buf.build/go/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// ValidateUnaryServerInterceptor returns a new unary server interceptor that validates incoming messages.
// If the request is invalid, clients may access a structured representation of the validation failure as an error detail.
func ValidateUnaryServerInterceptor(validator protovalidate.Validator) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		if err := validateMsg(req, validator); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func validateMsg(m interface{}, validator protovalidate.Validator) error {
	msg, ok := m.(proto.Message)
	if !ok {
		return status.Errorf(codes.Internal, "unsupported message type: %T", m)
	}
	err := validator.Validate(msg)
	if err == nil {
		return nil
	}
	var valErr *protovalidate.ValidationError
	if errors.As(err, &valErr) {
		// Message is invalid.
		errMsg := buildErrMsg(valErr)

		st := status.New(codes.InvalidArgument, errMsg)
		return st.Err()
	}
	// CEL expression doesn't compile or type-check.
	return status.Error(codes.Internal, err.Error())
}

func buildErrMsg(valErr *protovalidate.ValidationError) string {
	errMsg := ""
	violations := valErr.Violations

	if len(violations) > 0 && violations[0].Proto != nil && violations[0].Proto.Message != nil {

		errMsg = *violations[0].Proto.Message

		if errMsg == "" && violations[0].Proto.GetField() != nil {
			elements := violations[0].Proto.GetField().GetElements()
			if len(elements) == 0 {
				return errMsg
			}
			errMsg = fmt.Sprintf("field %s validation failed", elements[0].GetFieldName())
		}
	}
	return errMsg
}
