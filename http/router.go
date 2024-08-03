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
		api.GET("/auth", handleAuth(server))
		api.GET("/users", handleListUsers(server))
		api.GET("/sessions", handleListSessions(server))
		api.GET("/sessions/:id", handleGetSession(server))

		protected := api.Group("/")
		protected.Use(UseAuthMiddleware(server))
		{
			protected.POST("/users", handleAddUser(server))
			protected.POST("/sessions", handleAddSession(server))
			protected.POST("/sessions/:id/attendance-form", handleCreateSessionForm(server))
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
