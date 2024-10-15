package http

import (
	"net/http"
	"net/http/httptest"
	"rush/permission"
	"rush/server"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockServer struct {
	sessionToReturn  server.UserSession
	newTokenToReturn string
	errorToReturn    error
}

func (s *mockServer) GetUserSession(token string) (server.UserSession, string, error) {
	return s.sessionToReturn, s.newTokenToReturn, s.errorToReturn
}

func TestUseAuthMiddleware(t *testing.T) {
	t.Run("Should return a gin.HandlerFunc", func(t *testing.T) {
		middleware := UseAuthMiddleware(&mockServer{
			sessionToReturn:  server.UserSession{UserId: "user-id", Role: permission.RoleAdmin, ExpiresAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			newTokenToReturn: "new-token",
			errorToReturn:    nil,
		})

		assert.NotNil(t, middleware)
		assert.IsType(t, gin.HandlerFunc(nil), middleware)
	})

	t.Run("Should return 401 when cookie is not found", func(t *testing.T) {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)

		middleware := UseAuthMiddleware(&mockServer{
			sessionToReturn:  server.UserSession{UserId: "user-id", Role: permission.RoleAdmin, ExpiresAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			newTokenToReturn: "",
			errorToReturn:    nil,
		})
		middleware(ctx)

		assert.Equal(t, http.StatusUnauthorized, ctx.Writer.Status())
	})

	t.Run("Should return 401 when token is empty", func(t *testing.T) {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)
		ctx.Request.AddCookie(&http.Cookie{Name: authCookieName, Value: ""})

		middleware := UseAuthMiddleware(&mockServer{
			sessionToReturn:  server.UserSession{UserId: "user-id", Role: permission.RoleAdmin, ExpiresAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			newTokenToReturn: "",
			errorToReturn:    nil,
		})
		middleware(ctx)

		assert.Equal(t, http.StatusUnauthorized, ctx.Writer.Status())
	})

	t.Run("Should return 401 when token is invalid", func(t *testing.T) {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)
		ctx.Request.AddCookie(&http.Cookie{Name: authCookieName, Value: "token"})

		middleware := UseAuthMiddleware(&mockServer{
			sessionToReturn:  server.UserSession{UserId: "user-id", Role: permission.RoleAdmin, ExpiresAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			newTokenToReturn: "",
			errorToReturn:    &server.BadRequestError{},
		})
		middleware(ctx)

		assert.Equal(t, http.StatusUnauthorized, ctx.Writer.Status())
	})

	t.Run("Should return 500 when server returns internal server error", func(t *testing.T) {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)
		ctx.Request.AddCookie(&http.Cookie{Name: authCookieName, Value: "token"})

		middleware := UseAuthMiddleware(&mockServer{
			sessionToReturn:  server.UserSession{UserId: "user-id", Role: permission.RoleAdmin, ExpiresAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			newTokenToReturn: "",
			errorToReturn:    &server.InternalServerError{},
		})
		middleware(ctx)

		assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
	})

	t.Run("Should call next when token is valid", func(t *testing.T) {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)
		ctx.Request.AddCookie(&http.Cookie{Name: authCookieName, Value: "token"})

		middleware := UseAuthMiddleware(&mockServer{
			sessionToReturn:  server.UserSession{UserId: "user-id", Role: permission.RoleMember, ExpiresAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			newTokenToReturn: "new-token",
			errorToReturn:    nil,
		})
		middleware(ctx)

		assert.Equal(t, http.StatusOK, ctx.Writer.Status())
		assert.Equal(t, "user-id", ctx.GetString("userId"))
		userRole, exists := ctx.Get("userRole")
		assert.True(t, exists)
		assert.Equal(t, permission.RoleMember, userRole)

		middleware = UseAuthMiddleware(&mockServer{
			sessionToReturn:  server.UserSession{UserId: "user-id", Role: permission.RoleAdmin, ExpiresAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			newTokenToReturn: "new-token",
			errorToReturn:    nil,
		})
		middleware(ctx)

		assert.Equal(t, http.StatusOK, ctx.Writer.Status())
		assert.Equal(t, "user-id", ctx.GetString("userId"))
		userRole, exists = ctx.Get("userRole")
		assert.True(t, exists)
		assert.Equal(t, permission.RoleAdmin, userRole)
	})

	t.Run("Should pass a new token in the designated header when it's refreshed", func(t *testing.T) {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)
		ctx.Request.AddCookie(&http.Cookie{Name: authCookieName, Value: "token"})

		middleware := UseAuthMiddleware(&mockServer{
			sessionToReturn:  server.UserSession{UserId: "user-id", Role: permission.RoleMember, ExpiresAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			newTokenToReturn: "new-token",
			errorToReturn:    nil,
		})
		middleware(ctx)

		assert.Equal(t, http.StatusOK, ctx.Writer.Status())
		assert.Equal(t, "new-token", resRecorder.Header().Get(replaceCookieHeader))
	})
}

func TestRequireRole(t *testing.T) {
	t.Run("Should return a gin.HandlerFunc", func(t *testing.T) {
		middleware := RequireRole(permission.RoleAdmin)

		assert.NotNil(t, middleware)
		assert.IsType(t, gin.HandlerFunc(nil), middleware)
	})

	t.Run("Should return 500 when user role is not found", func(t *testing.T) {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)

		middleware := RequireRole(permission.RoleAdmin)
		middleware(ctx)

		assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
	})

	t.Run("Should return 500 when user role is invalid", func(t *testing.T) {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)
		ctx.Set(userRoleKey, "invalid")

		middleware := RequireRole(permission.RoleAdmin)
		middleware(ctx)

		assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
	})

	t.Run("Should return 403 when user role is not sufficient", func(t *testing.T) {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)
		ctx.Set(userRoleKey, permission.RoleMember)

		middleware := RequireRole(permission.RoleAdmin, permission.RoleSuperAdmin)
		middleware(ctx)

		assert.Equal(t, http.StatusForbidden, ctx.Writer.Status())
	})

	t.Run("Should call next when user role is sufficient", func(t *testing.T) {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)
		ctx.Set(userRoleKey, permission.RoleAdmin)

		middleware := RequireRole(permission.RoleAdmin, permission.RoleSuperAdmin)
		middleware(ctx)

		assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	})
}
