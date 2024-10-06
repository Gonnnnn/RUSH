package http

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"rush/permission"
	"rush/server"
)

func SetUpRouter(router *gin.Engine, server *server.Server) {
	api := router.Group("/api")
	{
		api.POST("/sign-in", handleSignIn(server))
		api.GET("/users", handleListUsers(server))
		api.GET("/sessions", handleListSessions(server))
		api.GET("/sessions/:id", handleGetSession(server))

		protected := api.Group("/")
		protected.Use(UseAuthMiddleware(server))
		{
			// handleAuth doesn't immplement anything. It relies on the middleware to check the token.
			protected.GET("/auth", handleAuth(server))

			protected.GET("/users/:id/attendances", handleGetAttendanceForUser(server))
			protected.GET("/users/:id", handleGetUser(server))
			protected.POST("/users", handleAddUser(server))

			// TODO(#138): Move it to the admin group after fixing the UI to handle permission denied error on it more properly.
			protected.GET("attendances/half-year", handleHalfYearAttendance(server))

			adminProtected := protected.Group("/")
			adminProtected.Use(RequireRole(permission.RoleAdmin, permission.RoleSuperAdmin))
			{
				adminProtected.POST("/sessions", handleAddSession(server))
				adminProtected.DELETE("/sessions/:id", handleDeleteSession(server))
				adminProtected.POST("/sessions/:id/attendance-form", handleCreateAttendanceForm(server))
				adminProtected.POST("/attendances/aggregate", handleAggregateAttendance(server))
				adminProtected.POST("/sessions/:id/close", handleApplyAttendance(server))
				adminProtected.POST("/sessions/:id/present", handleMarkUsersAsPresent(server))
			}
		}
	}

	staticDir := "./ui/dist"
	router.GET("/assets/*filepath", func(c *gin.Context) {
		http.FileServer(http.Dir(staticDir)).ServeHTTP(c.Writer, c.Request)
	})
	router.NoRoute(func(c *gin.Context) {
		c.File(filepath.Join(staticDir, "index.html"))
	})
}
