package http

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"rush/server"
)

type UsersPostRequest struct {
	Name       string  `json:"name"`
	University string  `json:"university"`
	Phone      string  `json:"phone"`
	Generation float64 `json:"generation"`
	IsActive   bool    `json:"is_active"`
}

type SessionsPostRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartsAt    time.Time `json:"starts_at"`
	Score       int       `json:"score"`
}

func SetUpRouter(router *gin.Engine, server *server.Server) {
	api := router.Group("/api")
	{
		api.GET("/users", func(c *gin.Context) {
			offset, err := strconv.Atoi(c.Query("offset"))
			if err != nil || offset < 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
				return
			}
			pageSize, err := strconv.Atoi(c.Query("pageSize"))
			if err != nil || pageSize < 1 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pageSize"})
				return
			}

			result, err := server.ListUsers(offset, pageSize)
			if err != nil {
				log.Printf("Error listing users: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			response := gin.H{
				"users":       result.Users,
				"is_end":      result.IsEnd,
				"total_count": result.TotalCount,
			}

			c.JSON(http.StatusOK, response)
		})

		api.POST("/users", func(c *gin.Context) {
			var req UsersPostRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if err := server.AddUser(
				req.Name,
				req.University,
				req.Phone,
				req.Generation,
				req.IsActive,
			); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "User added successfully"})
		})

		api.GET("/sessions", func(c *gin.Context) {
			sessions, err := server.GetAllSessions()
			if err != nil {
				log.Printf("Error getting sessions: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, sessions)
		})

		api.GET("/sessions/:id", func(c *gin.Context) {
			id := c.Param("id")
			session, err := server.GetSession(id)
			if err != nil {
				log.Printf("Error getting session: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, session)
		})

		api.POST("/sessions", func(c *gin.Context) {
			var req SessionsPostRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			id, err := server.AddSession(req.Name, req.Description, req.StartsAt, req.Score)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"id": id})
		})

		api.POST("/sessions/:id/attendance-form", func(c *gin.Context) {
			sessionId := c.Param("id")
			formUrl, err := server.CreateSessionForm(sessionId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"form_url": formUrl})
		})

		api.GET("/attendances", func(c *gin.Context) {
			reports, err := server.GetAllAttendanceReports()
			if err != nil {
				log.Printf("Error getting attendance reports: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, reports)
		})
	}

	// TODO: Get it from env variable.
	staticDir := "./ui/dist"
	router.GET("/assets/*filepath", func(c *gin.Context) {
		http.FileServer(http.Dir(staticDir)).ServeHTTP(c.Writer, c.Request)
	})
	router.NoRoute(func(c *gin.Context) {
		c.File(filepath.Join(staticDir, "index.html"))
	})
}
