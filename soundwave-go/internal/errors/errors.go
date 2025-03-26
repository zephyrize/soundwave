package errors

import "fmt"

type ServiceError struct {
	Code    int
	Message string
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("错误码: %d, 错误信息: %s", e.Code, e.Message)
}

var (
	ErrServiceNotFound    = &ServiceError{Code: 404, Message: "服务不存在"}
	ErrServiceUnavailable = &ServiceError{Code: 503, Message: "服务不可用"}
	ErrInvalidRequest     = &ServiceError{Code: 400, Message: "无效的请求参数"}
)
