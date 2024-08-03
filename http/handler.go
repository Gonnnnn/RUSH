package http

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"rush/server"
)

type SignInRequest struct {
	Token string `json:"token"`
}

func handleSignIn(server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SignInRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, err := server.SignIn(req.Token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func handleAuth(server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Authorized"})
	}
}

func handleListUsers(server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		c.JSON(http.StatusOK, gin.H{
			"users":       result.Users,
			"is_end":      result.IsEnd,
			"total_count": result.TotalCount,
		})
	}
}

type addUserRequest struct {
	Name       string  `json:"name"`
	University string  `json:"university"`
	Phone      string  `json:"phone"`
	Generation float64 `json:"generation"`
	IsActive   bool    `json:"is_active"`
}

func handleAddUser(server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req addUserRequest
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
	}
}

func handleListSessions(server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		result, err := server.ListSessions(offset, pageSize)
		if err != nil {
			log.Printf("Error getting sessions: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"sessions":    result.Sessions,
			"is_end":      result.IsEnd,
			"total_count": result.TotalCount,
		})
	}
}

func handleGetSession(server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		session, err := server.GetSession(id)
		if err != nil {
			log.Printf("Error getting session: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, session)
	}
}

type addSessionRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartsAt    time.Time `json:"starts_at"`
	Score       int       `json:"score"`
}

func handleAddSession(server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req addSessionRequest
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
	}
}

func handleCreateSessionForm(server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId := c.Param("id")
		formUrl, err := server.CreateSessionForm(sessionId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"form_url": formUrl})
	}
}
