package middlewares

import (
	"net/http"
	"strings"

	"github.com/Jake4-CX/portfolio-website-v2-backend/pkg/initializers"
	"github.com/Jake4-CX/portfolio-website-v2-backend/pkg/structs"
	"github.com/Jake4-CX/portfolio-website-v2-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

func RoleMiddleware(requiredRole structs.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		userId, err := utils.ValidateToken(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userId", userId)

		// Check if the user has the required role - check from the database
		user := structs.Users{}

		result := initializers.DB.First(&user, userId)

		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		if user.UserRole != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}