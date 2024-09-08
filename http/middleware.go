package http

import (
	"net/http"
	"rush/golang/array"
	"rush/permission"

	"github.com/gin-gonic/gin"
)

type userIdFetcher interface {
	GetUserIdentifier(token string) (string, permission.Role, error)
}

// The name of the cookie to store the rush authentication token.
const authCookieName = "rush-auth"
const userIdKey = "userId"
const userRoleKey = "userRole"

func UseAuthMiddleware(userIdFetcher userIdFetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(authCookieName)
		// https://pkg.go.dev/github.com/gin-gonic/gin#Context.Cookie
		// It only returns ErrNoCookie for errors.
		if err == http.ErrNoCookie {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Cookie not found"})
			c.Abort()
			return
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Empty token"})
			c.Abort()
			return
		}

		userId, role, err := userIdFetcher.GetUserIdentifier(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set(userIdKey, userId)
		c.Set(userRoleKey, role)
		c.Next()
	}
}

func RequireRole(requiredRoles ...permission.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get(userRoleKey)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User role should be found but it has not"})
			c.Abort()
			return
		}

		role, ok := userRole.(permission.Role)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user role"})
			c.Abort()
			return
		}

		if array.Contains(requiredRoles, role) {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}
