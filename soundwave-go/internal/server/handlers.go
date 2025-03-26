package server

import (
	"net/http"
	"soundwave-go/internal/registry"

	"github.com/gin-gonic/gin"
)

// RegisterService 处理服务注册请求
func (s *Server) RegisterService(c *gin.Context) {
	var req ServiceRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求参数: " + err.Error(),
		})
		return
	}

	service := &registry.Service{
		Name:     req.Name,
		ID:       req.ID,
		Hostname: req.Hostname,
		IP:       req.IP,
		Port:     req.Port,
		Metadata: req.Metadata,
		Version:  req.Version,
	}

	if err := s.registry.RegisterService(service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "服务注册失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "服务注册成功",
		"service": service,
	})
}

// DiscoverService 处理服务发现请求
func (s *Server) DiscoverService(c *gin.Context) {
	serviceName := c.Param("name")
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "服务名称不能为空",
		})
		return
	}

	services, err := s.registry.GetService(serviceName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"services": services,
	})
}

// ListServices 获取所有注册的服务
func (s *Server) ListServices(c *gin.Context) {
	services := s.registry.ListAllServices()
	c.JSON(http.StatusOK, gin.H{
		"services": services,
	})
}

// UpdateHeartbeat 处理服务心跳请求
func (s *Server) UpdateHeartbeat(c *gin.Context) {
	serviceName := c.Param("name")
	serviceID := c.Param("id")

	if serviceName == "" || serviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "服务名称和ID不能为空",
		})
		return
	}

	if err := s.registry.UpdateHeartbeat(serviceName, serviceID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "心跳更新成功",
	})
}

// GetServiceStats 获取服务统计信息
func (s *Server) GetServiceStats(c *gin.Context) {
	serviceName := c.Param("name")
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "服务名称不能为空",
		})
		return
	}

	stats, err := s.registry.GetServiceStats(serviceName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}
