package registry

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	// 服务健康检查的时间间隔
	healthCheckInterval = 10 * time.Second
	// 服务实例的过期时间
	serviceExpiration = 30 * time.Second
)

// 服务状态常量
const (
	// StatusUP 表示服务正常运行
	StatusUP ServiceStatus = "UP"
	// StatusDOWN 表示服务已离线
	StatusDOWN ServiceStatus = "DOWN"
	// StatusStarting 表示服务正在启动
	StatusStarting ServiceStatus = "STARTING"
	// StatusOutOfService 表示服务已手动下线
	StatusOutOfService ServiceStatus = "OUT_OF_SERVICE"
)

// ServiceStatus 定义服务状态类型
type ServiceStatus string

// Service 表示一个服务实例
type Service struct {
	Name          string            `json:"name"`
	ID            string            `json:"id"`
	Hostname      string            `json:"hostname"`
	IP            string            `json:"ip"` // 修改：使用IP替代address
	Port          int               `json:"port"`
	Metadata      map[string]string `json:"metadata"`
	Status        ServiceStatus     `json:"status"`
	LastHeartbeat time.Time         `json:"last_heartbeat"`
	Weight        int               `json:"weight"`
	StartTime     time.Time         `json:"start_time"`
	Version       string            `json:"version"`
}

// GetAddress 返回服务地址
func (s *Service) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}

// ValidateIP 验证IP地址格式
func (s *Service) ValidateIP() error {
	if s.IP == "" {
		return fmt.Errorf("IP地址不能为空")
	}

	// 验证IP地址格式
	ip := net.ParseIP(s.IP)
	if ip == nil {
		return fmt.Errorf("无效的IP地址格式: %s", s.IP)
	}

	return nil
}

// 新增：生成服务实例的唯一标识符
func (s *Service) UniqueID() string {
	return fmt.Sprintf("%s-%s-%s", s.Name, s.Hostname, s.ID)
}

// ServiceRegistry 服务注册中心
type ServiceRegistry struct {
	services    map[string]*Service // 修改为使用 uniqueID -> Service
	serviceMap  map[string][]string // 新增：服务名称 -> uniqueID列表的映射
	mutex       sync.RWMutex
	healthCheck HealthCheck
	balancer    LoadBalancer
}

// RegistryOption 定义注册中心的配置选项
type RegistryOption func(*ServiceRegistry)

// WithHealthCheck 设置健康检查器
func WithHealthCheck(hc HealthCheck) RegistryOption {
	return func(sr *ServiceRegistry) {
		sr.healthCheck = hc
	}
}

// WithLoadBalancer 设置负载均衡器
func WithLoadBalancer(lb LoadBalancer) RegistryOption {
	return func(sr *ServiceRegistry) {
		sr.balancer = lb
	}
}
