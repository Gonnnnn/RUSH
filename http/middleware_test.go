package http

import (
	"net/http"
	"net/http/httptest"
	"rush/permission"
	"testing"

	// assert
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockServer struct {
	idToReturn    string
	roleToReturn  permission.Role
	errorToReturn error
}

func (s *mockServer) GetUserIdentifier(token string) (string, permission.Role, error) {
	return s.idToReturn, s.roleToReturn, s.errorToReturn
}

func TestUseAuthMiddleware(t *testing.T) {
	t.Run("Should return a gin.HandlerFunc", func(t *testing.T) {
		middleware := UseAuthMiddleware(&mockServer{"token", permission.RoleAdmin, nil})

		assert.NotNil(t, middleware)
		assert.IsType(t, gin.HandlerFunc(nil), middleware)
	})

	t.Run("Should return 401 when cookie is not found", func(t *testing.T) {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)

		middleware := UseAuthMiddleware(&mockServer{"", permission.RoleAdmin, nil})
		middleware(ctx)

		assert.Equal(t, http.StatusUnauthorized, ctx.Writer.Status())
	})

	t.Run("Should return 401 when token is empty", func(t *testing.T) {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)
		ctx.Request.AddCookie(&http.Cookie{Name: authCookieName, Value: ""})

		middleware := UseAuthMiddleware(&mockServer{"", permission.RoleAdmin, nil})
		middleware(ctx)

		assert.Equal(t, http.StatusUnauthorized, ctx.Writer.Status())
	})

	t.Run("Should return 401 when token is invalid", func(t *testing.T) {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)
		ctx.Request.AddCookie(&http.Cookie{Name: authCookieName, Value: "token"})

		middleware := UseAuthMiddleware(&mockServer{"", permission.RoleAdmin, assert.AnError})
		middleware(ctx)

		assert.Equal(t, http.StatusUnauthorized, ctx.Writer.Status())
	})

	t.Run("Should call next when token is valid", func(t *testing.T) {
		resRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(resRecorder)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)
		ctx.Request.AddCookie(&http.Cookie{Name: authCookieName, Value: "token"})

		middleware := UseAuthMiddleware(&mockServer{"user-id", permission.RoleMember, nil})
		middleware(ctx)

		assert.Equal(t, http.StatusOK, ctx.Writer.Status())
		assert.Equal(t, "user-id", ctx.GetString("userId"))
		userRole, exists := ctx.Get("userRole")
		assert.True(t, exists)
		assert.Equal(t, permission.RoleMember, userRole)

		middleware = UseAuthMiddleware(&mockServer{"user-id", permission.RoleAdmin, nil})
		middleware(ctx)

		assert.Equal(t, http.StatusOK, ctx.Writer.Status())
		assert.Equal(t, "user-id", ctx.GetString("userId"))
		userRole, exists = ctx.Get("userRole")
		assert.True(t, exists)
		assert.Equal(t, permission.RoleAdmin, userRole)
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
