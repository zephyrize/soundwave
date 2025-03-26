package client

import "time"

// ClientConfig 客户端配置
type ClientConfig struct {
	ServiceName       string            // 服务名称
	ServiceID         string            // 服务ID
	IP                string            // 服务IP
	Port              int               // 服务端口
	Version           string            // 服务版本
	Metadata          map[string]string // 服务元数据
	RegistryURL       string            // 服务中心地址
	HeartbeatInterval time.Duration     // 心跳间隔
}

// DefaultConfig 返回默认配置
func DefaultConfig() *ClientConfig {
	return &ClientConfig{
		IP:                "127.0.0.1",
		Port:              8080,
		Version:           "1.0.0",
		RegistryURL:       "http://localhost:7777",
		HeartbeatInterval: 10 * time.Second,
		Metadata:          make(map[string]string),
	}
}
