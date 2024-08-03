package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	// assert
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockServer struct {
	valueToReturn bool
}

func (s *mockServer) IsTokenValid(token string) bool {
	return s.valueToReturn
}

func TestUseAuthMiddleware(t *testing.T) {
	t.Run("Should return a gin.HandlerFunc", func(t *testing.T) {
		middleware := UseAuthMiddleware(&mockServer{true})

		assert.NotNil(t, middleware)
		assert.IsType(t, gin.HandlerFunc(nil), middleware)
	})

	t.Run("Should return 401 when cookie is not found", func(t *testing.T) {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)

		middleware := UseAuthMiddleware(&mockServer{true})
		middleware(ctx)

		assert.Equal(t, http.StatusUnauthorized, ctx.Writer.Status())
	})

	t.Run("Should return 401 when token is empty", func(t *testing.T) {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)
		ctx.Request.AddCookie(&http.Cookie{Name: authCookieName, Value: ""})

		middleware := UseAuthMiddleware(&mockServer{true})
		middleware(ctx)

		assert.Equal(t, http.StatusUnauthorized, ctx.Writer.Status())
	})

	t.Run("Should return 401 when token is invalid", func(t *testing.T) {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)
		ctx.Request.AddCookie(&http.Cookie{Name: authCookieName, Value: "token"})

		middleware := UseAuthMiddleware(&mockServer{false})
		middleware(ctx)

		assert.Equal(t, http.StatusUnauthorized, ctx.Writer.Status())
	})

	t.Run("Should call next when token is valid", func(t *testing.T) {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)
		ctx.Request.AddCookie(&http.Cookie{Name: authCookieName, Value: "token"})

		middleware := UseAuthMiddleware(&mockServer{true})
		middleware(ctx)

		assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	})
}
