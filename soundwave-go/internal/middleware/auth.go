package middleware

import (
	"net/http"
	"soundwave-go/internal/config"
	"soundwave-go/internal/models"
	"soundwave-go/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthRequired(cfg *config.Config, requiredPermission models.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
			c.Abort()
			return
		}

		// 从 Bearer token 中提取 token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token格式"})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(cfg, parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		// 检查权限
		hasPermission := false
		for _, p := range claims.Permissions {
			if p == requiredPermission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有访问权限"})
			c.Abort()
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}
