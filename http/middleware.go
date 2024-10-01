package http

import (
	"net/http"
	"rush/golang/array"
	"rush/permission"
	"rush/server"

	"github.com/gin-gonic/gin"
)

type userSessionFetcher interface {
	GetUserSession(token string) (server.UserSession, string, error)
}

// The name of the cookie to store the rush authentication token.
const authCookieName = "rush-auth"
const userIdKey = "userId"
const userRoleKey = "userRole"
const replaceCookieHeader = "X-Replace-Cookie"

func UseAuthMiddleware(userSessionFetcher userSessionFetcher) gin.HandlerFunc {
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

		userSession, newToken, err := userSessionFetcher.GetUserSession(token)
		if err != nil {
			if getHttpStatus(err) == http.StatusInternalServerError {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if newToken != "" {
			c.Header(replaceCookieHeader, newToken)
		}

		c.Set(userIdKey, userSession.UserId)
		c.Set(userRoleKey, userSession.Role)
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
