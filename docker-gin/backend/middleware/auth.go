package middleware

import (
	"strings"
	"wordbook/utils"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(401, gin.H{"error": "未提供token"})
			c.Abort()
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "token格式错误"})
			c.Abort()
			return
		}

		userID, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(401, gin.H{"error": "token无效"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
