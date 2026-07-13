package grpcapi

import (
	"errors"
	"net/http"

	"github.com/ai-finops/ai-finops/packages/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toGRPCError(err error) error {
	var appErr *shared.AppError
	if errors.As(err, &appErr) {
		return status.Error(appErrToCode(appErr.Status), appErr.Message)
	}
	return status.Errorf(codes.Internal, "%v", err)
}

func appErrToCode(statusCode int) codes.Code {
	switch statusCode {
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.AlreadyExists
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	default:
		return codes.Internal
	}
}