package main

import (
	"flag"
	"os"
	"os/signal"
	"soundwave-go/internal/config"
	"soundwave-go/internal/logger"
	"soundwave-go/internal/server"
	"syscall"
)

var (
	configPath = flag.String("config", "configs/config.yaml", "配置文件路径")
)

func main() {
	flag.Parse()

	// 加载配置
	var cfg *config.Config

	if _, err := os.Stat(*configPath); err == nil {
		// 配置文件存在，尝试加载
		if loadedCfg, err := config.LoadConfig(*configPath); err != nil {
			logger.ErrorLogger.Printf("加载配置文件失败: %v, 将使用默认配置", err)
			cfg = config.DefaultConfig()
		} else {
			logger.InfoLogger.Printf("成功加载配置文件: %s", *configPath)
			cfg = loadedCfg
		}
	} else {
		// 配置文件不存在，使用默认配置
		logger.InfoLogger.Printf("配置文件不存在: %s, 将使用默认配置", *configPath)
		cfg = config.DefaultConfig()
	}

	logger.InfoLogger.Printf("当前配置: %+v", cfg)

	// 初始化服务器
	app := server.NewServer(cfg)

	// 处理优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// 启动服务器
		if err := app.Run(); err != nil {
			logger.ErrorLogger.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号
	<-quit
	logger.InfoLogger.Println("收到关闭信号...")

	// 执行清理操作
	app.Shutdown()
}
