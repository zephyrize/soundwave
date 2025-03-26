package registry

import (
	"math/rand"
	"time"
)

// LoadBalancer 负载均衡接口
type LoadBalancer interface {
	Select([]*Service) *Service
}

// RandomBalancer 随机负载均衡器
type RandomBalancer struct {
	rand *rand.Rand
}

func NewRandomBalancer() *RandomBalancer {
	return &RandomBalancer{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (rb *RandomBalancer) Select(services []*Service) *Service {
	if len(services) == 0 {
		return nil
	}
	return services[rb.rand.Intn(len(services))]
}
