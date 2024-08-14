package http

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"rush/server"
)

func SetUpRouter(router *gin.Engine, server *server.Server) {
	api := router.Group("/api")
	{
		api.POST("/sign-in", handleSignIn(server))
		api.GET("/users", handleListUsers(server))
		api.GET("/sessions", handleListSessions(server))
		api.GET("/sessions/:id", handleGetSession(server))
		api.GET("attendances/half-year", handleHalfYearAttendance(server))

		protected := api.Group("/")
		protected.Use(UseAuthMiddleware(server))
		{
			// handleAuth doesn't immplement anything. It relies on the middleware to check the token.
			protected.GET("/auth", handleAuth(server))

			protected.GET("/users/:id/attendances", handleGetAttendanceForUser(server))
			protected.GET("/users/:id", handleGetUser(server))
			protected.POST("/users", handleAddUser(server))

			protected.POST("/sessions", handleAddSession(server))
			protected.POST("/sessions/:id/attendance-form", handleCreateAttendanceForm(server))
			protected.POST("/sessions/:id/attendance", handleApplyAttendance(server))

			// protected.GET("attendances/half-year", handleHalfYearAttendance(server))
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
