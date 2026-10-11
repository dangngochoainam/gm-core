package apperror

import (
	"testing"

	commondtov1 "github.com/dangngochoainam/gm-core/application/common/v1/dto"
	"github.com/dangngochoainam/gopkg/grpcmiddlewarehelper"
	"github.com/dangngochoainam/gopkg/responsehelper"
	"google.golang.org/grpc/codes"
)

func TestSentinels(t *testing.T) {
	tests := []struct {
		name       string
		err        *responsehelper.Error
		code       codes.Code
		reason     string
		message    string
		systemCode int32
	}{
		{"ErrBadRequest", ErrBadRequest, codes.InvalidArgument, "BAD_REQUEST", "bad request", 40000},
		{"ErrUnauthorized", ErrUnauthorized, codes.Unauthenticated, "UNAUTHORIZED", "unauthorized", 40100},
		{"ErrForbidden", ErrForbidden, codes.PermissionDenied, "FORBIDDEN", "forbidden", 40300},
		{"ErrNotFound", ErrNotFound, codes.NotFound, "NOT_FOUND", "not found", 40400},
		{"ErrUserNotFound", ErrUserNotFound, codes.NotFound, "USER_NOT_FOUND", "user not found", 40401},
		{"ErrSessionNotFound", ErrSessionNotFound, codes.NotFound, "SESSION_NOT_FOUND", "session not found", 40402},
		{"ErrConflict", ErrConflict, codes.AlreadyExists, "CONFLICT", "conflict", 40900},
		{"ErrEmailExisted", ErrEmailExisted, codes.AlreadyExists, "EMAIL_EXISTED", "email existed", 40901},
		{"ErrAuthAccountExisted", ErrAuthAccountExisted, codes.AlreadyExists, "AUTH_ACCOUNT_EXISTED", "auth account existed", 40902},
		{"ErrTooManyRequests", ErrTooManyRequests, codes.ResourceExhausted, "TOO_MANY_REQUESTS", "too many requests", 42900},
		{"ErrMethodNotAllowed", ErrMethodNotAllowed, codes.Unimplemented, "METHOD_NOT_ALLOWED", "method not allowed", 50100},
		{"ErrInternalServerError", ErrInternalServerError, codes.Internal, "INTERNAL_SERVER_ERROR", "internal server error", 50000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Fatal("sentinel is nil")
			}
			if tt.err.Code != tt.code || tt.err.Reason != tt.reason || tt.err.Error() != tt.message {
				t.Errorf("got {%v, %q, %q}, want {%v, %q, %q}", tt.err.Code, tt.err.Reason, tt.err.Error(), tt.code, tt.reason, tt.message)
			}
			if got := commondtov1.SystemCode_value[tt.err.Reason]; got != tt.systemCode {
				t.Errorf("systemCode = %d, want %d", got, tt.systemCode)
			}
		})
	}
}

// TestGopkgReasonsInEnum keeps the reasons gopkg attaches on its own in step with
// the enum, which gopkg cannot import.
func TestGopkgReasonsInEnum(t *testing.T) {
	for reason, want := range map[string]int32{
		grpcmiddlewarehelper.ReasonValidationFailed: 40001,
		grpcmiddlewarehelper.ReasonDeadlineExceeded: 50400,
		grpcmiddlewarehelper.ReasonCanceled:         49900,
		grpcmiddlewarehelper.ReasonInternal:         50000,
	} {
		got, ok := commondtov1.SystemCode_value[reason]
		if !ok || got != want {
			t.Errorf("SystemCode_value[%q] = %d (present %v), want %d", reason, got, ok, want)
		}
	}
}
