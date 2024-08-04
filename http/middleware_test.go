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
	valueToReturn string
	errorToReturn error
}

func (s *mockServer) GetUserIdentifier(token string) (string, error) {
	return s.valueToReturn, s.errorToReturn
}

func TestUseAuthMiddleware(t *testing.T) {
	t.Run("Should return a gin.HandlerFunc", func(t *testing.T) {
		middleware := UseAuthMiddleware(&mockServer{"token", nil})

		assert.NotNil(t, middleware)
		assert.IsType(t, gin.HandlerFunc(nil), middleware)
	})

	t.Run("Should return 401 when cookie is not found", func(t *testing.T) {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)

		middleware := UseAuthMiddleware(&mockServer{"", nil})
		middleware(ctx)

		assert.Equal(t, http.StatusUnauthorized, ctx.Writer.Status())
	})

	t.Run("Should return 401 when token is empty", func(t *testing.T) {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)
		ctx.Request.AddCookie(&http.Cookie{Name: authCookieName, Value: ""})

		middleware := UseAuthMiddleware(&mockServer{"", assert.AnError})
		middleware(ctx)

		assert.Equal(t, http.StatusUnauthorized, ctx.Writer.Status())
	})

	t.Run("Should return 401 when token is invalid", func(t *testing.T) {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)
		ctx.Request.AddCookie(&http.Cookie{Name: authCookieName, Value: "token"})

		middleware := UseAuthMiddleware(&mockServer{"", assert.AnError})
		middleware(ctx)

		assert.Equal(t, http.StatusUnauthorized, ctx.Writer.Status())
	})

	t.Run("Should call next when token is valid", func(t *testing.T) {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)
		ctx.Request.AddCookie(&http.Cookie{Name: authCookieName, Value: "token"})

		middleware := UseAuthMiddleware(&mockServer{"token", nil})
		middleware(ctx)

		assert.Equal(t, http.StatusOK, ctx.Writer.Status())
		assert.Equal(t, "token", ctx.GetString("userId"))
	})
}
