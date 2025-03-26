package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 服务配置
type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`

	Registry struct {
		HeartbeatInterval time.Duration `yaml:"heartbeat_interval"`
		ServiceTTL        time.Duration `yaml:"service_ttl"`
	} `yaml:"registry"`

	MongoDB struct {
		URI         string `yaml:"uri"`
		Database    string `yaml:"database"`
		Collections struct {
			Users  string `yaml:"users"`
			Menus  string `yaml:"menus"`
			Tokens string `yaml:"tokens"`
		} `yaml:"collections"`
	} `yaml:"mongodb"`

	JWT struct {
		Secret      string `yaml:"secret"`
		ExpireHours int    `yaml:"expire_hours"`
	} `yaml:"jwt"`
}

// ServerConfig HTTP服务器配置
type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

// RegistryConfig 注册中心配置
type RegistryConfig struct {
	HeartbeatInterval time.Duration `yaml:"heartbeatInterval"`
	HeartbeatTimeout  time.Duration `yaml:"heartbeatTimeout"`
}

// LoadConfig 从文件加载配置
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &cfg, nil
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
		}{
			Port: 7777,
			Host: "0.0.0.0",
		},
		Registry: struct {
			HeartbeatInterval time.Duration `yaml:"heartbeat_interval"`
			ServiceTTL        time.Duration `yaml:"service_ttl"`
		}{
			HeartbeatInterval: 10 * time.Second,
			ServiceTTL:        30 * time.Second,
		},
		MongoDB: struct {
			URI         string `yaml:"uri"`
			Database    string `yaml:"database"`
			Collections struct {
				Users  string `yaml:"users"`
				Menus  string `yaml:"menus"`
				Tokens string `yaml:"tokens"`
			} `yaml:"collections"`
		}{
			URI:      "mongodb://admin:password123@localhost:27017",
			Database: "soundwave",
			Collections: struct {
				Users  string `yaml:"users"`
				Menus  string `yaml:"menus"`
				Tokens string `yaml:"tokens"`
			}{
				Users:  "users",
				Menus:  "menus",
				Tokens: "tokens",
			},
		},
		JWT: struct {
			Secret      string `yaml:"secret"`
			ExpireHours int    `yaml:"expire_hours"`
		}{
			Secret:      "your-secret-key",
			ExpireHours: 24,
		},
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 验证服务器配置
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("无效的端口号: %d", c.Server.Port)
	}
	if c.Server.Host == "" {
		return fmt.Errorf("主机地址不能为空")
	}

	// 验证注册中心配置
	if c.Registry.HeartbeatInterval <= 0 {
		return fmt.Errorf("心跳间隔必须大于0")
	}
	if c.Registry.ServiceTTL <= c.Registry.HeartbeatInterval {
		return fmt.Errorf("服务过期时间必须大于心跳间隔")
	}

	return nil
}
