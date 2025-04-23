package ecode

import "strconv"

// 自定义错误码范围: 10000-50000
type ErrorCode int

const (
	// 系统级错误: 10000-19999

	// 业务逻辑错误: 40000-49999
	ErrAuthFailed        ErrorCode = 40001
	ErrInvalidParam      ErrorCode = 40002
	ErrNotFound          ErrorCode = 40003
	ErrAlreadyExists     ErrorCode = 40004
	ErrPermissionDenied  ErrorCode = 40005
	ErrOperationFailed   ErrorCode = 40006
	ErrRateLimit         ErrorCode = 40007
	ErrResourceExhausted ErrorCode = 40008
	ErrInvalidOperation  ErrorCode = 40009
	ErrDataConflict      ErrorCode = 40010
)

var ErrorCode_name = map[int32]string{
	40001: "AUTH_FAILED",
	40002: "INVALID_PARAM",
	40003: "NOT_FOUND",
	40004: "ALREADY_EXISTS",
	40005: "PERMISSION_DENIED",
	40006: "OPERATION_FAILED",
	40007: "RATE_LIMIT",
	40008: "RESOURCE_EXHAUSTED",
	40009: "INVALID_OPERATION",
	40010: "DATA_CONFLICT",
}

func (x ErrorCode) String() string {
	s, ok := ErrorCode_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
