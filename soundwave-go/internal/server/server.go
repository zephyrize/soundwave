package server

import (
	"context"
	"fmt"
	"net/http"
	"soundwave-go/internal/config"
	"soundwave-go/internal/db"
	"soundwave-go/internal/logger"
	"soundwave-go/internal/middleware"
	"soundwave-go/internal/models"
	"soundwave-go/internal/registry"
	"soundwave-go/internal/service"

	"github.com/gin-gonic/gin"
)

// Server HTTP服务器结构
type Server struct {
	engine      *gin.Engine
	registry    *registry.ServiceRegistry
	config      *config.Config
	ctx         context.Context
	cancel      context.CancelFunc
	authService *service.AuthService
	menuService *service.MenuService
	userService *service.UserService
}

func NewServer(cfg *config.Config) *Server {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	r := gin.Default()
	registry := registry.NewServiceRegistry()
	ctx, cancel := context.WithCancel(context.Background())

	// 添加CORS中间件
	r.Use(middleware.CORS())

	// 初始化 MongoDB 连接
	mongodb, err := db.NewMongoDB(cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		logger.ErrorLogger.Fatalf("连接MongoDB失败: %v", err)
	}

	// 加载初始化数据
	initData, err := config.LoadInitData("configs/init_data.yaml")
	if err != nil {
		logger.ErrorLogger.Fatalf("加载初始化数据失败: %v", err)
	}

	// 初始化数据库
	if err := db.InitializeData(mongodb.Database(), initData); err != nil {
		logger.ErrorLogger.Fatalf("初始化数据失败: %v", err)
	}

	server := &Server{
		engine:      r,
		registry:    registry,
		config:      cfg,
		ctx:         ctx,
		cancel:      cancel,
		authService: service.NewAuthService(mongodb.Database(), cfg),
		menuService: service.NewMenuService(mongodb.Collection("menus")),
		userService: service.NewUserService(mongodb.Database()),
	}

	// 注册路由
	server.registerRoutes()
	// 启动健康检查
	registry.StartHealthCheck(ctx, cfg.Registry.HeartbeatInterval)

	logger.InfoLogger.Printf("服务器初始化完成，配置：%+v", cfg)
	return server
}

func (s *Server) registerRoutes() {
	// 服务注册接口
	s.engine.POST("/services", s.RegisterService)
	// 服务发现接口
	s.engine.GET("/services/:name", s.DiscoverService)
	// 获取所有服务列表
	s.engine.GET("/services", s.ListServices)
	// 服务心跳接口
	s.engine.PUT("/services/:name/:id/heartbeat", s.UpdateHeartbeat)
	// 服务统计信息
	s.engine.GET("/services/:name/stats", s.GetServiceStats)
	// 负载均衡获取服务
	s.engine.GET("/services/:name/instance", s.GetServiceInstance)

	// 认证相关路由
	auth := s.engine.Group("/auth")
	{
		auth.POST("/register", s.HandleRegister)
		auth.POST("/login", s.HandleLogin)
	}

	// 需要认证的路由
	api := s.engine.Group("/api")
	{
		menus := api.Group("/menus")
		menus.Use(middleware.AuthRequired(s.config, models.PermissionViewServices))
		{
			menus.GET("", s.GetUserMenus)
		}

		services := api.Group("/services")
		services.Use(middleware.AuthRequired(s.config, models.PermissionViewServices))
		{
			services.GET("", s.ListServices)
		}

		stats := api.Group("/stats")
		stats.Use(middleware.AuthRequired(s.config, models.PermissionViewStats))
		{
			stats.GET("", s.GetServiceStats)
		}

		// 用户相关路由
		user := api.Group("/user")
		user.Use(middleware.AuthRequired(s.config, models.PermissionViewServices))
		{
			user.POST("/change-password", s.HandleChangePassword)
		}

		// 用户管理路由
		users := api.Group("/users")
		users.Use(middleware.AuthRequired(s.config, models.PermissionManageUsers))
		{
			users.GET("", s.ListUsers)
			users.POST("", s.CreateUser)
			users.GET("/:id", s.GetUser)
			users.PUT("/:id", s.UpdateUser)
			users.DELETE("/:id", s.DeleteUser)
			users.POST("/:id/reset-password", s.ResetUserPassword)
		}
	}

	logger.InfoLogger.Println("路由注册完成")
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	logger.InfoLogger.Printf("服务器启动，监听地址：%s", addr)
	return s.engine.Run(addr)
}

// Shutdown 优雅关闭服务器
func (s *Server) Shutdown() {
	logger.InfoLogger.Println("服务器开始关闭...")
	if s.cancel != nil {
		s.cancel()
	}
	logger.InfoLogger.Println("服务器已关闭")
}

// GetServiceInstance 获取负载均衡后的服务实例
func (s *Server) GetServiceInstance(c *gin.Context) {
	serviceName := c.Param("name")
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "服务名称不能为空",
		})
		return
	}

	service, err := s.registry.GetServiceWithLoadBalancing(serviceName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"service": service,
	})
}
