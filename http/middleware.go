package http

import (
	"net/http"
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
