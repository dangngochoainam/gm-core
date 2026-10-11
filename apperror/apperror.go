// Package apperror holds the business error sentinels every service shares. Each
// one pairs a gRPC code with a reason named in commondtov1.SystemCode, so a
// reason the gateway cannot turn into a systemCode does not compile.
//
// This is the one hand-written Go package in gm-core, on purpose: the reasons
// belong next to the enum that numbers them, and gopkg stays neutral.
package apperror

import (
	commondtov1 "github.com/dangngochoainam/gm-core/application/common/v1/dto"
	"github.com/dangngochoainam/gopkg/responsehelper"
	"google.golang.org/grpc/codes"
)

var (
	ErrBadRequest          = newError(codes.InvalidArgument, commondtov1.SystemCode_BAD_REQUEST, "bad request")
	ErrUnauthorized        = newError(codes.Unauthenticated, commondtov1.SystemCode_UNAUTHORIZED, "unauthorized")
	ErrForbidden           = newError(codes.PermissionDenied, commondtov1.SystemCode_FORBIDDEN, "forbidden")
	ErrNotFound            = newError(codes.NotFound, commondtov1.SystemCode_NOT_FOUND, "not found")
	ErrUserNotFound        = newError(codes.NotFound, commondtov1.SystemCode_USER_NOT_FOUND, "user not found")
	ErrSessionNotFound     = newError(codes.NotFound, commondtov1.SystemCode_SESSION_NOT_FOUND, "session not found")
	ErrConflict            = newError(codes.AlreadyExists, commondtov1.SystemCode_CONFLICT, "conflict")
	ErrEmailExisted        = newError(codes.AlreadyExists, commondtov1.SystemCode_EMAIL_EXISTED, "email existed")
	ErrAuthAccountExisted  = newError(codes.AlreadyExists, commondtov1.SystemCode_AUTH_ACCOUNT_EXISTED, "auth account existed")
	ErrTooManyRequests     = newError(codes.ResourceExhausted, commondtov1.SystemCode_TOO_MANY_REQUESTS, "too many requests")
	ErrMethodNotAllowed    = newError(codes.Unimplemented, commondtov1.SystemCode_METHOD_NOT_ALLOWED, "method not allowed")
	ErrInternalServerError = newError(codes.Internal, commondtov1.SystemCode_INTERNAL_SERVER_ERROR, "internal server error")
)

func newError(code codes.Code, reason commondtov1.SystemCode, message string) *responsehelper.Error {
	return responsehelper.New(code, reason.String(), message)
}
