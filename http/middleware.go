package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type tokenValidator interface {
	IsTokenValid(token string) bool
}

var authCookieName = "rush-auth"

func UseAuthMiddleware(tokenValidator tokenValidator) gin.HandlerFunc {
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

		if !tokenValidator.IsTokenValid(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
