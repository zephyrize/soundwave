package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"soundwave-go/internal/logger"
	"time"
)

// Client 服务注册客户端
type Client struct {
	config *ClientConfig
	ctx    context.Context
	cancel context.CancelFunc
}

// NewClient 创建新的客户端实例
func NewClient(config *ClientConfig) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// 获取主机名
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("获取主机名失败: %v", err)
	}

	// 如果没有指定ServiceID，使用hostname作为后缀
	if config.ServiceID == "" {
		config.ServiceID = fmt.Sprintf("%s-%s", config.ServiceName, hostname)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		config: config,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// Start 启动服务注册和心跳
func (c *Client) Start() error {
	// 注册服务
	if err := c.register(); err != nil {
		return fmt.Errorf("服务注册失败: %v", err)
	}

	// 启动心跳
	go c.heartbeat()

	return nil
}

// Stop 停止服务
func (c *Client) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
}

func (c *Client) register() error {
	// 验证必要字段
	if c.config.ServiceName == "" {
		return fmt.Errorf("服务名称不能为空")
	}
	if c.config.Port <= 0 || c.config.Port > 65535 {
		return fmt.Errorf("无效的端口号: %d", c.config.Port)
	}

	hostname := getHostname()
	data := map[string]interface{}{
		"name":     c.config.ServiceName,
		"id":       c.config.ServiceID,
		"hostname": hostname,
		"ip":       c.config.IP,
		"port":     c.config.Port,
		"version":  c.config.Version,
		"metadata": c.config.Metadata,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("JSON编码失败: %v", err)
	}

	log.Printf("发送注册请求: %s", string(jsonData))

	resp, err := http.Post(
		fmt.Sprintf("%s/services", c.config.RegistryURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("服务注册失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	logger.InfoLogger.Printf("服务 %s 注册成功", c.config.ServiceName)
	return nil
}

func (c *Client) heartbeat() {
	ticker := time.NewTicker(c.config.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			if err := c.sendHeartbeat(); err != nil {
				logger.ErrorLogger.Printf("发送心跳失败: %v", err)
			}
		}
	}
}

func (c *Client) sendHeartbeat() error {
	url := fmt.Sprintf("%s/services/%s/%s/heartbeat",
		c.config.RegistryURL,
		c.config.ServiceName,
		c.config.ServiceID,
	)

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("心跳请求失败，状态码: %d", resp.StatusCode)
	}

	return nil
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown-host"
	}
	return hostname
}
