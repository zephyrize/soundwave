package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"soundwave-go/client"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Service struct {
		Name     string            `yaml:"name"`
		Port     int               `yaml:"port"`
		Version  string            `yaml:"version"`
		Metadata map[string]string `yaml:"metadata"`
	} `yaml:"service"`

	Registry struct {
		URL               string `yaml:"url"`
		HeartbeatInterval string `yaml:"heartbeat_interval"`
	} `yaml:"registry"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config.yaml", "配置文件路径")
	flag.Parse()

	// 加载配置文件
	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 解析心跳间隔
	heartbeatInterval, err := time.ParseDuration(config.Registry.HeartbeatInterval)
	if err != nil {
		log.Fatalf("解析心跳间隔失败: %v", err)
	}

	// 创建客户端配置
	clientConfig := &client.ClientConfig{
		ServiceName:       config.Service.Name,
		Port:              config.Service.Port,
		Version:           config.Service.Version,
		Metadata:          config.Service.Metadata,
		RegistryURL:       config.Registry.URL,
		HeartbeatInterval: heartbeatInterval,
		IP:                getLocalIP(), // 使用本机IP替代硬编码的IP
	}

	log.Printf("客户端配置: %+v", clientConfig) // 添加配置日志

	// 创建客户端实例
	c, err := client.NewClient(clientConfig)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 启动服务注册和心跳
	if err := c.Start(); err != nil {
		log.Fatalf("启动客户端失败: %v", err)
	}

	// 启动HTTP服务
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	go func() {
		addr := fmt.Sprintf(":%d", config.Service.Port)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Printf("HTTP服务启动失败: %v", err)
		}
	}()

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// 优雅关闭
	c.Stop()
	log.Println("服务已关闭")
}
