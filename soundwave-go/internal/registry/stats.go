package registry

import (
	"fmt"
	"time"
)

// ServiceStats 服务统计信息
type ServiceStats struct {
	TotalInstances     int           `json:"total_instances"`
	HealthyInstances   int           `json:"healthy_instances"`
	UnhealthyInstances int           `json:"unhealthy_instances"`
	AverageUptime      time.Duration `json:"average_uptime"`
	LastUpdateTime     time.Time     `json:"last_update_time"`
}

// GetServiceStats 获取服务统计信息
func (sr *ServiceRegistry) GetServiceStats(name string) (*ServiceStats, error) {
	sr.mutex.RLock()
	defer sr.mutex.RUnlock()

	uniqueIDs, exists := sr.serviceMap[name]
	if !exists {
		return nil, fmt.Errorf("服务 %s 不存在", name)
	}

	stats := &ServiceStats{
		TotalInstances: len(uniqueIDs),
		LastUpdateTime: time.Now(),
	}

	var totalUptime time.Duration
	for _, uniqueID := range uniqueIDs {
		if service, ok := sr.services[uniqueID]; ok {
			if sr.healthCheck.Check(service) {
				stats.HealthyInstances++
			} else {
				stats.UnhealthyInstances++
			}
			totalUptime += time.Since(service.StartTime)
		}
	}

	if stats.TotalInstances > 0 {
		stats.AverageUptime = totalUptime / time.Duration(stats.TotalInstances)
	}

	return stats, nil
}
