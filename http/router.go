package http

import (
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

	"rush/server"
)

type UsersPostRequest struct {
	Name       string `json:"name"`
	University string `json:"university"`
	Phone      string `json:"phone"`
	Generation string `json:"generation"`
	IsActive   bool   `json:"is_active"`
}

type SessionsPostRequest struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	StartsAt    *customTime `json:"starts_at"`
	Score       int         `json:"score"`
}

type SessionAttendanceFormPostRequest struct {
	FormTitle       string `json:"form_title"`
	FormDescription string `json:"form_description"`
}

type customTime struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
	Hour  int `json:"hour"`
	Min   int `json:"min"`
}

func SetUpRouter(router *gin.Engine, server *server.Server) {
	api := router.Group("/api")
	{
		api.GET("/users", func(c *gin.Context) {
			users, err := server.GetAllUsers()
			if err != nil {
				log.Printf("Error getting users: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, users)
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

			id, err := server.AddSession(
				req.Name,
				req.Description,
				time.Date(req.StartsAt.Year, time.Month(req.StartsAt.Month), req.StartsAt.Day, req.StartsAt.Hour, req.StartsAt.Min, 0 /* =sec */, 0 /* =nsec */, time.UTC),
				req.Score,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"id": id})
		})

		api.POST("/sessions/:id/attendance-form", func(c *gin.Context) {
			var req SessionAttendanceFormPostRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			sessionId := c.Param("id")
			formUrl, err := server.CreateSessionForm(sessionId, req.FormTitle, req.FormDescription)
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
