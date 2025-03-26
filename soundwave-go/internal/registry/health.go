package registry

import "time"

// HealthCheck 健康检查接口
type HealthCheck interface {
	Check(*Service) bool
}

// DefaultHealthCheck 默认健康检查实现
type DefaultHealthCheck struct {
	timeout time.Duration
}

func NewDefaultHealthCheck(timeout time.Duration) *DefaultHealthCheck {
	return &DefaultHealthCheck{
		timeout: timeout,
	}
}

func (hc *DefaultHealthCheck) Check(service *Service) bool {
	return service.Status == StatusUP && time.Since(service.LastHeartbeat) <= hc.timeout
}
