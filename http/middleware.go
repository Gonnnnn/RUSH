package http

import (
	"net/http"
	"rush/auth"

	"github.com/gin-gonic/gin"
)

type identifierFetcher interface {
	// Fetches the user identifier from the token to identify them.
	GetUserIdentifier(token string) (auth.UserIdentifier, error)
}

// TODO(#23): Attach it to the router once the UI is ready.
func ginAuthMiddleware(identifierFetcher identifierFetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		_, err := identifierFetcher.GetUserIdentifier(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
