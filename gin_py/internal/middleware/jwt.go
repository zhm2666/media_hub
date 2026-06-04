package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 验证 JWT Token 的中间件
// 从 Authorization header 中获取 Token 并验证
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "缺少 Authorization header",
			})
			c.Abort()
			return
		}

		// 检查 Bearer token 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Authorization 格式错误，应为: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Token 不能为空",
			})
			c.Abort()
			return
		}

		// 将 token 传递给后续处理函数
		// 实际的 token 验证由调用 user 服务的 API 来完成
		c.Set("token", tokenString)

		c.Next()
	}
}

// OptionalJWTAuthMiddleware 可选的 JWT 验证中间件
// 如果有 token 就验证，没有 token 也不阻止请求
func OptionalJWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			tokenString := parts[1]
			if tokenString != "" {
				c.Set("token", tokenString)
			}
		}

		c.Next()
	}
}
