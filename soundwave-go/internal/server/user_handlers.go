package server

import (
	"net/http"
	"soundwave-go/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// ListUsers 获取用户列表
func (s *Server) ListUsers(c *gin.Context) {
	users, err := s.userService.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 清除敏感信息
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"users":   users,
		"message": "获取成功",
	})
}

// CreateUser 创建用户
func (s *Server) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的请求数据",
		})
		return
	}

	// 验证用户输入
	if err := s.userService.ValidateUserInput(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// 检查是否已存在同名用户
	existingUser, _ := s.userService.GetUserByUsername(c.Request.Context(), user.Username)
	if existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户名已存在",
		})
		return
	}

	// 创建用户
	if err := s.userService.CreateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	// 清除敏感信息
	user.Password = ""
	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"data":    user,
		"message": "用户创建成功",
	})
}

// UpdateUser 更新用户信息
func (s *Server) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的请求数据",
		})
		return
	}

	// 不允许直接更新密码
	if _, exists := updates["password"]; exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "不能通过此接口更新密码",
		})
		return
	}

	// 检查是否更新用户名
	if username, exists := updates["username"]; exists {
		existingUser, _ := s.userService.GetUserByUsername(c.Request.Context(), username.(string))
		if existingUser != nil && existingUser.ID.Hex() != id {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "用户名已存在",
			})
			return
		}
	}

	if err := s.userService.UpdateUser(c.Request.Context(), id, bson.M(updates)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "用户信息更新成功",
	})
}

// DeleteUser 删除用户
func (s *Server) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := s.userService.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "用户删除成功",
	})
}

// GetUser 获取用户信息
func (s *Server) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := s.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "用户不存在",
		})
		return
	}

	// 清除敏感信息
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    user,
		"message": "获取成功",
	})
}

// ResetUserPassword 重置用户密码
func (s *Server) ResetUserPassword(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=6,max=32"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的密码格式",
		})
		return
	}

	if err := s.userService.UpdateUserPassword(c.Request.Context(), id, req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "密码重置成功",
	})
}
