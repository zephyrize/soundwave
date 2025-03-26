package registry

import (
	"context"
	"fmt"
	"soundwave-go/internal/logger"
	"time"
)

// NewServiceRegistry 创建新的服务注册中心
func NewServiceRegistry(opts ...RegistryOption) *ServiceRegistry {
	sr := &ServiceRegistry{
		services:    make(map[string]*Service),
		serviceMap:  make(map[string][]string),
		healthCheck: NewDefaultHealthCheck(serviceExpiration),
		balancer:    NewRandomBalancer(),
	}

	// 应用选项
	for _, opt := range opts {
		opt(sr)
	}

	return sr
}

// RegisterService 注册服务
func (sr *ServiceRegistry) RegisterService(service *Service) error {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()

	// 验证服务信息
	if service.Name == "" || service.ID == "" || service.Hostname == "" {
		logger.ErrorLogger.Printf("服务注册失败：信息不完整 %+v", service)
		return fmt.Errorf("服务名称、ID和主机名不能为空")
	}

	// 验证IP地址
	if err := service.ValidateIP(); err != nil {
		logger.ErrorLogger.Printf("服务注册失败：%v", err)
		return err
	}

	// 验证端口
	if service.Port <= 0 || service.Port > 65535 {
		return fmt.Errorf("无效的端口号: %d", service.Port)
	}

	// 设置服务状态和心跳时间
	service.Status = StatusUP
	service.LastHeartbeat = time.Now()
	if service.StartTime.IsZero() {
		service.StartTime = time.Now()
	}

	// 生成唯一标识符
	uniqueID := service.UniqueID()

	// 存储服务实例
	sr.services[uniqueID] = service

	// 更新服务名称到uniqueID的映射
	if _, exists := sr.serviceMap[service.Name]; !exists {
		sr.serviceMap[service.Name] = make([]string, 0)
	}

	// 检查是否已存在该uniqueID
	found := false
	for _, id := range sr.serviceMap[service.Name] {
		if id == uniqueID {
			found = true
			break
		}
	}
	if !found {
		sr.serviceMap[service.Name] = append(sr.serviceMap[service.Name], uniqueID)
	}

	logger.InfoLogger.Printf("注册服务实例：%s，地址：%s", uniqueID, service.GetAddress())
	return nil
}

// DeregisterService 注销服务
func (sr *ServiceRegistry) DeregisterService(name, id string) error {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()

	// 检查服务是否存在
	uniqueIDs, exists := sr.serviceMap[name]
	if !exists {
		return fmt.Errorf("服务 %s 不存在", name)
	}

	// 查找并移除指定的服务实例
	for i, uniqueID := range uniqueIDs {
		if service, ok := sr.services[uniqueID]; ok {
			if service.ID == id {
				// 从serviceMap中移除该uniqueID
				sr.serviceMap[name] = append(uniqueIDs[:i], uniqueIDs[i+1:]...)
				// 从services中删除该服务实例
				delete(sr.services, uniqueID)

				// 如果该服务没有实例了，则删除该服务条目
				if len(sr.serviceMap[name]) == 0 {
					delete(sr.serviceMap, name)
				}
				return nil
			}
		}
	}

	return fmt.Errorf("服务实例 %s 不存在", id)
}

// GetService 获取服务实例
func (sr *ServiceRegistry) GetService(name string) ([]*Service, error) {
	sr.mutex.RLock()
	defer sr.mutex.RUnlock()

	uniqueIDs, exists := sr.serviceMap[name]
	if !exists {
		return nil, fmt.Errorf("服务 %s 不存在", name)
	}

	// 获取所有活跃的服务实例
	activeServices := make([]*Service, 0)
	for _, uniqueID := range uniqueIDs {
		if service, ok := sr.services[uniqueID]; ok {
			if service.Status == StatusUP && time.Since(service.LastHeartbeat) <= serviceExpiration {
				activeServices = append(activeServices, service)
			}
		}
	}

	if len(activeServices) == 0 {
		return nil, fmt.Errorf("服务 %s 没有可用的实例", name)
	}

	return activeServices, nil
}

// UpdateHeartbeat 更新服务心跳时间
func (sr *ServiceRegistry) UpdateHeartbeat(name, id string) error {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()

	// 检查服务是否存在
	uniqueIDs, exists := sr.serviceMap[name]
	if !exists {
		return fmt.Errorf("服务 %s 不存在", name)
	}

	// 查找并更新服务实例
	for _, uniqueID := range uniqueIDs {
		if service, ok := sr.services[uniqueID]; ok {
			if service.ID == id {
				service.LastHeartbeat = time.Now()
				service.Status = StatusUP
				return nil
			}
		}
	}

	return fmt.Errorf("服务实例 %s 不存在", id)
}

// ListAllServices 获取所有注册的服务
func (sr *ServiceRegistry) ListAllServices() map[string][]*Service {
	sr.mutex.RLock()
	defer sr.mutex.RUnlock()

	// 创建服务列表的副本
	result := make(map[string][]*Service)

	// 遍历serviceMap获取所有服务
	for name, uniqueIDs := range sr.serviceMap {
		services := make([]*Service, 0)
		for _, uniqueID := range uniqueIDs {
			if service, ok := sr.services[uniqueID]; ok {
				services = append(services, service)
			}
		}
		if len(services) > 0 {
			result[name] = services
		}
	}

	return result
}

// StartHealthCheck 启动健康检查定时任务
func (sr *ServiceRegistry) StartHealthCheck(ctx context.Context, interval time.Duration) {
	logger.InfoLogger.Printf("启动健康检查，间隔时间：%v", interval)
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				logger.InfoLogger.Println("健康检查已停止")
				return
			case <-ticker.C:
				sr.checkServicesHealth()
			}
		}
	}()
}

// checkServicesHealth 检查所有服务的健康状态
func (sr *ServiceRegistry) checkServicesHealth() {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()

	now := time.Now()
	// 遍历所有服务
	for serviceName, uniqueIDs := range sr.serviceMap {
		activeUniqueIDs := make([]string, 0)

		// 检查每个服务实例
		for _, uniqueID := range uniqueIDs {
			if service, ok := sr.services[uniqueID]; ok {
				// 如果服务实例在过期时间内有心跳，则保留
				if service.Status == StatusUP && now.Sub(service.LastHeartbeat) <= serviceExpiration {
					activeUniqueIDs = append(activeUniqueIDs, uniqueID)
				} else {
					// 更新服务状态为离线
					service.Status = StatusDOWN
					// 从services中移除过期的服务实例
					delete(sr.services, uniqueID)
				}
			}
		}

		// 更新serviceMap
		if len(activeUniqueIDs) > 0 {
			sr.serviceMap[serviceName] = activeUniqueIDs
		} else {
			// 如果没有活跃实例，删除该服务
			delete(sr.serviceMap, serviceName)
		}
	}
}

// GetServiceWithLoadBalancing 使用负载均衡获取服务实例
func (sr *ServiceRegistry) GetServiceWithLoadBalancing(name string) (*Service, error) {
	services, err := sr.GetService(name)
	if err != nil {
		return nil, err
	}

	service := sr.balancer.Select(services)
	if service == nil {
		return nil, fmt.Errorf("没有可用的服务实例")
	}

	return service, nil
}
